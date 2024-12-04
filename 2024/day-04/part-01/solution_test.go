package main

import (
	"testing"
)

func TestFindWordOccurrences(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]rune
		word     string
		expected int
	}{
		{
			name: "Basic grid with multiple occurrences",
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
			word:     "XMAS",
			expected: 18,
		},
		{
			name: "Horizontal Right",
			grid: [][]rune{
				{'X', 'M', 'A', 'S', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Horizontal Left",
			grid: [][]rune{
				{'S', 'A', 'M', 'X', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Vertical Down",
			grid: [][]rune{
				{'X', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'M', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'A', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'S', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Vertical Up",
			grid: [][]rune{
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'S', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'A', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'M', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'X', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Diagonal Down-Right",
			grid: [][]rune{
				{'X', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'M', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'A', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'S', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Diagonal Down-Left",
			grid: [][]rune{
				{'O', 'O', 'O', 'O', 'O', 'O', 'X'},
				{'O', 'O', 'O', 'O', 'O', 'M', 'O'},
				{'O', 'O', 'O', 'O', 'A', 'O', 'O'},
				{'O', 'O', 'O', 'S', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Diagonal Up-Right",
			grid: [][]rune{
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'S', 'O', 'O', 'O'},
				{'O', 'O', 'A', 'O', 'O', 'O', 'O'},
				{'O', 'M', 'O', 'O', 'O', 'O', 'O'},
				{'X', 'O', 'O', 'O', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
		{
			name: "Diagonal Up-Left",
			grid: [][]rune{
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'O'},
				{'O', 'O', 'O', 'O', 'O', 'O', 'X'},
				{'O', 'O', 'O', 'O', 'O', 'M', 'O'},
				{'O', 'O', 'O', 'O', 'A', 'O', 'O'},
				{'O', 'O', 'O', 'S', 'O', 'O', 'O'},
			},
			word:     "XMAS",
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FindWordOccurrences(tc.grid, tc.word)
			if result != tc.expected {
				t.Errorf("Expected %d occurrences, got %d", tc.expected, result)
			}
		})
	}
}
