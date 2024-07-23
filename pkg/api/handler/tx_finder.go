package handler

import (
	"context"
	"encoding/hex"
	"log/slog"

	"github.com/bitcoin-sv/arc/config"
	"github.com/bitcoin-sv/arc/internal/validator"
	"github.com/bitcoin-sv/arc/internal/woc_client"
	"github.com/bitcoin-sv/arc/pkg/metamorph"
)

type txFinder struct {
	th         metamorph.TransactionHandler
	pc         *config.PeerRpcConfig
	l          *slog.Logger
	w          *woc_client.WocClient
	useMainnet bool
}

func (f txFinder) GetRawTxs(ctx context.Context, source validator.FindSourceFlag, ids []string) ([]validator.RawTx, error) {
	// NOTE: we can ignore ALL errors from providers, if one returns err we go to another

	// TODO: discuss if it's worth to have implementation for len(ids) == 1

	foundTxs := make([]validator.RawTx, 0, len(ids))
	var remainingIDs []string

	// first get transactions from the handler
	if source.Has(validator.SourceTransactionHandler) {
		txs, _ := f.th.GetManyTransactions(ctx, ids)
		for _, tx := range txs {
			rt := validator.RawTx{
				TxID:    tx.TxID,
				Bytes:   tx.Bytes,
				IsMined: tx.BlockHeight > 0,
			}

			foundTxs = append(foundTxs, rt)
		}

		// add remaining ids
		for _, id := range ids {
			found := false
			for _, tx := range foundTxs {
				if tx.TxID == id {
					found = true
					break
				}
			}

			if !found {
				remainingIDs = append(remainingIDs, id)
			}
		}
	}

	ids = remainingIDs[:]
	remainingIDs = make([]string, 0)

	// try to get remaining txs from the node
	if source.Has(validator.SourceNodes) {
		for _, id := range ids {
			nTx, err := getTransactionFromNode(f.pc, id)
			if err != nil {
				// do we really need this info?
				f.l.Warn("failed to get transaction from node", slog.String("id", id), slog.String("err", err.Error()))
			}
			if nTx != nil {
				rt, e := newRawTx(nTx.TxID, nTx.Hex, nTx.BlockHeight)
				if e != nil {
					return nil, e
				}

				foundTxs = append(foundTxs, rt)
			} else {
				remainingIDs = append(remainingIDs, id)
			}
		}
	}

	// at last try the WoC
	if source.Has(validator.SourceWoC) {
		wocTxs, _ := f.w.GetRawTxs(ctx, f.useMainnet, remainingIDs)
		for _, wTx := range wocTxs {
			rt, e := newRawTx(wTx.TxID, wTx.Hex, wTx.BlockHeight)
			if e != nil {
				return nil, e
			}

			foundTxs = append(foundTxs, rt)
		}
	}

	return foundTxs, nil
}

func newRawTx(id, hexTx string, blockH uint64) (validator.RawTx, error) {
	b, e := hex.DecodeString(hexTx)
	if e != nil {
		return validator.RawTx{}, e
	}

	rt := validator.RawTx{
		TxID:    id,
		Bytes:   b,
		IsMined: blockH > 0,
	}

	return rt, nil
}
