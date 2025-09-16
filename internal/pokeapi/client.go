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
	PokemonURL      = BaseURL + "pokemon/"
)

type Config struct {
	Next     *string
	Previous *string
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

func (c *Client) GetLocationArea(name string) (*LocationArea, error) {
	fullURL := LocationAreaURL + name
	bytes, err := c.FetchBytes(fullURL)
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

type LocationAreaList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) GetLocationAreaList(fullURL string) (*LocationAreaList, error) {
	bytes, err := c.FetchBytes(fullURL)
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

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
		BackDefault  string `json:"back_default"`
		FrontShiny   string `json:"front_shiny"`
		BackShiny    string `json:"back_shiny"`
	} `json:"sprites"`
}

func (c *Client) GetPokemon(name string) (*Pokemon, error) {
	fullURL := PokemonURL + name
	bytes, err := c.FetchBytes(fullURL)
	if err != nil {
		return nil, err
	}

	var pokemon Pokemon
	err = json.Unmarshal(bytes, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes: %v", err)
	}

	return &pokemon, nil
}
