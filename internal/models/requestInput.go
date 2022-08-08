package models

type RequestInput struct {
	RequestType string                 `json:"requestType"`
	Data        map[string]interface{} `json:"data"`
}
