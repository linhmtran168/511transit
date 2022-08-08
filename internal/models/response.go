package models

import (
	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
)

type OperatorsResponse struct {
	ResponseType string      `json:"responseType"`
	Data         []*Operator `json:"data"`
}

type TripUpdatesResponse struct {
	ResponseType string             `json:"responseType"`
	Data         []*gtfs.FeedEntity `json:"data"`
}
