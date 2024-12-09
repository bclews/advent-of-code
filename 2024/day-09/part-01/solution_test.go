package main

import "testing"

func TestSolve(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Example case",
			input:    []string{"2333133121414131402"},
			expected: "1928",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := solve(tc.input)
			if result != tc.expected {
				t.Errorf("For input %v, expected %s, got %s",
					tc.input, tc.expected, result)
			}
		})
	}
}
