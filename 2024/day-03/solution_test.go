package main

import (
	"testing"
)

func TestSolveCorruptedMemory(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Example from problem statement",
			input:    "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
			expected: 161,
		},
		{
			name:     "No valid instructions",
			input:    "random text with no mul(X,Y) instructions",
			expected: 0,
		},
		{
			name:     "Multiple valid instructions",
			input:    "mul(10,2)abc mul(5,3)def mul(7,6)",
			expected: 20 + 15 + 42,
		},
		{
			name:     "Malformed instructions ignored",
			input:    "mul[1,2] mul(3,4) do_mul(5,6) mul(7,8)",
			expected: 12 + 30 + 56,
		},
		{
			name:     "Single valid instruction",
			input:    "some text mul(12,3) more text",
			expected: 36,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matches := FindMultiplicationMatches(tc.input)
			sum := SumMultiplicationMatches(matches)
			if sum != tc.expected {
				t.Errorf("For input '%s': expected %d, got %d",
					tc.input, tc.expected, sum)
			}
		})
	}
}
