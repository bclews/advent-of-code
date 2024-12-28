package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

// Position represents a coordinate in the maze
type Position struct {
	row, col int
}

// Direction represents possible movement vectors
type Direction struct {
	deltaRow, deltaCol int
}

// MovementDirections defines all possible single-step movements
var MovementDirections = []Direction{
	{0, 1},  // right
	{1, 0},  // down
	{-1, 0}, // up
	{0, -1}, // left
}

// MazeConfig holds the maze configuration and parameters
type MazeConfig struct {
	walls         map[Position]struct{}
	start, finish Position
}

// CheatParams holds the parameters for calculating valid cheats
type CheatParams struct {
	maxDistance  int // maximum allowed cheat distance
	minTimeSaved int // minimum time that must be saved for a valid cheat
}

// ParseMaze reads the maze layout from input and returns the configuration
func ParseMaze(r io.Reader) (MazeConfig, error) {
	walls := make(map[Position]struct{})
	var start, finish Position

	scanner := bufio.NewScanner(r)
	for row := 0; scanner.Scan(); row++ {
		for col, ch := range scanner.Text() {
			pos := Position{row, col}
			switch ch {
			case '#':
				walls[pos] = struct{}{}
			case 'S':
				start = pos
			case 'E':
				finish = pos
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return MazeConfig{}, fmt.Errorf("reading maze: %w", err)
	}

	return MazeConfig{walls, start, finish}, nil
}

// FindShortestPath uses BFS to find the shortest valid path through the maze
func (mc *MazeConfig) FindShortestPath() []Position {
	queue := []Position{mc.start}
	visited := map[Position]struct{}{mc.start: {}}
	parent := make(map[Position]Position)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == mc.finish {
			return reconstructPath(parent, mc.start, mc.finish)
		}

		for _, dir := range MovementDirections {
			next := Position{
				row: current.row + dir.deltaRow,
				col: current.col + dir.deltaCol,
			}

			if _, isWall := mc.walls[next]; isWall {
				continue
			}

			if _, seen := visited[next]; !seen {
				visited[next] = struct{}{}
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}
	return nil
}

// reconstructPath builds the path from start to finish using the parent map
func reconstructPath(parent map[Position]Position, start, finish Position) []Position {
	path := []Position{finish}
	for curr := finish; curr != start; curr = parent[curr] {
		path = append(path, parent[curr])
	}
	slices.Reverse(path)
	return path
}

// CountValidCheats calculates the number of valid cheating opportunities
func CountValidCheats(path []Position, params CheatParams) int {
	validCheats := 0

	for i := 0; i < len(path)-1; i++ {
		for j := i + 1; j < len(path); j++ {
			cheatDistance := ManhattanDistance(path[i], path[j])
			if cheatDistance > 0 && cheatDistance <= params.maxDistance {
				timeSaved := j - i - cheatDistance
				if timeSaved >= params.minTimeSaved {
					validCheats++
				}
			}
		}
	}
	return validCheats
}

// ManhattanDistance calculates the Manhattan distance between two positions
func ManhattanDistance(p1, p2 Position) int {
	return abs(p2.row-p1.row) + abs(p2.col-p1.col)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	maze, err := ParseMaze(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing maze: %v\n", err)
		os.Exit(1)
	}

	path := maze.FindShortestPath()
	if path == nil {
		fmt.Fprintf(os.Stderr, "No valid path found through maze\n")
		os.Exit(1)
	}

	// Part One: 2-step cheats, need to save at least 100 steps
	partOne := CountValidCheats(path, CheatParams{
		maxDistance:  2,
		minTimeSaved: 100,
	})
	fmt.Printf("Part One: %d\n", partOne)

	// Part Two: 20-step cheats, need to save at least 100 steps
	partTwo := CountValidCheats(path, CheatParams{
		maxDistance:  20,
		minTimeSaved: 100,
	})
	fmt.Printf("Part Two: %d\n", partTwo)
}
