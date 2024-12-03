package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// FindMultiplicationMatches finds all valid mul(X,Y) instructions in the memory string
func FindMultiplicationMatches(memory string) [][]string {
	// Regex to match valid mul(X,Y) instructions
	// Ensures:
	// - Starts with 'mul'
	// - Parentheses with two 1-3 digit numbers separated by a comma
	pattern := `mul\((\d{1,3}),(\d{1,3})\)`

	// Compile the regex
	regex := regexp.MustCompile(pattern)

	// Return all matches
	return regex.FindAllStringSubmatch(memory, -1)
}

// SumMultiplicationMatches calculates the sum of products from multiplication matches
func SumMultiplicationMatches(matches [][]string) int {
	// Sum of multiplication results
	totalSum := 0

	// Process each valid match
	for _, match := range matches {
		// match[0] is full match, match[1] is first number, match[2] is second number
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])

		// Add product to total sum
		totalSum += x * y
	}

	return totalSum
}

func main() {
	filename := "input.txt"

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert the content to a string
	memory := string(content)

	// Find all matches first
	matches := FindMultiplicationMatches(memory)

	// Then sum the matches
	sum := SumMultiplicationMatches(matches)

	// Print the content
	fmt.Println("Sum of all valid multiplication results:", sum)
}
