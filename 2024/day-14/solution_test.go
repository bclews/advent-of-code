package main

import (
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		input    string
		expected Robot
		wantErr  bool
	}{
		{
			input: "p=0,4 v=3,-3",
			expected: Robot{
				position: Point{x: 0, y: 4},
				velocity: Point{x: 3, y: -3},
			},
			wantErr: false,
		},
		{
			input:    "invalid",
			expected: Robot{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		got, err := parseInput(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseInput(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (got.position != tt.expected.position || got.velocity != tt.expected.velocity) {
			t.Errorf("parseInput(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestCalculatePosition(t *testing.T) {
	// Test case from the problem example
	robot := Robot{
		position: Point{x: 2, y: 4},
		velocity: Point{x: 2, y: -3},
	}

	expectedPositions := []Point{
		{x: 2, y: 4},  // Initial
		{x: 4, y: 1},  // After 1 second
		{x: 6, y: 5},  // After 2 seconds
		{x: 8, y: 2},  // After 3 seconds
		{x: 10, y: 6}, // After 4 seconds
		{x: 1, y: 3},  // After 5 seconds
	}

	width, height := 11, 7
	for seconds, expected := range expectedPositions {
		got := robot.calculatePosition(seconds, width, height)
		if got != expected {
			t.Errorf("After %d seconds: got %v, want %v", seconds, got, expected)
		}
	}
}

func TestCalculateSafetyFactor(t *testing.T) {
	// Example from the problem statement
	input := []string{
		"p=0,4 v=3,-3",
		"p=6,3 v=-1,-3",
		"p=10,3 v=-1,2",
		"p=2,0 v=2,-1",
		"p=0,0 v=1,3",
		"p=3,0 v=-2,-2",
		"p=7,6 v=-1,-3",
		"p=3,0 v=-1,-2",
		"p=9,3 v=2,3",
		"p=7,3 v=-1,2",
		"p=2,4 v=2,-3",
		"p=9,5 v=-3,-3",
	}

	var robots []Robot
	for _, line := range input {
		robot, err := parseInput(line)
		if err != nil {
			t.Fatalf("Failed to parse input: %v", err)
		}
		robots = append(robots, robot)
	}

	got := calculateSafetyFactor(robots, 11, 7, 100)
	want := 12 // From the problem example

	if got != want {
		t.Errorf("calculateSafetyFactor() = %v, want %v", got, want)
	}
}
