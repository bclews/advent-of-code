package main

import (
	"strings"
	"testing"
)

// Unit tests
func TestParseInputFile(t *testing.T) {
	// Create a test input file
	input := `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`

	// Parse the file
	result, err := parseFile(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected result
	expected := [][]int{
		{7, 6, 4, 2, 1},
		{1, 2, 7, 8, 9},
		{9, 7, 6, 2, 1},
		{1, 3, 2, 4, 5},
		{8, 6, 4, 4, 1},
		{1, 3, 6, 7, 9},
	}

	// Validate the result
	if len(result) != len(expected) {
		t.Fatalf("Expected %d lines, got %d", len(expected), len(result))
	}

	// Compare each line
	for i := range expected {
		if len(result[i]) != len(expected[i]) {
			t.Fatalf("Line %d: expected %d numbers, got %d",
				i, len(expected[i]), len(result[i]))
		}

		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Fatalf("Mismatch at [%d][%d]: expected %d, got %d",
					i, j, expected[i][j], result[i][j])
			}
		}
	}
}

func TestReportSafety(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected bool
	}{
		// Original given examples
		{
			name:     "Decreasing report 7 6 4 2 1",
			input:    []int{7, 6, 4, 2, 1},
			expected: true,
		},
		{
			name:     "Unsafe - large increase 1 2 7 8 9",
			input:    []int{1, 2, 7, 8, 9},
			expected: false,
		},
		{
			name:     "Unsafe - large decrease 9 7 6 2 1",
			input:    []int{9, 7, 6, 2, 1},
			expected: false,
		},
		{
			name:     "Unsafe - inconsistent direction 1 3 2 4 5",
			input:    []int{1, 3, 2, 4, 5},
			expected: false,
		},
		{
			name:     "Unsafe - flat section 8 6 4 4 1",
			input:    []int{8, 6, 4, 4, 1},
			expected: false,
		},
		{
			name:     "Increasing report 1 3 6 7 9",
			input:    []int{1, 3, 6, 7, 9},
			expected: true,
		},
		// Additional edge cases
		{
			name:     "Empty report",
			input:    []int{},
			expected: false,
		},
		{
			name:     "Single element report",
			input:    []int{5},
			expected: false,
		},
		{
			name:     "Two equal elements",
			input:    []int{5, 5},
			expected: false,
		},
		{
			name:     "Increasing by exactly 1",
			input:    []int{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "Decreasing by exactly 1",
			input:    []int{5, 4, 3, 2, 1},
			expected: true,
		},
		{
			name:     "Increasing by maximum 3",
			input:    []int{1, 4, 7, 10, 13},
			expected: true,
		},
		{
			name:     "Decreasing by maximum 3",
			input:    []int{13, 10, 7, 4, 1},
			expected: true,
		},
		{
			name:     "Unsafe - increase beyond 3",
			input:    []int{1, 5, 9, 13},
			expected: false,
		},
		{
			name:     "Unsafe - decrease beyond 3",
			input:    []int{13, 9, 5, 1},
			expected: false,
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsReportSafe(tc.input)
			if result != tc.expected {
				t.Errorf("For %v: expected %v, got %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestCountSafeReports(t *testing.T) {
	testCases := []struct {
		name     string
		input    [][]int
		expected int
	}{
		{
			name: "Mixed reports",
			input: [][]int{
				{7, 6, 4, 2, 1},  // Safe (decreasing)
				{1, 2, 7, 8, 9},  // Unsafe
				{1, 3, 6, 7, 9},  // Safe (increasing)
				{9, 7, 6, 2, 1},  // Unsafe
				{1, 3, 2, 4, 5},  // Unsafe
				{8, 6, 4, 4, 1},  // Unsafe
				{1, 2, 4, 7, 10}, // Safe (increasing)
			},
			expected: 3,
		},
		{
			name:     "Empty input",
			input:    [][]int{},
			expected: 0,
		},
		{
			name: "All safe reports",
			input: [][]int{
				{1, 2, 3, 4, 5},
				{5, 4, 3, 2, 1},
				{1, 4, 7, 10, 13},
			},
			expected: 3,
		},
		{
			name: "All unsafe reports",
			input: [][]int{
				{1, 3, 2, 4, 5},
				{8, 6, 4, 4, 1},
				{1, 2, 7, 8, 9},
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CountSafeReports(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %d safe reports, got %d", tc.expected, result)
			}
		})
	}
}
