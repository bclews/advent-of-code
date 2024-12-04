package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Pattern represents a 3x3 grid search pattern
type Pattern [3][3]rune

var validPatterns = []Pattern{
	{
		{'M', '.', 'S'},
		{'.', 'A', '.'},
		{'M', '.', 'S'},
	},
	{
		{'M', '.', 'M'},
		{'.', 'A', '.'},
		{'S', '.', 'S'},
	},
	{
		{'S', '.', 'S'},
		{'.', 'A', '.'},
		{'M', '.', 'M'},
	},
	{
		{'S', '.', 'M'},
		{'.', 'A', '.'},
		{'S', '.', 'M'},
	},
}

// matchPattern checks if a 3x3 grid matches a pattern template
func matchPattern(grid [][]rune, x, y int, pattern Pattern) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gridCell := grid[x+i][y+j]
			patternCell := pattern[i][j]

			// Skip dot (wildcard) matches
			if patternCell == '.' {
				continue
			}

			// Check exact matches
			if gridCell != patternCell {
				return false
			}
		}
	}
	return true
}

func countXmasGrids(grid [][]rune) int {
	count := 0
	for x := 0; x <= len(grid)-3; x++ {
		for y := 0; y <= len(grid[0])-3; y++ {
			// Center cell must be 'A'
			if grid[x+1][y+1] != 'A' {
				continue
			}

			for _, pattern := range validPatterns {
				if matchPattern(grid, x, y, pattern) {
					count++
				}
			}
		}
	}
	return count
}

func parseFile(r io.Reader) ([][]rune, error) {
	var grid [][]rune

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineRunes := []rune(scanner.Text())
		grid = append(grid, lineRunes)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return grid, nil
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
	grid, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	xmasGrids := countXmasGrids(grid)
	fmt.Printf("Pattern Matches: %d\n", xmasGrids)
}
