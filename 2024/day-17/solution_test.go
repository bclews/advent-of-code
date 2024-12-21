package main

import (
	"slices"
	"strings"
	"testing"
)

// Tests based on examples from the problem
func TestComputer(t *testing.T) {
	// Part 1 tests
	tests := []struct {
		name     string
		a        int
		b        int
		c        int
		program  []int
		expected []int
		wantB    int // Expected value of register B after execution
	}{
		{
			name:     "Example 1 - BST with C=9",
			a:        0,
			b:        0,
			c:        9,
			program:  []int{2, 6},
			expected: []int{},
			wantB:    1,
		},
		{
			name:     "Example 2 - Multiple outputs with A=10",
			a:        10,
			b:        0,
			c:        0,
			program:  []int{5, 0, 5, 1, 5, 4},
			expected: []int{0, 1, 2},
			wantB:    0,
		},
		{
			name:     "Example 3 - Complex program with A=2024",
			a:        2024,
			b:        0,
			c:        0,
			program:  []int{0, 1, 5, 4, 3, 0},
			expected: []int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0},
			wantB:    0,
		},
		{
			name:     "Example 4 - BXL with B=29",
			a:        0,
			b:        29,
			c:        0,
			program:  []int{1, 7},
			expected: []int{},
			wantB:    26,
		},
		{
			name:     "Example 5 - BXC with B=2024, C=43690",
			a:        0,
			b:        2024,
			c:        43690,
			program:  []int{4, 0},
			expected: []int{},
			wantB:    44354,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			computer := NewComputer(tt.a, tt.b, tt.c, tt.program)
			output := computer.Run()

			// Convert output to string for easier comparison in error messages
			gotStr := strings.Join(formatOutput(output), ",")
			wantStr := strings.Join(formatOutput(tt.expected), ",")

			if !slices.Equal(output, tt.expected) {
				t.Errorf("output = %v, want %v", gotStr, wantStr)
			}

			if computer.B != tt.wantB {
				t.Errorf("Register B = %d, want %d", computer.B, tt.wantB)
			}
		})
	}

	// Part 2 test
	t.Run("Part 2 - Find self-output A value", func(t *testing.T) {
		program := []int{0, 3, 5, 4, 3, 0}
		expected := 117440

		result := findSelfReplicatingValue(0, 0, program)
		if result != expected {
			t.Errorf("findSelfReplicatingValue() = %d, want %d", result, expected)
		}

		// Verify the output matches the program
		computer := NewComputer(result, 0, 0, program)
		output := computer.Run()
		if !slices.Equal(output, program) {
			t.Errorf("Output %v does not match program %v with A=%d", output, program, result)
		}
	})
}
