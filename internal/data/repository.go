package data

import (
	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
	"github.com/linhmtran168/511transit/internal/models"
)

type DataRepository interface {
	GetOperators() ([]*models.Operator, error)
	GetTripUpdates(id string) ([]*gtfs.FeedEntity, error)
}
