package memory

import (
	"sync"
	"time"

	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/models"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type operatorsData struct {
	lastUpdated time.Time
	operators   []*models.Operator
	lock        sync.Mutex
}

type tripUpdatesData struct {
	// Use a concurrent map for thread-safe accessing trip data
	tripMap cmap.ConcurrentMap[*tripEntry]
}

type tripEntry struct {
	lastUpdated time.Time
	trips       []*gtfs.FeedEntity
}

var (
	// We only need one instance of memory cache for operators and trip updates
	// in memory to share between multiple goroutines.
	operatorsCache   *operatorsData
	tripUpdatesCache *tripUpdatesData
	operatorOnce     sync.Once
	tripUpdatesOnce  sync.Once
)

func newOperatorsData() *operatorsData {
	operatorOnce.Do(func() {
		operatorsCache = &operatorsData{}
	})
	return operatorsCache
}

func (store *operatorsData) GetCache(isExtended bool) ([]*models.Operator, bool) {
	timeout := OPERATORS_CACHE_TIMEOUT
	if isExtended {
		timeout = OPERATORS_CACHE_TIMEOUT * 5
	}
	if len(store.operators) > 0 && store.lastUpdated.Add(timeout).After(time.Now()) {
		return store.operators, true
	}

	return nil, false
}

func (store *operatorsData) UpdateCache(operator []*models.Operator) {
	store.operators = operator
	store.lastUpdated = time.Now()
}

func (store *operatorsData) Lock() {
	store.lock.Lock()
}

func (store *operatorsData) Unlock() {
	store.lock.Unlock()
}

func newTripUpdatesData() *tripUpdatesData {
	tripUpdatesOnce.Do(func() {
		tripUpdatesCache = &tripUpdatesData{
			tripMap: cmap.New[*tripEntry](),
		}
	})

	return tripUpdatesCache
}

func (store *tripUpdatesData) GetCache(operatorID string, isExtended bool) ([]*gtfs.FeedEntity, bool) {
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

func (store *tripUpdatesData) UpdateCache(operatorID string, trips []*gtfs.FeedEntity) {
	entry := &tripEntry{
		lastUpdated: time.Now(),
		trips:       trips,
	}
	store.tripMap.Set(operatorID, entry)
}
