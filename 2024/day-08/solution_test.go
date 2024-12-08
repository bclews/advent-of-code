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
			expected:  14,
			mapWidth:  12,
			mapHeight: 12,
		},
		{
			name: "Single antenna, no antinodes",
			inputMap: []string{
				"............",
				"............",
				".....a......",
				"............",
				"............",
			},
			expected:  0,
			mapWidth:  12,
			mapHeight: 5,
		},
		{
			name: "Two antennas, valid antinodes",
			inputMap: []string{
				"............",
				"............",
				".....a......",
				"............",
				".....a......",
				"............",
			},
			expected:  1,
			mapWidth:  12,
			mapHeight: 6,
		},
		{
			name: "Antennas with different frequencies",
			inputMap: []string{
				"............",
				".....a......",
				"............",
				".....b......",
				"............",
			},
			expected:  0,
			mapWidth:  12,
			mapHeight: 5,
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
