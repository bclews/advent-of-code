package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
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

// isValidLetterForPattern checks if a letter can match the pattern
func isValidLetterForPattern(r rune) bool {
	return unicode.IsUpper(r)
}

// matchPattern checks if a 3x3 grid matches a pattern template
func matchPattern(grid [][]rune, x, y int, pattern Pattern) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gridCell := grid[x+i][y+j]
			patternCell := pattern[i][j]

			// Skip dot (wildcard) matches
			if patternCell == '.' {
				if !isValidLetterForPattern(gridCell) {
					return false
				}
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

// printSubgrid prints the 3x3 subgrid at the given coordinates
func printSubgrid(grid [][]rune, x, y int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			fmt.Printf("%c ", grid[x+i][y+j])
		}
		fmt.Println()
	}
	fmt.Println()
}

// findPatterns searches the grid for matching 3x3 patterns
func findPatterns(grid [][]rune) []string {
	var results []string

	for x := 0; x <= len(grid)-3; x++ {
		for y := 0; y <= len(grid[0])-3; y++ {
			// Center cell must be 'A'
			if grid[x+1][y+1] != 'A' {
				continue
			}

			for _, pattern := range validPatterns {
				if matchPattern(grid, x, y, pattern) {
					matchStr := fmt.Sprintf("Found at (%d, %d)", x, y)
					results = append(results, matchStr)

					// Print the matched subgrid for debugging
					fmt.Println(matchStr)
					printSubgrid(grid, x, y)
				}
			}
		}
	}

	return results
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

	results := findPatterns(grid)
	fmt.Printf("Pattern Matches: %d\n", len(results))
}
