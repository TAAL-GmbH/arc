package blocktx

import (
	"context"
	"log/slog"
	"time"

	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/internal/blocktx/store"
	"github.com/libsv/go-p2p"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bitcoin-sv/arc/internal/grpc_opts"
)

// Server type carries the logger within it.
type Server struct {
	blocktx_api.UnsafeBlockTxAPIServer
	grpc_opts.GrpcServer

	logger                        *slog.Logger
	pm                            p2p.PeerManagerI
	store                         store.BlocktxStore
	maxAllowedBlockHeightMismatch int
}

// NewServer will return a server instance with the logger stored within it.
func NewServer(prometheusEndpoint string, maxMsgSize int, logger *slog.Logger,
	store store.BlocktxStore, pm p2p.PeerManagerI, maxAllowedBlockHeightMismatch int, tracingEnabled bool) (*Server, error) {

	logger = logger.With(slog.String("module", "server"))

	grpcServer, err := grpc_opts.NewGrpcServer(logger, "blocktx", prometheusEndpoint, maxMsgSize, tracingEnabled)
	if err != nil {
		return nil, err
	}

	s := &Server{
		GrpcServer:                    grpcServer,
		store:                         store,
		logger:                        logger,
		pm:                            pm,
		maxAllowedBlockHeightMismatch: maxAllowedBlockHeightMismatch,
	}

	blocktx_api.RegisterBlockTxAPIServer(s.GrpcServer.Srv, s)
	reflection.Register(s.GrpcServer.Srv)

	return s, nil
}

func (s *Server) Health(_ context.Context, _ *emptypb.Empty) (*blocktx_api.HealthResponse, error) {
	return &blocktx_api.HealthResponse{
		Ok:        true,
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

func (s *Server) ClearTransactions(ctx context.Context, clearData *blocktx_api.ClearData) (*blocktx_api.RowsAffectedResponse, error) {
	return s.store.ClearBlocktxTable(ctx, clearData.GetRetentionDays(), "transactions")
}

func (s *Server) ClearBlocks(ctx context.Context, clearData *blocktx_api.ClearData) (*blocktx_api.RowsAffectedResponse, error) {
	return s.store.ClearBlocktxTable(ctx, clearData.GetRetentionDays(), "blocks")
}

func (s *Server) ClearBlockTransactionsMap(ctx context.Context, clearData *blocktx_api.ClearData) (*blocktx_api.RowsAffectedResponse, error) {
	return s.store.ClearBlocktxTable(ctx, clearData.GetRetentionDays(), "block_transactions_map")
}

func (s *Server) DelUnfinishedBlockProcessing(ctx context.Context, req *blocktx_api.DelUnfinishedBlockProcessingRequest) (*blocktx_api.RowsAffectedResponse, error) {
	bhs, err := s.store.GetBlockHashesProcessingInProgress(ctx, req.GetProcessedBy())
	if err != nil {
		return &blocktx_api.RowsAffectedResponse{}, err
	}

	var rowsTotal int64
	for _, bh := range bhs {
		rows, err := s.store.DelBlockProcessing(ctx, bh, req.GetProcessedBy())
		if err != nil {
			return &blocktx_api.RowsAffectedResponse{}, err
		}

		rowsTotal += rows
	}

	return &blocktx_api.RowsAffectedResponse{Rows: rowsTotal}, nil
}

func (s *Server) VerifyMerkleRoots(ctx context.Context, req *blocktx_api.MerkleRootsVerificationRequest) (*blocktx_api.MerkleRootVerificationResponse, error) {
	return s.store.VerifyMerkleRoots(ctx, req.GetMerkleRoots(), s.maxAllowedBlockHeightMismatch)
}
