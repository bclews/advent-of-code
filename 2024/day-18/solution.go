package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var GridSize = 70

// Position represents a coordinate in the memory grid
type Position struct {
	x, y int
}

// Grid constants
var (
	startPos = Position{0, 0}
	exitPos  = Position{GridSize, GridSize}
	// Possible movements: right, down, left, up
	directions = []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
)

// represents the state of corrupted memory positions
type MemoryGrid struct {
	corrupted map[Position]struct{}
}

// creates a new memory grid
func NewMemoryGrid() *MemoryGrid {
	return &MemoryGrid{
		corrupted: make(map[Position]struct{}),
	}
}

// marks a position as corrupted
func (g *MemoryGrid) AddCorruption(pos Position) {
	g.corrupted[pos] = struct{}{}
}

// checks if a position is corrupted... no really, it does
func (g *MemoryGrid) IsCorrupted(pos Position) bool {
	_, exists := g.corrupted[pos]
	return exists
}

// checks if a position is within grid bounds
func IsValidPosition(pos Position) bool {
	return pos.x >= 0 && pos.x <= GridSize &&
		pos.y >= 0 && pos.y <= GridSize
}

// finds the shortest path from start to exit using BFS
func (g *MemoryGrid) FindShortestPath() []Position {
	queue := []Position{startPos}
	visited := map[Position]struct{}{startPos: {}}
	parent := make(map[Position]Position)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == exitPos {
			return reconstructPath(parent)
		}

		for _, dir := range directions {
			next := Position{
				x: current.x + dir.x,
				y: current.y + dir.y,
			}

			if !IsValidPosition(next) ||
				g.IsCorrupted(next) ||
				isVisited(visited, next) {
				continue
			}

			visited[next] = struct{}{}
			parent[next] = current
			queue = append(queue, next)
		}
	}

	return nil
}

// builds the path from exit to start using parent mapping
func reconstructPath(parent map[Position]Position) []Position {
	var path []Position
	for pos := exitPos; pos != startPos; pos = parent[pos] {
		path = append(path, pos)
	}
	return path
}

func isVisited(visited map[Position]struct{}, pos Position) bool {
	_, exists := visited[pos]
	return exists
}

// reads corrupted memory positions from stdin
func parseCorruptedPositions(r io.Reader) ([]Position, error) {
	var positions []Position
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		positions = append(positions, Position{x, y})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return positions, nil
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	corruptedPositions, err := parseCorruptedPositions(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing corrupted positions: %v\n", err)
	}

	grid := NewMemoryGrid()

	// Part 1: Find shortest path after first 1024 corrupted positions
	for i := 0; i < 1024 && i < len(corruptedPositions); i++ {
		grid.AddCorruption(corruptedPositions[i])
	}

	if path := grid.FindShortestPath(); path != nil {
		fmt.Println("Part one:", len(path))
	}

	// Part 2: Find position that blocks path to exit
	for _, pos := range corruptedPositions[1024:] {
		grid.AddCorruption(pos)
		if path := grid.FindShortestPath(); path == nil {
			fmt.Printf("Part two:\n%d,%d\n", pos.x, pos.y)
			return
		}
	}
}

