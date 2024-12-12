package main

import (
	"reflect"
	"testing"
)

func TestBlinkStones(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		blinks   int
		expected int
	}{
		{
			name:     "Single Blink Example",
			input:    []int{0, 1, 10, 99, 999},
			blinks:   1,
			expected: 7,
		},
		{
			name:     "Longer Example After 1 Blink",
			input:    []int{125, 17},
			blinks:   1,
			expected: 3,
		},
		{
			name:     "Longer Example After 2 Blinks",
			input:    []int{125, 17},
			blinks:   2,
			expected: 4,
		},
		{
			name:     "Longer Example After 3 Blinks",
			input:    []int{125, 17},
			blinks:   3,
			expected: 5,
		},
		{
			name:     "Longer Example After 4 Blinks",
			input:    []int{125, 17},
			blinks:   4,
			expected: 9,
		},
		{
			name:     "Longer Example After 5 Blinks",
			input:    []int{125, 17},
			blinks:   5,
			expected: 13,
		},
		{
			name:     "Longer Example After 6 Blinks",
			input:    []int{125, 17},
			blinks:   6,
			expected: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stoneCount := simulateBlinks(tt.input, tt.blinks)
			if !reflect.DeepEqual(stoneCount, tt.expected) {
				t.Errorf("got %v, want %v", stoneCount, tt.expected)
			}
		})
	}
}
