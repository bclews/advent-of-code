package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	empty = '.'
	wall  = '#'
	box   = 'O'
	robot = '@'
)

type Point struct {
	row, col int
}

type Direction struct {
	dRow, dCol int
}

var directions = map[byte]Direction{
	'<': {0, -1},
	'>': {0, 1},
	'^': {-1, 0},
	'v': {1, 0},
}

func parseFile(r io.Reader) (map[Point]byte, []byte, Point, error) {
	var robotPos Point
	grid := make(map[Point]byte)

	scanner := bufio.NewScanner(r)

	// Parse grid
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		if line == "" {
			break
		}

		for col, char := range line {
			pos := Point{row, col}
			if char == robot {
				robotPos = pos
				grid[pos] = empty
			} else {
				grid[pos] = byte(char)
			}
		}
	}

	var moves []byte
	for scanner.Scan() {
		moves = append(moves, []byte(scanner.Text())...)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, Point{}, err
	}

	return grid, moves, robotPos, nil
}

// findNextEmpty finds the next empty cell in the given direction
func findNextEmpty(grid map[Point]byte, pos Point, dir Direction) (Point, error) {
	for {
		pos = Point{pos.row + dir.dRow, pos.col + dir.dCol}
		switch grid[pos] {
		case empty:
			return pos, nil
		case wall:
			return Point{}, errors.New("wall encountered")
		}
	}
}

// calculateScore computes the score based on box positions
func calculateScore(grid map[Point]byte) int {
	var score int
	for pos, char := range grid {
		if char == box {
			score += 100*pos.row + pos.col
		}
	}
	return score
}

// Attempts to move the robot in the given direction
func Solve(grid map[Point]byte, move byte, robotPos Point) Point {
	dir, exists := directions[move]
	if !exists {
		return robotPos
	}

	nextPos := Point{robotPos.row + dir.dRow, robotPos.col + dir.dCol}

	// Check if next position contains a box
	if grid[nextPos] == box {
		if emptyPos, err := findNextEmpty(grid, robotPos, dir); err == nil {
			// Move the box
			grid[nextPos] = empty
			grid[emptyPos] = box
			return nextPos
		}
		return robotPos
	}

	// Move robot if next position is empty
	if grid[nextPos] == empty {
		return nextPos
	}

	return robotPos
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the file
	grid, moves, robotPos, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	for _, move := range moves {
		robotPos = Solve(grid, move, robotPos)
	}

	fmt.Println(calculateScore(grid))
}
