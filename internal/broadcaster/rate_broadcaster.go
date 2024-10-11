package broadcaster

import (
	"context"
	cRand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	sdkTx "github.com/bitcoin-sv/go-sdk/transaction"

	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/pkg/keyset"
)

var (
	ErrFailedToGetBalance       = errors.New("failed to get balance")
	ErrKeyHasUnconfirmedBalance = errors.New("key has unconfirmed balance")
	ErrFailedToGetUTXOs         = errors.New("failed to get utxos")
	ErrTooHighSubmissionRate    = errors.New("submission rate is too high")
	ErrTooSmallUTXOSet          = errors.New("utxo set is too small")
	ErrFailedToAddInput         = errors.New("failed to add input")
	ErrFailedToAddOutput        = errors.New("failed to add output")
	ErrNotEnoughUTXOs           = errors.New("not enough utxos with sufficient funds left")
	ErrNotEnoughUTXOsForBatch   = errors.New("not enough utxos with sufficient funds left for another batch")
	ErrFailedToFillInputs       = errors.New("failed to fill inputs")
)

type UTXORateBroadcaster struct {
	Broadcaster
	totalTxs         int64
	connectionCount  int64
	shutdown         chan struct{}
	utxoCh           chan *sdkTx.UTXO
	wg               sync.WaitGroup
	satoshiMap       sync.Map
	ks               *keyset.KeySet
	rateTxsPerSecond int
	limit            int64
}

func NewRateBroadcaster(logger *slog.Logger, client ArcClient, ks *keyset.KeySet, utxoClient UtxoClient, isTestnet bool, rateTxsPerSecond int, limit int64, opts ...func(p *Broadcaster)) (*UTXORateBroadcaster, error) {
	b, err := NewBroadcaster(logger, client, utxoClient, isTestnet, opts...)
	if err != nil {
		return nil, err
	}
	rb := &UTXORateBroadcaster{
		Broadcaster:      b,
		shutdown:         make(chan struct{}, 1),
		utxoCh:           nil,
		wg:               sync.WaitGroup{},
		satoshiMap:       sync.Map{},
		ks:               ks,
		totalTxs:         0,
		connectionCount:  0,
		rateTxsPerSecond: rateTxsPerSecond,
		limit:            limit,
	}

	return rb, nil
}

func (b *UTXORateBroadcaster) Start() error {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.shutdown:
				b.cancelAll()
			case <-b.ctx.Done():
				return
			}
		}
	}()

	_, unconfirmed, err := b.utxoClient.GetBalanceWithRetries(b.ctx, b.ks.Address(!b.isTestnet), 1*time.Second, 5)
	if err != nil {
		return errors.Join(ErrFailedToGetBalance, err)
	}
	if math.Abs(float64(unconfirmed)) > 0 {
		return errors.Join(ErrKeyHasUnconfirmedBalance, fmt.Errorf("address %s, unconfirmed amount %d", b.ks.Address(!b.isTestnet), unconfirmed))
	}
	b.logger.Info("Start broadcasting", slog.String("wait for status", b.waitForStatus.String()), slog.String("op return", b.opReturn))

	utxoSet, err := b.utxoClient.GetUTXOsWithRetries(b.ctx, b.ks.Script, b.ks.Address(!b.isTestnet), 1*time.Second, 5)
	if err != nil {
		return errors.Join(ErrFailedToGetUTXOs, err)
	}

	submitBatchesPerSecond := float64(b.rateTxsPerSecond) / float64(b.batchSize)

	if submitBatchesPerSecond > millisecondsPerSecond {
		return errors.Join(ErrTooHighSubmissionRate, fmt.Errorf("submission rate %d [txs/s] and batch size %d [txs] result in submission frequency %.2f greater than 1000 [/s]", b.rateTxsPerSecond, b.batchSize, submitBatchesPerSecond))
	}

	if len(utxoSet) < b.batchSize {
		return errors.Join(ErrTooSmallUTXOSet, fmt.Errorf("size of utxo set %d is smaller than requested batch size %d - create more utxos first", len(utxoSet), b.batchSize))
	}

	b.utxoCh = make(chan *sdkTx.UTXO, 100000)
	for _, utxo := range utxoSet {
		b.utxoCh <- utxo
	}

	submitBatchInterval := time.Duration(millisecondsPerSecond/float64(submitBatchesPerSecond)) * time.Millisecond
	submitBatchTicker := time.NewTicker(submitBatchInterval)

	errCh := make(chan error, 100)

	b.wg.Add(1)
	go func() {
		defer func() {
			b.logger.Info("shutting down broadcaster")
			b.wg.Done()
		}()

		for {
			select {
			case <-b.ctx.Done():
				return
			case <-submitBatchTicker.C:

				txs, err := b.createSelfPayingTxs()
				if err != nil {
					b.logger.Error("failed to create self paying txs", slog.String("err", err.Error()))
					b.shutdown <- struct{}{}
					continue
				}

				if b.limit > 0 && atomic.LoadInt64(&b.totalTxs) >= b.limit {
					b.logger.Info("limit reached", slog.Int64("total", atomic.LoadInt64(&b.totalTxs)), slog.Int64("limit", b.limit))
					b.shutdown <- struct{}{}
				}

				b.broadcastBatchAsync(txs, errCh, b.waitForStatus)

			case responseErr := <-errCh:
				b.logger.Error("failed to submit transactions", slog.String("err", responseErr.Error()))
			}
		}
	}()

	return nil
}

func (b *UTXORateBroadcaster) createSelfPayingTxs() (sdkTx.Transactions, error) {
	txs := make(sdkTx.Transactions, 0, b.batchSize)

utxoLoop:
	for {
		select {
		case <-b.ctx.Done():
			return txs, nil
		case utxo := <-b.utxoCh:
			tx := sdkTx.NewTransaction()
			amount := utxo.Satoshis

			err := tx.AddInputsFromUTXOs(utxo)
			if err != nil {
				return nil, errors.Join(ErrFailedToAddInput, err)
			}

			if b.opReturn != "" {
				err = tx.AddOpReturnOutput([]byte(b.opReturn))
				if err != nil {
					return nil, fmt.Errorf("failed to add OP_RETURN output: %v", err)
				}
			}

			if b.sizeJitterMax > 0 {
				// Add additional inputs to the transaction
				src := rand.NewSource(time.Now().UnixNano())
				r := rand.New(src)
				numOfInputs := r.Intn(10)

				for i := 0; i < numOfInputs; i++ {
					additionalUtxo, ok := <-b.utxoCh
					if !ok {
						return nil, ErrNotEnoughUTXOsForBatch
					}

					err = tx.AddInputsFromUTXOs(additionalUtxo)
					if err != nil {
						return nil, errors.Join(ErrFailedToAddInput, err)
					}

					amount += additionalUtxo.Satoshis
				}

				// Add additional OP_RETURN with random data to the transaction
				dataSize := r.Intn(b.sizeJitterMax)
				if err != nil {
					return nil, fmt.Errorf("failed to generate random number for filling OP_RETURN: %v", err)
				}

				randomBytes := make([]byte, dataSize)
				_, err = cRand.Read(randomBytes)
				if err != nil {
					return nil, fmt.Errorf("failed to fill OP_RETURN with random bytes: %v", err)
				}

				testHeader := []byte(" size-jitter random bytes - ")
				if b.opReturn != "" {
					testHeader = append([]byte(b.opReturn), testHeader...)
				}

				err = tx.AddOpReturnOutput(append(testHeader, randomBytes...))
				if err != nil {
					return nil, fmt.Errorf("failed to add OP_RETURN output: %v", err)
				}
			}

			fee, err := b.EstimateFee(tx)
			if err != nil {
				return nil, err
			}

			if amount <= fee {
				if len(b.utxoCh) == 0 {
					return nil, ErrNotEnoughUTXOs
				}

				if len(b.utxoCh) < b.batchSize {
					return nil, ErrNotEnoughUTXOsForBatch
				}

				continue
			}
			amount -= fee

			err = PayTo(tx, b.ks.Script, amount)
			if err != nil {
				return nil, errors.Join(ErrFailedToAddOutput, err)
			}

			err = SignAllInputs(tx, b.ks.PrivateKey)
			if err != nil {
				return nil, errors.Join(ErrFailedToFillInputs, err)
			}

			b.satoshiMap.Store(tx.TxID(), tx.Outputs[0].Satoshis)

			txs = append(txs, tx)

			if len(txs) >= b.batchSize {
				break utxoLoop
			}
		}
	}

	return txs, nil
}

// EstimateFee estimates the fee for a transaction
// based on the estimated size of the transaction
func (b *UTXORateBroadcaster) EstimateFee(tx *sdkTx.Transaction) (uint64, error) {
	size := EstimateSize(tx)

	fee, err := b.feeModel.ComputeFeeBasedOnSize(uint64(size))
	if err != nil {
		return 0, err
	}

	return fee, nil
}

func (b *UTXORateBroadcaster) broadcastBatchAsync(txs sdkTx.Transactions, errCh chan error, waitForStatus metamorph_api.Status) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()

		ctx, cancel := context.WithTimeout(b.ctx, 20*time.Second)
		defer cancel()

		atomic.AddInt64(&b.connectionCount, 1)

		resp, err := b.client.BroadcastTransactions(ctx, txs, waitForStatus, b.callbackURL, b.callbackToken, b.fullStatusUpdates, false)
		if err != nil {

			// In case of error put utxos back in channel
			for _, tx := range txs {
				for _, input := range tx.Inputs {
					unusedUtxo := &sdkTx.UTXO{
						TxID:          input.SourceTXID,
						Vout:          0,
						LockingScript: b.ks.Script,
						Satoshis:      *input.SourceTxSatoshis(),
					}
					b.utxoCh <- unusedUtxo
				}
			}

			if errors.Is(err, context.Canceled) {
				atomic.AddInt64(&b.connectionCount, -1)
				return
			}
			errCh <- err
		}

		atomic.AddInt64(&b.connectionCount, -1)

		for _, res := range resp {

			txIDBytes, err := hex.DecodeString(res.Txid)
			if err != nil {
				b.logger.Error("failed to decode txid", slog.String("err", err.Error()))
				continue
			}

			sat, found := b.satoshiMap.Load(res.Txid)
			satoshis, isValid := sat.(uint64)

			if found && isValid {
				newUtxo := &sdkTx.UTXO{
					TxID:          txIDBytes,
					Vout:          0,
					LockingScript: b.ks.Script,
					Satoshis:      satoshis,
				}
				b.utxoCh <- newUtxo
			}

			b.satoshiMap.Delete(res.Txid)

			atomic.AddInt64(&b.totalTxs, 1)
		}
	}()
}

func (b *UTXORateBroadcaster) Shutdown() {
	b.cancelAll()

	b.wg.Wait()
}

func (b *UTXORateBroadcaster) Wait() {
	b.wg.Wait()
}

func (b *UTXORateBroadcaster) GetLimit() int64 {
	return b.limit
}

func (b *UTXORateBroadcaster) GetTxCount() int64 {
	return atomic.LoadInt64(&b.totalTxs)
}

func (b *UTXORateBroadcaster) GetConnectionCount() int64 {
	return atomic.LoadInt64(&b.connectionCount)
}

func (b *UTXORateBroadcaster) GetUtxoSetLen() int {
	return len(b.utxoCh)
}
