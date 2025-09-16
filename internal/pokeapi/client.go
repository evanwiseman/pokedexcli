package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/evanwiseman/pokedexcli/internal/pokecache"
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

type Client struct {
	baseURL    string
	httpClient *http.Client
	cache      *pokecache.Cache
}

func NewClient() *Client {
	return &Client{
		baseURL:    BaseURL,
		httpClient: &http.Client{},
		cache:      pokecache.NewCache(5 * time.Second),
	}
}

func (c *Client) FetchBytes(url string) ([]byte, error) {
	bytes, ok := c.cache.Get(url)
	if ok {
		return bytes, nil
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	c.cache.Add(url, body)
	return body, nil
}

func (c *Client) GetLocationAreas(url string) (*LocationAreas, error) {
	bytes, err := c.FetchBytes(url)
	if err != nil {
		return nil, err
	}

	var locationAreas LocationAreas
	err = json.Unmarshal(bytes, &locationAreas)
	if err != nil {
		return nil, err
	}
	return &locationAreas, nil
}
