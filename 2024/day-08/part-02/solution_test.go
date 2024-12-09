package main

import (
	"testing"
)

func TestFindAntinodes(t *testing.T) {
	tests := []struct {
		name      string
		inputMap  []string
		expected  int
		mapWidth  int
		mapHeight int
	}{
		{
			name: "Example input",
			inputMap: []string{
				"............",
				"........0...",
				".....0......",
				".......0....",
				"....0.......",
				"......A.....",
				"............",
				"............",
				"........A...",
				".........A..",
				"............",
				"............",
			},
			expected:  34,
			mapWidth:  12,
			mapHeight: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input map
			antennas := parseMap(tt.inputMap)

			// Calculate antinodes
			result := findAntinodes(antennas, tt.mapWidth, tt.mapHeight)

			// Verify the result
			if len(result) != tt.expected {
				t.Errorf("Expected %d antinodes, got %d", tt.expected, len(result))
			}
		})
	}
}
