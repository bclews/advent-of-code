package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func FindMultiplicationMatchesPartOne(memory string) [][]string {
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

func FindMultiplicationMatchesPartTwo(memory string) [][]string {
	// This pattern is designed to match three different types of strings:
	// 1. `mul(x,y)` where `x` and `y` are one or more digits (`\d+`).
	// 2. `do()`, which is a literal string.
	// 3. `don't()`, which is also a literal string.
	pattern := `(mul\((\d+),(\d+)\)|do\(\)|don't\(\))`

	// Compile the regex
	regex := regexp.MustCompile(pattern)

	// Return all matches
	return regex.FindAllStringSubmatch(memory, -1)
}

// SumMultiplicationMatchesPartOne calculates the sum of products from multiplication matches
func SumMultiplicationMatchesPartOne(matches [][]string) int {
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

// SumMultiplicationMatchesPartOne calculates the sum of products from multiplication matches
func SumMultiplicationMatchesPartTwo(matches [][]string) int {
	enabled := true
	total := 0
	for _, match := range matches {
		switch match[0] {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				total += sti(match[2]) * sti(match[3])
			}
		}
	}

	return total
}

func sti(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}
	return i
}

func main() {
	filename := "input.txt"

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert the content to a string
	memory := string(content)

	// -- Part One --
	matchesPartOne := FindMultiplicationMatchesPartOne(memory)
	sumPartOne := SumMultiplicationMatchesPartOne(matchesPartOne)
	fmt.Println("Sum of all valid multiplication results for part one:", sumPartOne)

	// -- Part Two --
	matchesPartTwo := FindMultiplicationMatchesPartTwo(memory)
	sumPartTwo := SumMultiplicationMatchesPartTwo(matchesPartTwo)
	fmt.Println("Sum of all valid multiplication results for part two:", sumPartTwo)
}
