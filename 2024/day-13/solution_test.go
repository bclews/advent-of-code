package main

import (
	"testing"
)

func TestSolve(t *testing.T) {
	tests := []struct {
		clawMachine ClawMachine
		expected    int
	}{
		{
			clawMachine: ClawMachine{
				ButtonA: Pair{x: 94, y: 34},
				ButtonB: Pair{x: 22, y: 67},
				Prize:   Pair{x: 8400, y: 5400},
			},
			expected: 280,
		},
		{
			clawMachine: ClawMachine{
				ButtonA: Pair{x: 26, y: 66},
				ButtonB: Pair{x: 67, y: 21},
				Prize:   Pair{x: 12748, y: 12176},
			},
			expected: 0,
		},
		{
			clawMachine: ClawMachine{
				ButtonA: Pair{x: 17, y: 86},
				ButtonB: Pair{x: 84, y: 37},
				Prize:   Pair{x: 7870, y: 6450},
			},
			expected: 200,
		},
		{
			clawMachine: ClawMachine{
				ButtonA: Pair{x: 69, y: 23},
				ButtonB: Pair{x: 27, y: 71},
				Prize:   Pair{x: 18641, y: 10279},
			},
			expected: 0,
		},
	}

	for _, test := range tests {
		result := solve(test.clawMachine)
		if result != test.expected {
			t.Errorf("solve(%v) = %d; expected %d", test.clawMachine, result, test.expected)
		}
	}
}
