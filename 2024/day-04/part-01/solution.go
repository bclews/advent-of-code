package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func FindWordOccurrences(grid [][]rune, word string) int {
	if len(grid) == 0 || len(grid[0]) == 0 || len(word) == 0 {
		return 0
	}

	// Define the 8 possible directions (dx, dy)
	directions := [][2]int{
		{0, 1},   // Right
		{0, -1},  // Left
		{1, 0},   // Down
		{-1, 0},  // Up
		{1, 1},   // Down-Right
		{1, -1},  // Down-Left
		{-1, 1},  // Up-Right
		{-1, -1}, // Up-Left
	}

	wordLen := len(word)
	wordRunes := []rune(word)
	rows := len(grid)
	cols := len(grid[0])
	count := 0

	// Helper function to check if the word exists starting at (x, y) in a given direction
	isValid := func(x, y, dx, dy int) bool {
		for i := 0; i < wordLen; i++ {
			nx, ny := x+i*dx, y+i*dy
			if nx < 0 || ny < 0 || nx >= rows || ny >= cols || grid[nx][ny] != wordRunes[i] {
				return false
			}
		}
		return true
	}

	// Iterate through each cell in the grid
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// If the current cell matches the first letter of the word
			if grid[i][j] == wordRunes[0] {
				// Check all 8 directions
				for _, dir := range directions {
					if isValid(i, j, dir[0], dir[1]) {
						count++
					}
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

	// Find the word "XMAS" in the grid
	word := "XMAS"
	result := FindWordOccurrences(grid, word)
	fmt.Printf("The word '%s' occurs %d times in the grid.\n", word, result)
}
