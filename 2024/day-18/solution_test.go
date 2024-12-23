package main

import (
	"testing"
)

func TestShortestPathWithCorruptedPositions(t *testing.T) {
	GridSize = 6
	exitPos = Position{GridSize, GridSize}

	grid := NewMemoryGrid()

	// Adding the first 12 corrupted positions from the example
	corruptedPositions := []Position{
		{5, 4},
		{4, 2},
		{4, 5},
		{3, 0},
		{2, 1},
		{6, 3},
		{2, 4},
		{1, 5},
		{0, 6},
		{3, 3},
		{2, 6},
		{5, 1},
	}

	for _, pos := range corruptedPositions {
		grid.AddCorruption(pos)
	}

	// Calculate shortest path
	path := grid.FindShortestPath()

	if path == nil {
		t.Fatalf("Expected a path to the exit, but none was found")
	}

	expectedLength := 22
	if len(path) != expectedLength {
		t.Errorf("Expected path length %d, but got %d", expectedLength, len(path))
	}
}

func TestFirstByteBlockingExit(t *testing.T) {
	GridSize = 6
	exitPos = Position{GridSize, GridSize}

	grid := NewMemoryGrid()

	// Adding corrupted positions until the path is blocked
	corruptedPositions := []Position{
		{5, 4},
		{4, 2},
		{4, 5},
		{3, 0},
		{2, 1},
		{6, 3},
		{2, 4},
		{1, 5},
		{0, 6},
		{3, 3},
		{2, 6},
		{5, 1},
		{1, 2},
		{5, 5},
		{2, 5},
		{6, 5},
		{1, 4},
		{0, 4},
		{6, 4},
		{1, 1}, // up to here path exists
		{6, 1}, // this byte blocks the path
	}

	var blockingByte Position
	for _, pos := range corruptedPositions {
		grid.AddCorruption(pos)
		if path := grid.FindShortestPath(); path == nil {
			blockingByte = pos
			break
		}
	}

	expectedBlockingByte := Position{6, 1}
	if blockingByte != expectedBlockingByte {
		t.Errorf("Expected blocking byte at %v, but got %v", expectedBlockingByte, blockingByte)
	}
}

