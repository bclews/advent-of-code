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
			name: "Single trailhead with score 1",
			grid: [][]int{
				{0, 1, 2, 3},
				{1, 2, 3, 4},
				{8, 7, 6, 5},
				{9, 8, 7, 6},
			},
			expected: 1,
		},
		{
			name: "Trailhead with score 2",
			grid: [][]int{
				{'.', '.', '.', 0, '.', '.', '.'},
				{'.', '.', '.', 1, '.', '.', '.'},
				{'.', '.', '.', 2, '.', '.', '.'},
				{6, 5, 4, 3, 4, 5, 6},
				{7, '.', '.', '.', '.', '.', 7},
				{8, '.', '.', '.', '.', '.', 8},
				{9, '.', '.', '.', '.', '.', 9},
			},
			expected: 2,
		},
		{
			name: "Trailhead with 4 reachable 9s",
			grid: [][]int{
				{'.', '.', 9, 0, '.', '.', 9},
				{'.', '.', '.', 1, '.', 9, 8},
				{'.', '.', '.', 2, '.', '.', 7},
				{6, 5, 4, 3, 4, 5, 6},
				{7, 6, 5, '.', 9, 8, 7},
				{8, 7, 6, '.', '.', '.', '.'},
				{9, 8, 7, '.', '.', '.', '.'},
			},
			expected: 4,
		},
		{
			name: "Two trailheads with scores 1 and 2",
			grid: [][]int{
				{1, 0, '.', '.', 9, '.', '.'},
				{2, '.', '.', '.', 8, '.', '.'},
				{3, '.', '.', '.', 7, '.', '.'},
				{4, 5, 6, 7, 6, 5, 4},
				{'.', '.', '.', 8, '.', '.', 3},
				{'.', '.', '.', 9, '.', '.', 2},
				{'.', '.', '.', '.', '.', 0, 1},
			},
			expected: 3,
		},
		{
			name: "Large example with sum of scores 36",
			grid: [][]int{
				{8, 9, 0, 1, 0, 1, 2, 3},
				{7, 8, 1, 2, 1, 2, 1, 4},
				{8, 7, 4, 3, 0, 9, 6, 5},
				{9, 6, 5, 4, 9, 8, 7, 4},
				{4, 5, 6, 7, 8, 9, 0, 3},
				{3, 2, 0, 1, 9, 0, 1, 2},
				{0, 1, 3, 2, 9, 8, 0, 1},
				{1, 0, 4, 5, 6, 7, 3, 2},
			},
			expected: 36,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CalculateTrailheadScores(test.grid)
			if result != test.expected {
				t.Errorf("CalculateTrailheadScores() = %d; want %d", result, test.expected)
			}
		})
	}
}
