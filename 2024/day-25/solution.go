package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	numPins       = 5 // Number of pins in each lock/key
	schematicRows = 7 // Height of each schematic diagram
	minPinSpace   = 2 // Minimum space needed between key and lock pins
)

// represents the heights of pins in a lock or key
type PinHeights [numPins]int

// Schematic represents a complete lock or key pattern
type Schematic struct {
	heights PinHeights
	isLock  bool
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	locks, keys, error := parseFile(file)
	if error != nil {
		fmt.Println("Error parsing file:", error)
		return
	}

	compatiblePairs := countCompatiblePairs(locks, keys)
	fmt.Println(compatiblePairs)
}

// counts how many lock-key pairs can work together
func countCompatiblePairs(locks, keys []PinHeights) int {
	var count int
	for _, key := range keys {
		for _, lock := range locks {
			if areCompatible(key, lock) {
				count++
			}
		}
	}
	return count
}

// checks if a key and lock can work together without pin collision
func areCompatible(key, lock PinHeights) bool {
	for pinIndex := 0; pinIndex < numPins; pinIndex++ {
		totalHeight := key[pinIndex] + lock[pinIndex]
		if totalHeight > schematicRows-minPinSpace {
			return false
		}
	}
	return true
}

// reads and processes the schematic input from stdin
func parseFile(r io.Reader) ([]PinHeights, []PinHeights, error) {
	var locks, keys []PinHeights
	scanner := bufio.NewScanner(r)

	for {
		schematic := readSchematic(scanner)
		if schematic == nil {
			break
		}

		if schematic.isLock {
			locks = append(locks, schematic.heights)
		} else {
			keys = append(keys, schematic.heights)
		}

		// Skip empty line between schematics
		scanner.Scan()
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return locks, keys, nil
}

// reads a single schematic diagram and converts it to pin heights
func readSchematic(scanner *bufio.Scanner) *Schematic {
	var rows [schematicRows]string

	// Read all rows of the schematic
	for i := 0; i < schematicRows; i++ {
		if !scanner.Scan() {
			return nil
		}
		rows[i] = scanner.Text()
	}

	// Calculate pin heights and determine if it's a lock
	heights := calculatePinHeights(rows)
	isLock := rows[0] == "#####" && rows[schematicRows-1] == "....."

	return &Schematic{
		heights: heights,
		isLock:  isLock,
	}
}

// converts a schematic diagram to pin heights
func calculatePinHeights(rows [schematicRows]string) PinHeights {
	var heights PinHeights

	for col := 0; col < numPins; col++ {
		var pinHeight int
		// Count '#' characters in each column, excluding top and bottom rows
		for row := 1; row < schematicRows-1; row++ {
			if rows[row][col] == '#' {
				pinHeight++
			}
		}
		heights[col] = pinHeight
	}

	return heights
}

