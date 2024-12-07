package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		numbers   []int
		operators []rune
		expected  int
	}{
		// Original test cases
		{[]int{10, 19}, []rune{'*'}, 190},
		{[]int{81, 40, 27}, []rune{'+', '*'}, 3267},
		{[]int{11, 6, 16, 20}, []rune{'+', '*', '+'}, 292},

		// Concatenation test cases
		{[]int{15, 6}, []rune{'|'}, 156},
		{[]int{6, 8, 6, 15}, []rune{'*', '|', '*'}, 7290},
		{[]int{17, 8, 14}, []rune{'|', '+'}, 192},
	}

	for _, test := range tests {
		result := evaluate(test.numbers, test.operators)
		if result != test.expected {
			t.Errorf("evaluate(%v, %v) = %d; want %d", test.numbers, test.operators, result, test.expected)
		}
	}
}

func TestCalculateCalibrationResult(t *testing.T) {
	// Test cases for part one and part two
	equations := map[int][]int{
		// Part one equations
		190:  {10, 19},
		3267: {81, 40, 27},
		292:  {11, 6, 16, 20},
		83:   {17, 5},

		// Part two additional equations
		156:  {15, 6},
		7290: {6, 8, 6, 15},
		192:  {17, 8, 14},
	}

	// Total calibration result now includes part two equations
	expected := 11387
	result := calculateCalibrationResult(equations)

	if result != expected {
		t.Errorf("calculateCalibrationResult() = %d; want %d", result, expected)
	}
}

func TestGenerateOperatorCombinations(t *testing.T) {
	// Test to ensure the number of combinations is correct for different lengths
	testCases := []struct {
		n               int
		expectedCombLen int
	}{
		{1, 3},
		{2, 9},
		{3, 27},
	}

	for _, tc := range testCases {
		combinations := generateOperatorCombinations(tc.n)
		if len(combinations) != tc.expectedCombLen {
			t.Errorf("generateOperatorCombinations(%d) produced %d combinations; want %d",
				tc.n, len(combinations), tc.expectedCombLen)
		}

		// Verify that each combination has the correct length
		for _, combo := range combinations {
			if len(combo) != tc.n {
				t.Errorf("Combination %v has incorrect length %d; want %d", combo, len(combo), tc.n)
			}
		}
	}
}
