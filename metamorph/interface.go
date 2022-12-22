package metamorph

import "github.com/TAAL-GmbH/arc/metamorph/metamorph_api"

type ProcessorI interface {
	LoadUnseen()
	ProcessTransaction(req *ProcessorRequest)
	SendStatusForTransaction(hashStr string, status metamorph_api.Status, err error) bool
	GetStats() *ProcessorStats
}
