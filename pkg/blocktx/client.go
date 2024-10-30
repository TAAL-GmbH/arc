package blocktx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/bitcoin-sv/arc/internal/blocktx/blocktx_api"
	"github.com/bitcoin-sv/arc/internal/grpc_opts"
)

// check if Client implements all necessary interfaces
var _ BlocktxClient = &Client{}
var _ MerkleRootsVerifier = &Client{}

type BlocktxClient interface {
	Health(ctx context.Context) error
	ClearTransactions(ctx context.Context, retentionDays int32) (int64, error)
	ClearBlocks(ctx context.Context, retentionDays int32) (int64, error)
	ClearBlockTransactionsMap(ctx context.Context, retentionDays int32) (int64, error)
	DelUnfinishedBlockProcessing(ctx context.Context, processedBy string) (int64, error)
}

// MerkleRootsVerifier verifies the merkle roots existance in blocktx db and returns unverified block heights.
type MerkleRootsVerifier interface {
	// VerifyMerkleRoots verifies the merkle roots existance in blocktx db and returns unverified block heights.
	VerifyMerkleRoots(ctx context.Context, merkleRootVerificationRequest []MerkleRootVerificationRequest) ([]uint64, error)
}

type MerkleRootVerificationRequest struct {
	MerkleRoot  string
	BlockHeight uint64
}

type Client struct {
	client blocktx_api.BlockTxAPIClient
}

func NewClient(client blocktx_api.BlockTxAPIClient) *Client {
	btc := &Client{
		client: client,
	}

	return btc
}

func (btc *Client) Health(ctx context.Context) error {
	_, err := btc.client.Health(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (btc *Client) DelUnfinishedBlockProcessing(ctx context.Context, processedBy string) (int64, error) {
	resp, err := btc.client.DelUnfinishedBlockProcessing(ctx, &blocktx_api.DelUnfinishedBlockProcessingRequest{ProcessedBy: processedBy})
	if err != nil {
		return 0, err
	}
	return resp.Rows, nil
}

func (btc *Client) ClearTransactions(ctx context.Context, retentionDays int32) (int64, error) {
	resp, err := btc.client.ClearTransactions(ctx, &blocktx_api.ClearData{RetentionDays: retentionDays})
	if err != nil {
		return 0, err
	}
	return resp.Rows, nil
}

func (btc *Client) ClearBlocks(ctx context.Context, retentionDays int32) (int64, error) {
	resp, err := btc.client.ClearBlocks(ctx, &blocktx_api.ClearData{RetentionDays: retentionDays})
	if err != nil {
		return 0, err
	}
	return resp.Rows, nil
}

func (btc *Client) ClearBlockTransactionsMap(ctx context.Context, retentionDays int32) (int64, error) {
	resp, err := btc.client.ClearBlockTransactionsMap(ctx, &blocktx_api.ClearData{RetentionDays: retentionDays})
	if err != nil {
		return 0, err
	}
	return resp.Rows, nil
}

func (btc *Client) VerifyMerkleRoots(ctx context.Context, merkleRootVerificationRequest []MerkleRootVerificationRequest) ([]uint64, error) {
	merkleRoots := make([]*blocktx_api.MerkleRootVerificationRequest, 0)

	for _, mr := range merkleRootVerificationRequest {
		merkleRoots = append(merkleRoots, &blocktx_api.MerkleRootVerificationRequest{MerkleRoot: mr.MerkleRoot, BlockHeight: mr.BlockHeight})
	}

	resp, err := btc.client.VerifyMerkleRoots(ctx, &blocktx_api.MerkleRootsVerificationRequest{MerkleRoots: merkleRoots})
	if err != nil {
		return nil, err
	}

	return resp.UnverifiedBlockHeights, nil
}

func DialGRPC(address string, prometheusEndpoint string, grpcMessageSize int, tracingEnabled bool) (*grpc.ClientConn, error) {
	dialOpts, err := grpc_opts.GetGRPCClientOpts(prometheusEndpoint, grpcMessageSize, tracingEnabled)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(address, dialOpts...)
}
