package broadcaster

import (
	"container/list"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/pkg/keyset"
	sdkTx "github.com/bitcoin-sv/go-sdk/transaction"
)

type UTXOCreator struct {
	Broadcaster
	keySet *keyset.KeySet 
	wg     sync.WaitGroup
}

func (b *UTXOCreator) Shutdown() {
	
	b.cancelAll() 
	
	b.wg.Wait()
}
func NewUTXOCreator(logger *slog.Logger, client ArcClient, keySet *keyset.KeySet, utxoClient UtxoClient, isTestnet bool, opts ...func(p *Broadcaster)) (*UTXOCreator, error) {
	b, err := NewBroadcaster(logger, client, utxoClient, isTestnet, opts...)
	if err != nil {
		return nil, err
	}

	creator := &UTXOCreator{
		Broadcaster: b,
		keySet:      keySet, 
	}

	return creator, nil
}

func (b *UTXOCreator) CreateUtxos(requestedOutputs int, requestedSatoshisPerOutput uint64) error {
	b.logger.Info("creating utxos", slog.String("address", b.keySet.Address(!b.isTestnet)))

	requestedOutputsSatoshis := int64(requestedOutputs) * int64(requestedSatoshisPerOutput)

	confirmed, unconfirmed, err := b.utxoClient.GetBalanceWithRetries(b.ctx, b.keySet.Address(!b.isTestnet), 1*time.Second, 5)
	if err != nil {
		return err
	}

	balance := confirmed + unconfirmed

	if requestedOutputsSatoshis > balance {
		return fmt.Errorf("requested total of satoshis %d exceeds balance on funding keyset %d", requestedOutputsSatoshis, balance)
	}

	utxos, err := b.utxoClient.GetUTXOsWithRetries(b.ctx, b.keySet.Script, b.keySet.Address(!b.isTestnet), 1*time.Second, 5)
	if err != nil {
		return err
	}

	utxoSet := list.New()
	for _, utxo := range utxos {
		
		if utxo.Satoshis >= requestedSatoshisPerOutput {
			utxoSet.PushBack(utxo)
			continue
		}
	}

	
	if utxoSet.Len() >= requestedOutputs {
		b.logger.Info("utxo set", slog.Int("ready", utxoSet.Len()), slog.Int("requested", requestedOutputs), slog.Uint64("satoshis", requestedSatoshisPerOutput))
		return nil
	}

	satoshiMap := map[string][]splittingOutput{}
	lastUtxoSetLen := 0

	
	for {
		if lastUtxoSetLen >= utxoSet.Len() {
			b.logger.Error("utxo set length hasn't changed since last iteration")
			break
		}
		lastUtxoSetLen = utxoSet.Len()

		
		if utxoSet.Len() >= requestedOutputs {
			break
		}

		b.logger.Info("splitting outputs", slog.Int("ready", utxoSet.Len()), slog.Int("requested", requestedOutputs), slog.Uint64("satoshis", requestedSatoshisPerOutput))

		// create splitting txs
		txsSplitBatches, err := b.splitOutputs(requestedOutputs, requestedSatoshisPerOutput, utxoSet, satoshiMap, b.keySet)
		if err != nil {
			return err
		}

		for i, batch := range txsSplitBatches {
			nrOutputs := 0
			nrInputs := 0
			for _, txBatch := range batch {
				nrOutputs += len(txBatch.Outputs)
				nrInputs += len(txBatch.Inputs)
			}

			b.logger.Info(fmt.Sprintf("broadcasting splitting batch %d/%d", i+1, len(txsSplitBatches)), slog.Int("size", len(batch)), slog.Int("inputs", nrInputs), slog.Int("outputs", nrOutputs))

			resp, err := b.client.BroadcastTransactions(context.Background(), batch, metamorph_api.Status_SEEN_ON_NETWORK, "", "", false, false)
			if err != nil {
				return fmt.Errorf("failed to broadcast tx: %v", err)
			}

			for _, res := range resp {
				if res.Status == metamorph_api.Status_REJECTED || res.Status == metamorph_api.Status_SEEN_IN_ORPHAN_MEMPOOL {
					b.logger.Error("splitting tx was not successful", slog.String("status", res.Status.String()), slog.String("hash", res.Txid), slog.String("reason", res.RejectReason))
					for _, tx := range batch {
						if tx.TxID() == res.Txid {
							b.logger.Debug(tx.String())
							break
						}
					}
					continue
				}

				txIDBytes, err := hex.DecodeString(res.Txid)
				if err != nil {
					b.logger.Error("failed to decode txid", slog.String("err", err.Error()))
					continue
				}

				foundOutputs, found := satoshiMap[res.Txid]
				if !found {
					b.logger.Error("output not found", slog.String("hash", res.Txid))
					continue
				}

				for _, foundOutput := range foundOutputs {
					newUtxo := &sdkTx.UTXO{
						TxID:          txIDBytes,
						Vout:          foundOutput.vout,
						LockingScript: b.keySet.Script,
						Satoshis:      foundOutput.satoshis,
					}

					utxoSet.PushBack(newUtxo)
				}
				delete(satoshiMap, res.Txid)
			}

			// do not performance test ARC when creating the utxos
			time.Sleep(100 * time.Millisecond)
		}
	}

	b.logger.Info("utxo set", slog.Int("ready", utxoSet.Len()), slog.Int("requested", requestedOutputs), slog.Uint64("satoshis", requestedSatoshisPerOutput))
	return nil
}

func (b *UTXOCreator) splitOutputs(requestedOutputs int, requestedSatoshisPerOutput uint64, utxoSet *list.List, satoshiMap map[string][]splittingOutput, fundingKeySet *keyset.KeySet) ([]sdkTx.Transactions, error) {
	txsSplitBatches := make([]sdkTx.Transactions, 0)
	txsSplit := make(sdkTx.Transactions, 0)
	outputs := utxoSet.Len()
	var err error

	var next *list.Element
	for front := utxoSet.Front(); front != nil; front = next {
		next = front.Next()

		if outputs >= requestedOutputs {
			break
		}

		utxo, ok := front.Value.(*sdkTx.UTXO)
		if !ok {
			return nil, errors.New("failed to parse value to utxo")
		}

		tx := sdkTx.NewTransaction()
		err = tx.AddInputsFromUTXOs(utxo)
		if err != nil {
			return nil, err
		}

		// only split if splitting increases nr of outputs
		const feeMargin = 50
		if utxo.Satoshis < 2*requestedSatoshisPerOutput+feeMargin {
			continue
		}

		addedOutputs, err := b.splitToFundingKeyset(tx, utxo.Satoshis, requestedSatoshisPerOutput, requestedOutputs-outputs, fundingKeySet)
		if err != nil {
			return nil, err
		}
		utxoSet.Remove(front)

		outputs += addedOutputs

		txsSplit = append(txsSplit, tx)

		txOutputs := make([]splittingOutput, len(tx.Outputs))
		for i, txOutput := range tx.Outputs {
			txOutputs[i] = splittingOutput{satoshis: txOutput.Satoshis, vout: uint32(i)}
		}

		satoshiMap[tx.TxID()] = txOutputs

		if len(txsSplit) == b.batchSize {
			txsSplitBatches = append(txsSplitBatches, txsSplit)
			txsSplit = make(sdkTx.Transactions, 0)
		}
	}

	if len(txsSplit) > 0 {
		txsSplitBatches = append(txsSplitBatches, txsSplit)
	}
	return txsSplitBatches, nil
}

func (b *UTXOCreator) splitToFundingKeyset(tx *sdkTx.Transaction, splitSatoshis uint64, requestedSatoshis uint64, requestedOutputs int, fundingKeySet *keyset.KeySet) (int, error) {
	if requestedSatoshis > splitSatoshis {
		return 0, fmt.Errorf("requested satoshis %d greater than satoshis to be split %d", requestedSatoshis, splitSatoshis)
	}

	counter := 0
	var err error
	var fee uint64

	remaining := int64(splitSatoshis)

	for remaining > int64(requestedSatoshis) && counter < requestedOutputs {
		fee, err = b.feeModel.ComputeFee(tx)
		if err != nil {
			return 0, err
		}
		if uint64(remaining)-requestedSatoshis < fee {
			break
		}

		err = PayTo(tx, fundingKeySet.Script, requestedSatoshis)
		if err != nil {
			return 0, err
		}

		remaining -= int64(requestedSatoshis)
		counter++
	}

	fee, err = b.feeModel.ComputeFee(tx)
	if err != nil {
		return 0, err
	}

	err = PayTo(tx, fundingKeySet.Script, uint64(remaining)-fee)
	if err != nil {
		return 0, err
	}

	err = SignAllInputs(tx, fundingKeySet.PrivateKey)
	if err != nil {
		return 0, err
	}

	return counter, nil
}

type splittingOutput struct {
	satoshis uint64
	vout     uint32
}
