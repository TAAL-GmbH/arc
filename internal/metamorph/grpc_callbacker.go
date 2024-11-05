package metamorph

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"

	"github.com/bitcoin-sv/arc/internal/callbacker/callbacker_api"
	"github.com/bitcoin-sv/arc/internal/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/internal/metamorph/store"
	"github.com/bitcoin-sv/arc/internal/tracing"
)

type GrpcCallbacker struct {
	cc                callbacker_api.CallbackerAPIClient
	l                 *slog.Logger
	tracingEnabled    bool
	tracingAttributes []attribute.KeyValue
}

func WithCallbackerTracer(attr ...attribute.KeyValue) func(*GrpcCallbacker) {
	return func(p *GrpcCallbacker) {
		p.tracingEnabled = true
		if len(attr) > 0 {
			p.tracingAttributes = append(p.tracingAttributes, attr...)
		}
	}
}

type CallbackerOption func(s *GrpcCallbacker)

func NewGrpcCallbacker(api callbacker_api.CallbackerAPIClient, logger *slog.Logger, opts ...CallbackerOption) GrpcCallbacker {
	c := GrpcCallbacker{
		cc: api,
		l:  logger,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

func (c GrpcCallbacker) SendCallback(ctx context.Context, data *store.Data) {
	ctx, span := tracing.StartTracing(ctx, "SendCallback", c.tracingEnabled, c.tracingAttributes...)
	defer tracing.EndTracing(span)

	if len(data.Callbacks) == 0 {
		return
	}

	in := toGrpcInput(data)
	if in == nil {
		return
	}

	_, err := c.cc.SendCallback(ctx, in)
	if err != nil {
		c.l.Error("sending callback failed", slog.String("err", err.Error()), slog.Any("input", in))
	}
}

func toGrpcInput(d *store.Data) *callbacker_api.SendCallbackRequest {
	routings := make([]*callbacker_api.CallbackRouting, 0, len(d.Callbacks))

	for _, c := range d.Callbacks {
		if c.CallbackURL != "" {
			routings = append(routings, &callbacker_api.CallbackRouting{
				Url:        c.CallbackURL,
				Token:      c.CallbackToken,
				AllowBatch: c.AllowBatch,
			})
		}
	}

	if len(routings) == 0 {
		return nil
	}

	in := callbacker_api.SendCallbackRequest{
		CallbackRoutings: routings,

		Txid:         d.Hash.String(),
		Status:       callbacker_api.Status(d.Status),
		MerklePath:   d.MerklePath,
		ExtraInfo:    getCallbackExtraInfo(d),
		CompetingTxs: getCallbackCompetitingTxs(d),

		BlockHash:   getCallbackBlockHash(d),
		BlockHeight: d.BlockHeight,
	}

	return &in
}

func getCallbackExtraInfo(d *store.Data) string {
	if d.Status == metamorph_api.Status_MINED && len(d.CompetingTxs) > 0 {
		return minedDoubleSpendMsg
	}

	return d.RejectReason
}

func getCallbackCompetitingTxs(d *store.Data) []string {
	if d.Status == metamorph_api.Status_MINED {
		return nil
	}

	return d.CompetingTxs
}

func getCallbackBlockHash(d *store.Data) string {
	if d.BlockHash == nil {
		return ""
	}

	return d.BlockHash.String()
}
