package pokeapi

import (
	"strings"
	"testing"
)

func TestFetchBytes(t *testing.T) {
	client := NewClient()

	bytes, err := client.FetchBytes("https://www.example.com")
	if err != nil {
		t.Fatalf("FetchBytes returned error: %v", err)
	}
	if len(bytes) == 0 {
		t.Fatalf("FetchBytes returned empty body")
	}
	if !strings.HasPrefix(strings.ToLower(string(bytes)), "<!doctype html>") {
		t.Fatalf("FetchBytes returned unexpected content: %s", string(bytes)[:50])
	}
}

func TestGetLocationArea(t *testing.T) {
	client := NewClient()

	area, err := client.GetLocationArea("canalave-city-area")
	if err != nil {
		t.Fatalf("GetLocationArea returned error: %v", err)
	}
	if area.Name != "canalave-city-area" {
		t.Errorf("expected name 'canalave-city-area', got %s", area.Name)
	}
	if len(area.PokemonEncounters) == 0 {
		t.Errorf("expected at least 1 Pokémon encounter, got 0")
	}
}

func TestGetLocationAreaList(t *testing.T) {
	client := NewClient()

	areas, err := client.GetLocationAreaList(LocationAreaURL)
	if err != nil {
		t.Fatalf("GetLocationAreaList returned error: %v", err)
	}
	if areas.Count == 0 {
		t.Errorf("expected non-zero count of areas")
	}
	if len(areas.Results) == 0 {
		t.Errorf("expected some results in the list")
	}
	if areas.Next == nil {
		t.Errorf("expected a next page URL")
	}
}
