package metamorph

import (
	"encoding/json"
	"errors"

	"github.com/libsv/go-p2p/chaincfg/chainhash"

	"github.com/bitcoin-sv/arc/internal/metamorph/store"
)

type StatusUpdateMap map[chainhash.Hash]store.UpdateStatus

const CacheStatusUpdateKey = "status-updates"

var (
	ErrFailedToSerialize   = errors.New("failed to serialize value")
	ErrFailedToDeserialize = errors.New("failed to deserialize value")
)

func (p *Processor) GetProcessorMapSize() int {
	return p.responseProcessor.getMapLen()
}

func (p *Processor) updateStatusMap(statusUpdate store.UpdateStatus) (StatusUpdateMap, error) {
	statusUpdatesMap := p.getStatusUpdateMap()
	foundStatusUpdate, found := statusUpdatesMap[statusUpdate.Hash]

	if !found || shouldUpdateStatus(statusUpdate, foundStatusUpdate) {
		if len(statusUpdate.CompetingTxs) > 0 {
			statusUpdate.CompetingTxs = mergeUnique(statusUpdate.CompetingTxs, foundStatusUpdate.CompetingTxs)
		}

		statusUpdatesMap[statusUpdate.Hash] = statusUpdate
	}

	err := p.setStatusUpdateMap(statusUpdatesMap)
	if err != nil {
		return statusUpdatesMap, err
	}

	return statusUpdatesMap, nil
}

func (p *Processor) setStatusUpdateMap(statusUpdatesMap StatusUpdateMap) error {
	bytes, err := serializeStatusMap(statusUpdatesMap)
	if err != nil {
		return err
	}

	err = p.cacheStore.Set(CacheStatusUpdateKey, bytes, processStatusUpdatesIntervalDefault)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) getStatusUpdateMap() StatusUpdateMap {
	existingMap, err := p.cacheStore.Get(CacheStatusUpdateKey)

	if err == nil {
		statusUpdatesMap, err := deserializeStatusMap(existingMap)
		if err == nil {
			return statusUpdatesMap
		}
	}

	// If the key doesn't exist or there was an error unmarshalling the value return new map
	return make(StatusUpdateMap)
}

func (p *Processor) clearStatusUpdateMap() error {
	if p.cacheStore == nil {
		return nil
	}

	return p.cacheStore.Del(CacheStatusUpdateKey)

}

func shouldUpdateStatus(new, found store.UpdateStatus) bool {
	if new.Status > found.Status {
		return true
	}

	if new.Status == found.Status && !unorderedEqual(new.CompetingTxs, found.CompetingTxs) {
		return true
	}

	return false
}

// unorderedEqual checks if two string slices contain
// the same elements, regardless of order
func unorderedEqual(sliceOne, sliceTwo []string) bool {
	if len(sliceOne) != len(sliceTwo) {
		return false
	}

	exists := make(map[string]bool)

	for _, value := range sliceOne {
		exists[value] = true
	}

	for _, value := range sliceTwo {
		if !exists[value] {
			return false
		}
	}

	return true
}

// mergeUnique merges two string arrays into one with unique values
func mergeUnique(arr1, arr2 []string) []string {
	valueSet := make(map[string]struct{})

	for _, value := range arr1 {
		valueSet[value] = struct{}{}
	}

	for _, value := range arr2 {
		valueSet[value] = struct{}{}
	}

	uniqueSlice := make([]string, 0, len(valueSet))
	for key := range valueSet {
		uniqueSlice = append(uniqueSlice, key)
	}

	return uniqueSlice
}

func serializeStatusMap(updateStatusMap StatusUpdateMap) ([]byte, error) { //nolint:unused
	serializeMap := make(map[string]store.UpdateStatus)
	for k, v := range updateStatusMap {
		serializeMap[k.String()] = v
	}

	bytes, err := json.Marshal(serializeMap)
	if err != nil {
		return nil, errors.Join(ErrFailedToSerialize, err)
	}
	return bytes, nil
}

func deserializeStatusMap(data []byte) (StatusUpdateMap, error) { //nolint:unused
	serializeMap := make(map[string]store.UpdateStatus)
	updateStatusMap := make(StatusUpdateMap)

	err := json.Unmarshal(data, &serializeMap)
	if err != nil {
		return nil, errors.Join(ErrFailedToDeserialize, err)
	}

	for k, v := range serializeMap {
		hash, err := chainhash.NewHashFromStr(k)
		if err != nil {
			return nil, errors.Join(ErrFailedToDeserialize, err)
		}
		updateStatusMap[*hash] = v
	}

	return updateStatusMap, nil
}
