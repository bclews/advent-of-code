package main

import (
	"testing"
)

func TestCalculateTrailheadScores(t *testing.T) {
	tests := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name: "Single trailhead with rating 3",
			grid: [][]int{
				{'.', '.', '.', '.', '.', 0, '.'},
				{'.', '.', 4, 3, 2, 1, '.'},
				{'.', '.', 5, '.', '.', 2, '.'},
				{'.', '.', 6, 5, 4, 3, '.'},
				{'.', '.', 7, '.', '.', 4, '.'},
				{'.', '.', 8, 7, 6, 5, '.'},
				{'.', '.', 9, '.', '.', '.', '.'},
			},
			expected: 3,
		},
		{
			name: "Single trailhead with rating 13",
			grid: [][]int{
				{'.', '.', 9, 0, '.', '.', 9},
				{'.', '.', '.', 1, '.', 9, 8},
				{'.', '.', '.', 2, '.', '.', 7},
				{6, 5, 4, 3, 4, 5, 6},
				{7, 6, 5, '.', 9, 8, 7},
				{8, 7, 6, '.', '.', '.', '.'},
				{9, 8, 7, '.', '.', '.', '.'},
			},
			expected: 13,
		},
		{
			name: "Single trailhead with rating 227",
			grid: [][]int{
				{0, 1, 2, 3, 4, 5},
				{1, 2, 3, 4, 5, 6},
				{2, 3, 4, 5, 6, 7},
				{3, 4, 5, 6, 7, 8},
				{4, '.', 6, 7, 8, 9},
				{5, 6, 7, 8, 9, '.'},
			},
			expected: 227,
		},
		{
			name: "Larger example with sum of trailhead ratings 81",
			grid: [][]int{
				{8, 9, 0, 1, 0, 1, 2, 3},
				{7, 8, 1, 2, 1, 8, 7, 4},
				{8, 7, 4, 3, 0, 9, 6, 5},
				{9, 6, 5, 4, 9, 8, 7, 4},
				{4, 5, 6, 7, 8, 9, 0, 3},
				{3, 2, 0, 1, 9, 0, 1, 2},
				{0, 1, 3, 2, 9, 8, 0, 1},
				{1, 0, 4, 5, 6, 7, 3, 2},
			},
			expected: 81,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTrailheadScores(tt.grid)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
