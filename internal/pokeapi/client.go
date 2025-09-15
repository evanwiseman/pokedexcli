package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	BaseURL          = "https://pokeapi.co/api/v2/"
	LocationAreasURL = BaseURL + "location-area/"
)

type Config struct {
	Next     *string
	Previous *string
}

type LocationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(url string) (LocationAreas, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationAreas{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreas{}, err
	}

	var locationAreas LocationAreas
	err = json.Unmarshal(body, &locationAreas)
	if err != nil {
		return LocationAreas{}, err
	}
	return locationAreas, nil
}
