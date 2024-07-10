package handler

import (
	"context"

	"github.com/bitcoin-sv/arc/internal/beef"
	"github.com/bitcoin-sv/arc/pkg/blocktx"
)

type merkleVerifier struct {
	v blocktx.MerkleRootsVerifier
}

func (a merkleVerifier) Verify(ctx context.Context, request []beef.MerkleRootVerificationRequest) ([]uint64, error) {
	blocktxReq := convertToBlocktxMerkleVerRequest(request)
	return a.v.VerifyMerkleRoots(ctx, blocktxReq)
}

func convertToBlocktxMerkleVerRequest(mrReq []beef.MerkleRootVerificationRequest) []blocktx.MerkleRootVerificationRequest {
	merkleRoots := make([]blocktx.MerkleRootVerificationRequest, 0, len(mrReq))

	for _, mr := range mrReq {
		merkleRoots = append(merkleRoots, blocktx.MerkleRootVerificationRequest{
			BlockHeight: mr.BlockHeight,
			MerkleRoot:  mr.MerkleRoot,
		})
	}

	return merkleRoots
}
