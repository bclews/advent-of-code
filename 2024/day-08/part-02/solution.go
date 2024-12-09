package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// Read the input map from file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the input map
	input, err := parseFile(file)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Calculate map dimensions
	mapHeight := len(input)
	mapWidth := 0
	if mapHeight > 0 {
		mapWidth = len(input[0])
	}

	// Parse the input map
	antennas := parseMap(input)

	// Find all unique antinodes
	antinodes := findAntinodes(antennas, mapWidth, mapHeight)

	// Print the result
	fmt.Println("Unique antinode locations:", len(antinodes))
}

// Reads the input map from a text file
func parseFile(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Parse the input map into a frequency-to-positions mapping
func parseMap(input []string) map[rune][][2]int {
	antennas := make(map[rune][][2]int)

	for x, row := range input {
		for y, char := range row {
			if char != '.' {
				antennas[char] = append(antennas[char], [2]int{x, y})
			}
		}
	}

	return antennas
}

// Find all unique antinodes in the map
func findAntinodes(antennas map[rune][][2]int, mapWidth, mapHeight int) map[[2]int]bool {
	antinodes := make(map[[2]int]bool)

	for _, positions := range antennas {
		// Check all pairs of positions for the same frequency
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				findAntinodesForPair(positions[i], positions[j], antinodes, mapWidth, mapHeight)
			}
		}
	}

	return antinodes
}

// Find antinodes for a pair of antenna positions
func findAntinodesForPair(p1, p2 [2]int, antinodes map[[2]int]bool, mapWidth, mapHeight int) {
	x1, y1 := p1[0], p1[1]
	x2, y2 := p2[0], p2[1]

	// Include the original antenna positions in antinodes if they form a line
	for x := 0; x < mapHeight; x++ {
		for y := 0; y < mapWidth; y++ {
			// Check if the point forms a line with the two antenna positions
			if isInLine(x1, y1, x2, y2, x, y) {
				antinodes[[2]int{x, y}] = true
			}
		}
	}
}

// Check if a point (x,y) is in the same line as the line through (x1,y1) and (x2,y2)
func isInLine(x1, y1, x2, y2, x, y int) bool {
	// Cross product method to check colinearity
	return (y2-y1)*(x-x1) == (y-y1)*(x2-x1)
}
