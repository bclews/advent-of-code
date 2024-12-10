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
	visited9Positions := make(map[Position]bool)

	// Depth-first search to find all reachable 9-height positions
	var dfs func(r, c, prevHeight int)
	dfs = func(r, c, prevHeight int) {
		// Check bounds and validity of movement
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return
		}

		// Check if current cell is impassable (represented by '.')
		if grid[r][c] == '.' {
			return
		}

		currentHeight := grid[r][c]

		// Check if height change is valid (exactly 1)
		if prevHeight != -1 && currentHeight != prevHeight+1 {
			return
		}

		// Mark 9-height positions
		if currentHeight == 9 {
			visited9Positions[Position{r, c}] = true
		}

		// Try all 4 directions (up, down, left, right)
		dfs(r+1, c, currentHeight)
		dfs(r-1, c, currentHeight)
		dfs(r, c+1, currentHeight)
		dfs(r, c-1, currentHeight)
	}

	// Start DFS from the trailhead with initial height -1
	dfs(start.row, start.col, -1)

	return len(visited9Positions)
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
