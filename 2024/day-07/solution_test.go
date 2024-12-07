package main

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		numbers   []int
		operators []rune
		expected  int
	}{
		{[]int{10, 19}, []rune{'*'}, 190},
		{[]int{81, 40, 27}, []rune{'+', '*'}, 3267},
		{[]int{11, 6, 16, 20}, []rune{'+', '*', '+'}, 292},
	}

	for _, test := range tests {
		result := evaluate(test.numbers, test.operators)
		if result != test.expected {
			t.Errorf("evaluate(%v, %v) = %d; want %d", test.numbers, test.operators, result, test.expected)
		}
	}
}

func TestCalculateCalibrationResult(t *testing.T) {
	equations := map[int][]int{
		190:  {10, 19},
		3267: {81, 40, 27},
		292:  {11, 6, 16, 20},
		83:   {17, 5},
	}

	expected := 3749
	result := calculateCalibrationResult(equations)

	if result != expected {
		t.Errorf("calculateCalibrationResult() = %d; want %d", result, expected)
	}
}
