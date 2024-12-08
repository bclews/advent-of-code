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

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
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
		// Check all pairs of antennas with the same frequency
		for i := 0; i < len(positions); i++ {
			for j := i + 1; j < len(positions); j++ {
				p1, p2 := positions[i], positions[j]

				// Check if they create valid antinodes
				addAntinodes(p1, p2, antinodes, mapWidth, mapHeight)
			}
		}
	}

	return antinodes
}

// Add valid antinodes created by a pair of antennas to the set
func addAntinodes(p1, p2 [2]int, antinodes map[[2]int]bool, mapWidth, mapHeight int) {
	// Calculate midpoints for potential antinodes
	x1, y1 := p1[0], p1[1]
	x2, y2 := p2[0], p2[1]

	// Midpoints to check (in both directions)
	candidates := [][2]int{
		{2*x1 - x2, 2*y1 - y2},
		{2*x2 - x1, 2*y2 - y1},
	}

	// Add valid antinodes within bounds
	for _, c := range candidates {
		if c[0] >= 0 && c[0] < mapHeight && c[1] >= 0 && c[1] < mapWidth {
			antinodes[c] = true
		}
	}
}
