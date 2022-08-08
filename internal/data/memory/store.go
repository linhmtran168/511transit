package memory

import (
	"sync"
	"time"

	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/models"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const OPERATORS_CACHE_TIMEOUT = time.Hour
const TRIP_UPDATE_CACHE_TIMEOUT = time.Minute

type operatorsData struct {
	lastUpdated time.Time
	operators   []*models.Operator
	lock        sync.Mutex
}

func (store *operatorsData) getCache(isExtended bool) ([]*models.Operator, bool) {
	timeout := OPERATORS_CACHE_TIMEOUT
	if isExtended {
		timeout = OPERATORS_CACHE_TIMEOUT * 5
	}
	if len(store.operators) > 0 && store.lastUpdated.Add(timeout).After(time.Now()) {
		return store.operators, true
	}

	return nil, false
}

func (store *operatorsData) updateCache(operator []*models.Operator) {
	store.operators = operator
	store.lastUpdated = time.Now()
}

type tripUpdatesData struct {
	// Use a concurrent map for thread-safe accessing trip data
	tripMap cmap.ConcurrentMap[*tripEntry]
}

type tripEntry struct {
	lastUpdated time.Time
	trips       []*gtfs.FeedEntity
}

func (store *tripUpdatesData) getCache(operatorID string, isExtended bool) ([]*gtfs.FeedEntity, bool) {
	if operatorID == "" {
		return nil, false
	}

	timeout := TRIP_UPDATE_CACHE_TIMEOUT
	if isExtended {
		timeout = TRIP_UPDATE_CACHE_TIMEOUT * 5
	}
	if entry, ok := store.tripMap.Get(operatorID); ok && entry.lastUpdated.Add(timeout).After(time.Now()) {
		return entry.trips, true
	}

	return nil, false
}

func (store *tripUpdatesData) updateCache(operatorID string, trips []*gtfs.FeedEntity) {
	entry := &tripEntry{
		lastUpdated: time.Now(),
		trips:       trips,
	}
	store.tripMap.Set(operatorID, entry)
}
