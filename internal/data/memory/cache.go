//go:generate mockgen -source=./cache.go -destination=../../mock/cache_mock.go -package=mock
package memory

import (
	"time"

	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/models"
)

const OPERATORS_CACHE_TIMEOUT = time.Hour
const TRIP_UPDATE_CACHE_TIMEOUT = time.Minute

type OperatorCache interface {
	GetCache(isExtended bool) ([]*models.Operator, bool)
	UpdateCache(operator []*models.Operator)
	Lock()
	Unlock()
}

type TripUpdateCache interface {
	GetCache(operatorID string, isExtended bool) ([]*gtfs.FeedEntity, bool)
	UpdateCache(operatorID string, trips []*gtfs.FeedEntity)
}
