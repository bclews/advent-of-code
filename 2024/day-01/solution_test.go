package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Test case with valid input
	input := `3   4
4   3
2   5
1   3
3   9
3   3`

	column1, column2, err := parseFile(strings.NewReader(input))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedCol1 := []int{3, 4, 2, 1, 3, 3}
	expectedCol2 := []int{4, 3, 5, 3, 9, 3}

	if !reflect.DeepEqual(column1, expectedCol1) {
		t.Errorf("Column 1 mismatch. Got %v, want %v", column1, expectedCol1)
	}

	if !reflect.DeepEqual(column2, expectedCol2) {
		t.Errorf("Column 2 mismatch. Got %v, want %v", column2, expectedCol2)
	}
}

func TestParseFileInvalidFormat(t *testing.T) {
	// Test case with invalid input
	input := `3 4 extra
invalid`

	_, _, err := parseFile(strings.NewReader(input))

	if err == nil {
		t.Error("Expected an error, but got none")
	}
}

func TestSortColumns(t *testing.T) {
	// Test sorting functionality
	input1 := []int{3, 4, 2, 1, 3, 3}
	input2 := []int{4, 3, 5, 3, 9, 3}

	sorted1, sorted2 := sortColumns(input1, input2)

	expectedSorted1 := []int{1, 2, 3, 3, 3, 4}
	expectedSorted2 := []int{3, 3, 3, 4, 5, 9}

	if !reflect.DeepEqual(sorted1, expectedSorted1) {
		t.Errorf("Sorted Column 1 mismatch. Got %v, want %v", sorted1, expectedSorted1)
	}

	if !reflect.DeepEqual(sorted2, expectedSorted2) {
		t.Errorf("Sorted Column 2 mismatch. Got %v, want %v", sorted2, expectedSorted2)
	}
}

func TestCalculatePairedDistance(t *testing.T) {
	// Test case scenarios
	testCases := []struct {
		name          string
		col1          []int
		col2          []int
		expectedTotal int
		expectedPairs []Pair
	}{
		{
			name:          "Standard Case",
			col1:          []int{1, 2, 3, 3, 3, 4},
			col2:          []int{3, 3, 3, 4, 5, 9},
			expectedTotal: 11,
			expectedPairs: []Pair{
				{Left: 1, Right: 3, Distance: 2},
				{Left: 2, Right: 3, Distance: 1},
				{Left: 3, Right: 3, Distance: 0},
				{Left: 3, Right: 4, Distance: 1},
				{Left: 3, Right: 5, Distance: 2},
				{Left: 4, Right: 9, Distance: 5},
			},
		},
		{
			name:          "Mismatched Lengths",
			col1:          []int{1, 2},
			col2:          []int{3, 3, 4},
			expectedTotal: 0,
			expectedPairs: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totalDistance, pairs := calculatePairedDistance(tc.col1, tc.col2)

			if tc.name == "Mismatched Lengths" && totalDistance != 0 {
				t.Errorf("Expected 0 total distance for mismatched lengths")
				return
			}

			if tc.name == "Standard Case" {
				if totalDistance != tc.expectedTotal {
					t.Errorf("Expected total distance %d, got %d", tc.expectedTotal, totalDistance)
				}

				// Compare pairs
				if !reflect.DeepEqual(pairs, tc.expectedPairs) {
					t.Errorf("Pairs mismatch. Got %v, want %v", pairs, tc.expectedPairs)
				}
			}
		})
	}
}

func TestCalculateSimilarityScore(t *testing.T) {
	testCases := []struct {
		name          string
		col1          []int
		col2          []int
		expectedScore int
	}{
		{
			name:          "Standard Case",
			col1:          []int{3, 4, 2, 1, 3, 3},
			col2:          []int{4, 3, 5, 3, 9, 3},
			expectedScore: 31, // Manually calculated
		},
		{
			name:          "No Matches",
			col1:          []int{1, 2, 3},
			col2:          []int{4, 5, 6},
			expectedScore: 0,
		},
		{
			name:          "All Matches",
			col1:          []int{1, 1, 1},
			col2:          []int{1, 1, 1},
			expectedScore: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			similarityScore := calculateSimilarityScore(tc.col1, tc.col2)

			if similarityScore != tc.expectedScore {
				t.Errorf("Expected similarity score %d, got %d", tc.expectedScore, similarityScore)
			}
		})
	}
}
