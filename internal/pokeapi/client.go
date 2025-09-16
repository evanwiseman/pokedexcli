package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/evanwiseman/pokedexcli/internal/pokecache"
)

const (
	BaseURL         = "https://pokeapi.co/api/v2/"
	LocationAreaURL = BaseURL + "location-area/"
)

type Config struct {
	Next     *string
	Previous *string
}

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type LocationAreaList struct {
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
		return nil, fmt.Errorf("error get url %v: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}
	c.cache.Add(url, body)
	return body, nil
}

func (c *Client) GetLocationArea(name string) (*LocationArea, error) {
	url := LocationAreaURL + name
	bytes, err := c.FetchBytes(url)
	if err != nil {
		return nil, err
	}

	var area LocationArea
	err = json.Unmarshal(bytes, &area)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling bytes: %v", err)
	}
	return &area, nil
}

func (c *Client) GetLocationAreaList(url string) (*LocationAreaList, error) {
	bytes, err := c.FetchBytes(url)
	if err != nil {
		return nil, err
	}

	var areas LocationAreaList
	err = json.Unmarshal(bytes, &areas)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling bytes: %v", err)
	}
	return &areas, nil
}
