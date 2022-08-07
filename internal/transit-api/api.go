package transitapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TransitAPI struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

func NewTransitAPI() *TransitAPI {
	return &TransitAPI{
		client:  &http.Client{},
		baseUrl: os.Getenv("API_ENDPOINT"),
		apiKey:  os.Getenv("API_KEY"),
	}
}

func (t *TransitAPI) GetOperators() ([]byte, error) {
	url := fmt.Sprintf("%s?api_key=%s", t.baseUrl, t.apiKey)
	return t.sendGetRequest(url)
}

func (t *TransitAPI) GetTripUpdates(operatorID string) ([]byte, error) {
	url := fmt.Sprintf("%s?api_key=%s&agency=%s", t.baseUrl, t.apiKey, operatorID)
	return t.sendGetRequest(url)
}

func (t *TransitAPI) sendGetRequest(url string) ([]byte, error) {
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
