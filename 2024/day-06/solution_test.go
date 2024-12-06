package main

import (
	"testing"
)

// TestParseMap tests the initial map parsing functionality
func TestParseMap(t *testing.T) {
	testCases := []struct {
		name           string
		inputMap       []string
		expectedStartX int
		expectedStartY int
		expectedDir    Direction
	}{
		{
			name: "Simple Map with Up-Facing Guard",
			inputMap: []string{
				"....#.....",
				".........#",
				"..........",
				"..#.......",
				".......#..",
				"..........",
				".#..^.....",
				"........#.",
				"#.........",
				"......#...",
			},
			expectedStartX: 4,
			expectedStartY: 6,
			expectedDir:    Up,
		},
		{
			name: "Map with Guard in Different Position",
			inputMap: []string{
				"#..........",
				"....^......",
				"...........",
			},
			expectedStartX: 4,
			expectedStartY: 1,
			expectedDir:    Up,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMap(tc.inputMap)

			if m.GuardX != tc.expectedStartX {
				t.Errorf("Expected guard X position %d, got %d", tc.expectedStartX, m.GuardX)
			}

			if m.GuardY != tc.expectedStartY {
				t.Errorf("Expected guard Y position %d, got %d", tc.expectedStartY, m.GuardY)
			}

			if m.GuardDir != tc.expectedDir {
				t.Errorf("Expected guard direction %v, got %v", tc.expectedDir, m.GuardDir)
			}
		})
	}
}

// TestGuardMovement tests the guard's movement and turning logic
func TestGuardMovement(t *testing.T) {
	testCases := []struct {
		name           string
		inputMap       []string
		expectedVisits int
	}{
		{
			name: "Simple Open Path",
			inputMap: []string{
				".........",
				"...^.....",
				".........",
			},
			expectedVisits: 2,
		},
		{
			name: "Path with Obstacles",
			inputMap: []string{
				"....#.....",
				".........#",
				"..........",
				"..#.......",
				".......#..",
				"..........",
				".#..^.....",
				"........#.",
				"#.........",
				"......#...",
			},
			expectedVisits: 41, // From the problem description example
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMap(tc.inputMap)
			visits := m.SimulateGuardPatrol()

			if visits != tc.expectedVisits {
				t.Errorf("Expected %d visited positions, got %d", tc.expectedVisits, visits)
			}
		})
	}
}

// TestObstacleDetection tests the guard's ability to detect and avoid obstacles
func TestObstacleDetection(t *testing.T) {
	testCases := []struct {
		name             string
		inputMap         []string
		expectedFinalDir Direction
	}{
		{
			name: "Obstacle Ahead Causes Turn",
			inputMap: []string{
				"..#......",
				"..^......",
				"..#......",
			},
			expectedFinalDir: Right,
		},
		{
			name: "Multiple Obstacle Turns",
			inputMap: []string{
				"..#.#....",
				"..^......",
				".........",
			},
			expectedFinalDir: Right,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMap(tc.inputMap)
			m.SimulateGuardPatrol()

			if m.GuardDir != tc.expectedFinalDir {
				t.Errorf("Expected final direction %v, got %v", tc.expectedFinalDir, m.GuardDir)
			}
		})
	}
}
