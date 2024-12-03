package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	// Regex to match valid mul(X,Y) instructions
	// Ensures:
	// - Starts with 'mul'
	// - Parentheses with two 1-3 digit numbers separated by a comma
	regexPartOne = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

	// This pattern is designed to match three different types of strings:
	// 1. `mul(x,y)` where `x` and `y` are one or more digits (`\d+`).
	// 2. `do()`, which is a literal string.
	// 3. `don't()`, which is also a literal string.
	regexPartTwo = regexp.MustCompile(`(mul\((\d+),(\d+)\)|do\(\)|don't\(\))`)
)

func SumMultiplicationMatchesPartOne(matches [][]string) int {
	totalSum := 0
	for _, match := range matches {
		totalSum += sti(match[1]) * sti(match[2])
	}

	return totalSum
}

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
		fmt.Printf("Conversion error: %v\n", err)
		return 0
	}
	return i
}

func main() {
	filename := "input.txt"

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	memory := string(content)

	// -- Part One --
	matchesPartOne := regexPartOne.FindAllStringSubmatch(memory, -1)
	sumPartOne := SumMultiplicationMatchesPartOne(matchesPartOne)
	fmt.Println("Sum of all valid multiplication results for part one:", sumPartOne)

	// -- Part Two --
	matchesPartTwo := regexPartTwo.FindAllStringSubmatch(memory, -1)
	sumPartTwo := SumMultiplicationMatchesPartTwo(matchesPartTwo)
	fmt.Println("Sum of all valid multiplication results for part two:", sumPartTwo)
}
