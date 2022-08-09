package models

import (
	gtfs "github.com/linhmtran168/511transit/api/protos/gtfs-realtime"
)

type OperatorsResponse struct {
	ResponseType string      `json:"responseType"`
	Data         []*Operator `json:"data"`
}

type TripUpdatesResponse struct {
	OperatorID   string         `json:"operatorId"`
	ResponseType string         `json:"responseType"`
	Data         TripUpdateData `json:"data"`
}

type TripUpdateData struct {
	OperatorID  string             `json:"operatorId"`
	TripUpdates []*gtfs.FeedEntity `json:"tripUpdates"`
}
