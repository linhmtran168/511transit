package transitapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TransitAPI interface {
	GetOperators() ([]byte, error)
	GetTripUpdates(operatorID string) ([]byte, error)
}

type TransitAPIClient struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

func NewTransitAPI() TransitAPI {
	return &TransitAPIClient{
		client:  &http.Client{},
		baseUrl: os.Getenv("API_ENDPOINT"),
		apiKey:  os.Getenv("API_KEY"),
	}
}

func (t *TransitAPIClient) GetOperators() ([]byte, error) {
	url := fmt.Sprintf("%s/gtfsoperators?api_key=%s&format=json", t.baseUrl, t.apiKey)
	response, err := t.sendGetRequest(url)
	// Clean UTF-8 mark for json
	response = bytes.TrimPrefix(response, []byte("\xef\xbb\xbf"))
	return response, err
}

func (t *TransitAPIClient) GetTripUpdates(operatorID string) ([]byte, error) {
	url := fmt.Sprintf("%s/tripupdates?api_key=%s&agency=%s", t.baseUrl, t.apiKey, operatorID)
	return t.sendGetRequest(url)
}

func (t *TransitAPIClient) sendGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
