package openweathermap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"qonto-observability/internal/domain"
)

const (
	baseURL = "https://api.openweathermap.org/data/2.5/weather"
)

type WResponse struct {
	Name string `json: "name"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func GetWeather(city *domain.City, apiKey string) (*WResponse, error){
	params := url.Values{}
	params.Add("q", fmt.Sprintf("%s,%s", city.Name, city.Country))
	params.Add("units", "metric")
	params.Add("appid", apiKey)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur: statut HTTP %d, réponse: %s", resp.StatusCode, string(body))
	}

	var weather WResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}