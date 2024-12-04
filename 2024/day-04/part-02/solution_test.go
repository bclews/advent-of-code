package main

import (
	"testing"
)

// Test the first scenario with the provided 10x10 grid
func TestCountXmasGridsWithLargeGrid(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]rune
		expected int
	}{
		{
			name: "Example grid with multiple occurrences",
			grid: [][]rune{
				{'M', 'M', 'M', 'S', 'X', 'X', 'M', 'A', 'S', 'M'},
				{'M', 'S', 'A', 'M', 'X', 'M', 'S', 'M', 'S', 'A'},
				{'A', 'M', 'X', 'S', 'X', 'M', 'A', 'A', 'M', 'M'},
				{'M', 'S', 'A', 'M', 'A', 'S', 'M', 'S', 'M', 'X'},
				{'X', 'M', 'A', 'S', 'A', 'M', 'X', 'A', 'M', 'M'},
				{'X', 'X', 'A', 'M', 'M', 'X', 'X', 'A', 'M', 'A'},
				{'S', 'M', 'S', 'M', 'S', 'A', 'S', 'X', 'S', 'S'},
				{'S', 'A', 'X', 'A', 'M', 'A', 'S', 'A', 'A', 'A'},
				{'M', 'A', 'M', 'M', 'M', 'X', 'M', 'M', 'M', 'M'},
				{'M', 'X', 'M', 'X', 'A', 'X', 'M', 'A', 'S', 'X'},
			},
			expected: 9,
		},
		{
			name: "Pattern 1: M . S / . A . / M . S",
			grid: [][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'M', '.', 'S', 'X', 'X', 'X'},
				{'X', '.', 'A', '.', 'X', 'X', 'X'},
				{'X', 'M', '.', 'S', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			expected: 1,
		},
		{
			name: "Pattern 2: M . M / . A . / S . S",
			grid: [][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'M', '.', 'M', 'X', 'X', 'X'},
				{'X', '.', 'A', '.', 'X', 'X', 'X'},
				{'X', 'S', '.', 'S', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			expected: 1,
		},
		{
			name: "Pattern 3: S . S / . A . / M . M",
			grid: [][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'S', '.', 'S', 'X', 'X', 'X'},
				{'X', '.', 'A', '.', 'X', 'X', 'X'},
				{'X', 'M', '.', 'M', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			expected: 1,
		},
		{
			name: "Pattern 4: S . M / . A . / S . M",
			grid: [][]rune{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'S', '.', 'M', 'X', 'X', 'X'},
				{'X', '.', 'A', '.', 'X', 'X', 'X'},
				{'X', 'S', '.', 'M', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := countXmasGrids(tc.grid)
			if result != tc.expected {
				t.Errorf("Expected %d occurrences, got %d", tc.expected, result)
			}
		})
	}
}
