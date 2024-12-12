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
		expected []int
	}{
		{
			name:     "Single Blink Example",
			input:    []int{0, 1, 10, 99, 999},
			blinks:   1,
			expected: []int{1, 2024, 1, 0, 9, 9, 2021976},
		},
		{
			name:     "Longer Example After 1 Blink",
			input:    []int{125, 17},
			blinks:   1,
			expected: []int{253000, 1, 7},
		},
		{
			name:     "Longer Example After 2 Blinks",
			input:    []int{125, 17},
			blinks:   2,
			expected: []int{253, 0, 2024, 14168},
		},
		{
			name:     "Longer Example After 3 Blinks",
			input:    []int{125, 17},
			blinks:   3,
			expected: []int{512072, 1, 20, 24, 28676032},
		},
		{
			name:     "Longer Example After 4 Blinks",
			input:    []int{125, 17},
			blinks:   4,
			expected: []int{512, 72, 2024, 2, 0, 2, 4, 2867, 6032},
		},
		{
			name:     "Longer Example After 5 Blinks",
			input:    []int{125, 17},
			blinks:   5,
			expected: []int{1036288, 7, 2, 20, 24, 4048, 1, 4048, 8096, 28, 67, 60, 32},
		},
		{
			name:     "Longer Example After 6 Blinks",
			input:    []int{125, 17},
			blinks:   6,
			expected: []int{2097446912, 14168, 4048, 2, 0, 2, 4, 40, 48, 2024, 40, 48, 80, 96, 2, 8, 6, 7, 6, 0, 3, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := simulateBlinks(tt.input, tt.blinks)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
