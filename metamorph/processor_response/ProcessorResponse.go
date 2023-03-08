package processor_response

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/TAAL-GmbH/arc/metamorph/metamorph_api"
	"github.com/libsv/go-p2p"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/ordishs/go-utils"
	"github.com/ordishs/gocore"
	"github.com/sasha-s/go-deadlock"
)

type StatusAndError struct {
	Hash   *chainhash.Hash
	Status metamorph_api.Status
	Err    error
}

type ProcessorResponseStatusUpdate struct {
	Status      metamorph_api.Status
	Source      string
	StatusErr   error
	UpdateStore func() error
	Callback    func(err error)
}

type ProcessorResponse struct {
	callerCh              chan StatusAndError
	NoStats               bool
	statusUpdateCh        chan *ProcessorResponseStatusUpdate
	Hash                  *chainhash.Hash
	Start                 time.Time
	Retries               atomic.Uint32
	LastStatusUpdateNanos atomic.Int64
	// The following fields are protected by the mutex
	mu             deadlock.RWMutex
	Err            error
	AnnouncedPeers []p2p.PeerI
	Status         metamorph_api.Status
	Log            []ProcessorResponseLog
}

func NewProcessorResponse(hash *chainhash.Hash) *ProcessorResponse {
	return newProcessorResponse(hash, metamorph_api.Status_UNKNOWN, nil)
}

// NewProcessorResponseWithStatus creates a new ProcessorResponse with the given status.
// It is used when restoring the ProcessorResponseMap from the database.
func NewProcessorResponseWithStatus(hash *chainhash.Hash, status metamorph_api.Status) *ProcessorResponse {
	return newProcessorResponse(hash, status, nil)
}

func NewProcessorResponseWithChannel(hash *chainhash.Hash, ch chan StatusAndError) *ProcessorResponse {
	return newProcessorResponse(hash, metamorph_api.Status_UNKNOWN, ch)
}

func newProcessorResponse(hash *chainhash.Hash, status metamorph_api.Status, ch chan StatusAndError) *ProcessorResponse {
	pr := &ProcessorResponse{
		Start:          time.Now(),
		Hash:           hash,
		Status:         status,
		callerCh:       ch,
		statusUpdateCh: make(chan *ProcessorResponseStatusUpdate, 10),
		Log: []ProcessorResponseLog{
			{
				DeltaT: 0,
				Status: status.String(),
			},
		},
	}
	pr.LastStatusUpdateNanos.Store(pr.Start.UnixNano())

	go func() {
		for statusUpdate := range pr.statusUpdateCh {
			pr.updateStatus(statusUpdate)
		}
	}()

	return pr
}

func (r *ProcessorResponse) UpdateStatus(statusUpdate *ProcessorResponseStatusUpdate) {
	r.statusUpdateCh <- statusUpdate
}

func (r *ProcessorResponse) updateStatus(statusUpdate *ProcessorResponseStatusUpdate) {
	if statusUpdate.Status != metamorph_api.Status_MINED && r.Status >= statusUpdate.Status && statusUpdate.StatusErr == nil {
		r.addLog(statusUpdate.Status, statusUpdate.Source, "duplicate")
		return
	}

	if statusUpdate.UpdateStore != nil {
		if err := statusUpdate.UpdateStore(); err != nil {
			r.setErr(err, "processorResponse")
			statusUpdate.Callback(err)
			return
		}
	}

	if statusUpdate.StatusErr != nil {
		_ = r.setStatusAndError(statusUpdate.Status, statusUpdate.StatusErr, statusUpdate.Source)
	} else {
		_ = r.setStatus(statusUpdate.Status, statusUpdate.Source)
	}

	statKey := fmt.Sprintf("%d: %s", statusUpdate.Status, statusUpdate.Status.String())
	r.LastStatusUpdateNanos.Store(gocore.NewStat("processorResponse").NewStat(statKey).AddTime(r.LastStatusUpdateNanos.Load()))

	if !r.NoStats {
		statusUpdate.Callback(nil)
	}
}

func (r *ProcessorResponse) SetPeers(peers []p2p.PeerI) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.AnnouncedPeers = peers
}

func (r *ProcessorResponse) GetPeers() []p2p.PeerI {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.AnnouncedPeers
}

func (r *ProcessorResponse) String() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.Err != nil {
		return fmt.Sprintf("%v: %s [%s] %s", r.Hash, r.Start.Format(time.RFC3339), r.Status.String(), r.Err.Error())
	}
	return fmt.Sprintf("%v: %s [%s]", r.Hash, r.Start.Format(time.RFC3339), r.Status.String())
}

func (r *ProcessorResponse) setStatus(status metamorph_api.Status, source string) bool {
	r.Status = status

	sae := StatusAndError{
		Hash:   r.Hash,
		Status: r.Status,
	}

	ch := r.callerCh

	r.addLog(status, source, "")

	if ch != nil {
		return utils.SafeSend(ch, sae)
	}

	return true
}

func (r *ProcessorResponse) GetStatus() metamorph_api.Status {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Status
}

func (r *ProcessorResponse) setErr(err error, source string) bool {
	r.mu.Lock()

	r.Err = err

	sae := StatusAndError{
		Hash:   r.Hash,
		Status: r.Status,
		Err:    err,
	}

	ch := r.callerCh

	r.addLog(r.Status, source, err.Error())

	r.mu.Unlock()

	if ch != nil {
		return utils.SafeSend(ch, sae)
	}

	return true
}

func (r *ProcessorResponse) setStatusAndError(status metamorph_api.Status, err error, source string) bool {
	r.Status = status
	r.Err = err

	sae := StatusAndError{
		Hash:   r.Hash,
		Status: r.Status,
		Err:    err,
	}

	ch := r.callerCh

	r.addLog(status, source, err.Error())

	if ch != nil {
		return utils.SafeSend(ch, sae)
	}

	return true
}

func (r *ProcessorResponse) GetErr() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.Err
}

func (r *ProcessorResponse) GetRetries() uint32 {
	return r.Retries.Load()
}

func (r *ProcessorResponse) IncrementRetry() uint32 {
	r.Retries.Add(1)
	return r.Retries.Load()
}

func (r *ProcessorResponse) AddLog(status metamorph_api.Status, source string, info string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.addLog(status, source, info)
}

func (r *ProcessorResponse) addLog(status metamorph_api.Status, source string, info string) {
	r.Log = append(r.Log, ProcessorResponseLog{
		DeltaT: time.Since(r.Start).Nanoseconds(),
		Status: status.String(),
		Source: source,
		Info:   info,
	})
}
