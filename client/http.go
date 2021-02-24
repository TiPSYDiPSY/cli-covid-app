package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPCClient struct {
	country string
	client  *http.Client
}

type CountryData struct {
	Confirmed float64 `json:"Confirmed"`
	Deaths    float64 `json:"Deaths"`
	Recovered float64 `json:"Recovered"`
	Active    float64 `json:"Active"`
}

func NewHTTPClient(country string) HTTPCClient {
	return HTTPCClient{
		country: country,
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c HTTPCClient) GetData() ([]CountryData, error) {
	var countriesData []CountryData
	res, err := c.client.Get(fmt.Sprintf("https://api.covid19api.com/live/country/%s/status/confirmed/date/%s", c.country,
		time.Now().AddDate(0, -6, 0).Format("2006-1-02")))
	if err != nil {
		return nil, fmt.Errorf("unable to get data: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service currenty unavailable : %v", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("unable to read data: %v", err)
	}

	if err := json.Unmarshal(body, &countriesData); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %v", err)
	}

	return countriesData, nil
}
