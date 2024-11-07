package txfinder

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"runtime"

	"github.com/ordishs/go-bitcoin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/bitcoin-sv/arc/config"
	"github.com/bitcoin-sv/arc/internal/tracing"
	"github.com/bitcoin-sv/arc/internal/validator"
	"github.com/bitcoin-sv/arc/internal/woc_client"
	"github.com/bitcoin-sv/arc/pkg/metamorph"
)

type Finder struct {
	transactionHandler metamorph.TransactionHandler
	bitcoinClient      *bitcoin.Bitcoind
	wocClient          *woc_client.WocClient
	logger             *slog.Logger
	tracingEnabled     bool
	tracingAttributes  []attribute.KeyValue
}

func WithTracerFinder(attr ...attribute.KeyValue) func(s *Finder) {
	return func(p *Finder) {
		p.tracingEnabled = true
		if len(attr) > 0 {
			p.tracingAttributes = append(p.tracingAttributes, attr...)
		}
		_, file, _, ok := runtime.Caller(1)
		if ok {
			p.tracingAttributes = append(p.tracingAttributes, attribute.String("file", file))
		}
	}
}

func New(th metamorph.TransactionHandler, pc *config.PeerRPCConfig, w *woc_client.WocClient, l *slog.Logger, opts ...func(f *Finder)) *Finder {
	l = l.With(slog.String("module", "tx-finder"))
	var n *bitcoin.Bitcoind

	if pc != nil {
		rpcURL, err := url.Parse(fmt.Sprintf("rpc://%s:%s@%s:%d", pc.User, pc.Password, pc.Host, pc.Port))
		if err != nil {
			l.Warn("cannot parse node rpc url. Finder will not use node as source")
		} else {
			// get the transaction from the bitcoin node rpc
			n, err = bitcoin.NewFromURL(rpcURL, false)
			if err != nil {
				l.Warn("cannot create node client. Finder will not use node as source")
			}
		}
	}

	f := &Finder{
		transactionHandler: th,
		bitcoinClient:      n,
		wocClient:          w,
		logger:             l,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f Finder) GetRawTxs(ctx context.Context, source validator.FindSourceFlag, ids []string) ([]validator.RawTx, error) {
	ctx, span := tracing.StartTracing(ctx, "Finder_GetRawTxs", f.tracingEnabled, f.tracingAttributes...)
	defer tracing.EndTracing(span)

	// NOTE: we can ignore ALL errors from providers, if one returns err we go to another
	foundTxs := make([]validator.RawTx, 0, len(ids))
	var remainingIDs []string

	// first get transactions from the handler
	if source.Has(validator.SourceTransactionHandler) {
		var thGetRawTxSpan trace.Span
		ctx, thGetRawTxSpan = tracing.StartTracing(ctx, "TransactionHandler_GetRawTxs", f.tracingEnabled, f.tracingAttributes...)
		txs, thErr := f.transactionHandler.GetTransactions(ctx, ids)
		tracing.EndTracing(thGetRawTxSpan)
		for _, tx := range txs {
			rt := validator.RawTx{
				TxID:    tx.TxID,
				Bytes:   tx.Bytes,
				IsMined: tx.BlockHeight > 0,
			}

			foundTxs = append(foundTxs, rt)
		}

		// add remaining ids
		remainingIDs = outerRightJoin(foundTxs, ids)
		if len(remainingIDs) > 0 || thErr != nil {
			f.logger.WarnContext(ctx, "couldn't find transactions in TransactionHandler", slog.Any("ids", remainingIDs), slog.Any("source-err", thErr))
		}

		ids = remainingIDs[:]
		remainingIDs = nil
	}

	// try to get remaining txs from the node
	if source.Has(validator.SourceNodes) && f.bitcoinClient != nil {
		var nErr error
		for _, id := range ids {
			var bitcoinGetRawTxSpan trace.Span
			ctx, bitcoinGetRawTxSpan = tracing.StartTracing(ctx, "Bitcoind_GetRawTxs", f.tracingEnabled, f.tracingAttributes...)
			nTx, err := f.bitcoinClient.GetRawTransaction(id)
			tracing.EndTracing(bitcoinGetRawTxSpan)
			if err != nil {
				nErr = errors.Join(nErr, fmt.Errorf("%s: %w", id, err))
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

		if len(remainingIDs) > 0 || nErr != nil {
			f.logger.WarnContext(ctx, "couldn't find transactions in node", slog.Any("ids", remainingIDs), slog.Any("source-error", nErr))
		}

		ids = remainingIDs[:]
		remainingIDs = nil
	}

	// at last try the WoC
	if source.Has(validator.SourceWoC) && len(ids) > 0 {
		var wocSpan trace.Span
		ctx, wocSpan = tracing.StartTracing(ctx, "WocClient_GetRawTxs", f.tracingEnabled, f.tracingAttributes...)
		wocTxs, wocErr := f.wocClient.GetRawTxs(ctx, ids)
		defer tracing.EndTracing(wocSpan)

		for _, wTx := range wocTxs {
			if wTx.Error != "" {
				wocErr = errors.Join(wocErr, fmt.Errorf("returned error data tx %s: %s", wTx.TxID, wTx.Error))
				continue
			}

			rt, e := newRawTx(wTx.TxID, wTx.Hex, wTx.BlockHeight)
			if e != nil {
				return nil, e
			}

			foundTxs = append(foundTxs, rt)
		}

		// add remaining ids
		remainingIDs = outerRightJoin(foundTxs, ids)
		if len(remainingIDs) > 0 || wocErr != nil {
			f.logger.WarnContext(ctx, "couldn't find transactions in WoC", slog.Any("ids", remainingIDs), slog.Any("source-error", wocErr))
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

func outerRightJoin(left []validator.RawTx, right []string) []string {
	var outerRight []string

	for _, id := range right {
		found := false
		for _, tx := range left {
			if tx.TxID == id {
				found = true
				break
			}
		}

		if !found {
			outerRight = append(outerRight, id)
		}
	}

	return outerRight
}
