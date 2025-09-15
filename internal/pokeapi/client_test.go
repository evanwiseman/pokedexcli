package pokeapi

import (
	"testing"
)

func TestGetLocationAreas(t *testing.T) {
	url := LocationAreasURL
	locationAreas, err := GetLocationAreas(url)
	if err != nil {
		t.Fatalf("Unable to resolve %s", url)
	}

	if len(locationAreas.Results) != 20 {
		t.Fatalf("Expected 20 results got %d", len(locationAreas.Results))
	}
}
