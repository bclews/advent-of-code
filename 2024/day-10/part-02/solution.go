package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Position struct {
	row, col int
}

func parseFile(r io.Reader) ([][]int, error) {
	var result [][]int
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		var row []int
		for _, char := range line {
			num, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("failed to parse character %c: %v", char, err)
			}
			row = append(row, num)
		}
		result = append(result, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func CalculateTrailheadScores(grid [][]int) int {
	// Find all trailheads (positions with height 0)
	trailheads := findTrailheads(grid)

	// Calculate scores for each trailhead
	totalScore := 0
	for _, trailhead := range trailheads {
		totalScore += calculateTrailheadScore(grid, trailhead)
	}

	return totalScore
}

func findTrailheads(grid [][]int) []Position {
	trailheads := []Position{}

	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			// Check if the current position is a trailhead (height 0)
			if grid[r][c] == 0 {
				trailheads = append(trailheads, Position{r, c})
			}
		}
	}

	return trailheads
}

func calculateTrailheadScore(grid [][]int, start Position) int {
	rows, cols := len(grid), len(grid[0])
	trailCount := 0

	// Function to check if a trail is valid
	isValidTrail := func(trail []Position) bool {
		if len(trail) == 0 {
			return false
		}

		// Must start at trailhead (height 0)
		if grid[trail[0].row][trail[0].col] != 0 {
			return false
		}

		// Must end at height 9
		lastPos := trail[len(trail)-1]
		if grid[lastPos.row][lastPos.col] != 9 {
			return false
		}

		// Check height increases exactly by 1 at each step
		for i := 1; i < len(trail); i++ {
			prevHeight := grid[trail[i-1].row][trail[i-1].col]
			currHeight := grid[trail[i].row][trail[i].col]

			// Ensure height increases by exactly 1
			if currHeight != prevHeight+1 {
				return false
			}

			// Ensure move is orthogonal (not diagonal)
			rowDiff := trail[i].row - trail[i-1].row
			colDiff := trail[i].col - trail[i-1].col
			if (rowDiff != 0 && colDiff != 0) || (abs(rowDiff)+abs(colDiff) != 1) {
				return false
			}
		}

		return true
	}

	// Recursive function to find all possible trails
	var findTrails func(current Position, trail []Position, visited map[Position]bool)
	findTrails = func(current Position, trail []Position, visited map[Position]bool) {
		// Check bounds
		if current.row < 0 || current.row >= rows ||
			current.col < 0 || current.col >= cols {
			return
		}

		// Check for impassable tile
		if grid[current.row][current.col] == '.' {
			return
		}

		// Prevent revisiting positions in the same trail
		if visited[current] {
			return
		}

		// Add current position to trail
		newTrail := append(trail, current)
		newVisited := copyMap(visited)
		newVisited[current] = true

		// If trail is valid, increment trail count
		if isValidTrail(newTrail) {
			trailCount++
		}

		// Try moving in 4 directions
		directions := []Position{
			{current.row + 1, current.col},
			{current.row - 1, current.col},
			{current.row, current.col + 1},
			{current.row, current.col - 1},
		}

		for _, dir := range directions {
			// Ensure height increases by exactly 1
			if dir.row >= 0 && dir.row < rows &&
				dir.col >= 0 && dir.col < cols &&
				grid[dir.row][dir.col] != '.' {
				currentHeight := grid[current.row][current.col]
				nextHeight := grid[dir.row][dir.col]

				if nextHeight == currentHeight+1 {
					findTrails(dir, newTrail, newVisited)
				}
			}
		}
	}

	// Start from trailhead with no prior trail
	findTrails(start, []Position{}, make(map[Position]bool))

	return trailCount
}

// Helper function to create a deep copy of a map
func copyMap(m map[Position]bool) map[Position]bool {
	newMap := make(map[Position]bool)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
	input, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Calculate the total score
	totalScore := CalculateTrailheadScores(input)
	fmt.Println("Total score:", totalScore)
}
