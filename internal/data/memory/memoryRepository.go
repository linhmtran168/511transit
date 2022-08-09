package memory

import (
	"encoding/json"

	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/data"
	"github.com/linhmtran168/511transit/internal/models"
	transitapi "github.com/linhmtran168/511transit/internal/transit-api"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type MemoryRepository struct {
	apiClient     transitapi.TransitAPI
	operatorStore *operatorsData
	tripStore     *tripUpdatesData
}

func NewMemoryRepository(apiClient transitapi.TransitAPI) data.DataRepository {
	tripStore := &tripUpdatesData{
		tripMap: cmap.New[*tripEntry](),
	}
	return &MemoryRepository{
		apiClient:     apiClient,
		operatorStore: &operatorsData{},
		tripStore:     tripStore,
	}
}

func (p *MemoryRepository) GetOperators() ([]*models.Operator, error) {
	p.operatorStore.lock.Lock()
	defer p.operatorStore.lock.Unlock()
	if cachedOperators, ok := p.operatorStore.getCache(false); ok {
		log.Info().Msg("Operators cache exists, use data from cache")
		return cachedOperators, nil
	}

	response, err := p.apiClient.GetOperators()
	if err != nil {
		// we extended the cache in case of error
		if cachedOperators, ok := p.operatorStore.getCache(true); ok {
			log.Warn().Err(err).Msg("Failed to get operators from API, using stalled cache")
			return cachedOperators, nil
		}

		return nil, err
	}

	var operators []*models.Operator
	if err := json.Unmarshal(response, &operators); err != nil {
		return nil, err
	}
	p.operatorStore.updateCache(operators)

	return p.operatorStore.operators, nil
}

func (p *MemoryRepository) GetTripUpdates(operatorID string) ([]*gtfs.FeedEntity, error) {
	if operatorID == "" {
		return []*gtfs.FeedEntity{}, nil
	}
	if cachedTrips, ok := p.tripStore.getCache(operatorID, false); ok {
		log.Info().Msg("Trip update cache exists, use data from cache")
		return cachedTrips, nil
	}
	response, err := p.apiClient.GetTripUpdates(operatorID)
	if err != nil {
		// we extended the cache in case of error
		if cachedTrips, ok := p.tripStore.getCache(operatorID, true); ok {
			log.Warn().Err(err).Msg("Failed to get trips from API, using stalled cache")
			return cachedTrips, nil
		}

		return nil, err
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(response, &feed)
	if err != nil {
		return nil, err
	}

	if len(feed.Entity) > 0 {
		p.tripStore.updateCache(operatorID, feed.Entity)
	}

	return feed.Entity, nil
}
