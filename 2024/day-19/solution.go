package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// isDesignPossible checks if a given design can be formed using the towel patterns.
func isDesignPossible(patterns []string, design string) bool {
	memo := make(map[string]bool)

	var canForm func(target string) bool
	canForm = func(target string) bool {
		if target == "" {
			return true
		}
		if res, ok := memo[target]; ok {
			return res
		}

		for _, pattern := range patterns {
			if strings.HasPrefix(target, pattern) {
				if canForm(target[len(pattern):]) {
					memo[target] = true
					return true
				}
			}
		}

		memo[target] = false
		return false
	}

	return canForm(design)
}

// countPossibleDesigns counts how many designs can be formed from the patterns.
func countPossibleDesigns(patterns []string, designs []string) int {
	count := 0
	for _, design := range designs {
		if isDesignPossible(patterns, design) {
			count++
		}
	}
	return count
}

// countWaysToFormDesign counts all unique ways a design can be formed using the patterns.
func countWaysToFormDesign(patterns []string, design string) int {
	memo := make(map[string]int)

	var countWays func(target string) int
	countWays = func(target string) int {
		if target == "" {
			return 1
		}
		if res, ok := memo[target]; ok {
			return res
		}

		ways := 0
		for _, pattern := range patterns {
			if strings.HasPrefix(target, pattern) {
				ways += countWays(target[len(pattern):])
			}
		}

		memo[target] = ways
		return ways
	}

	return countWays(design)
}

// countTotalWays counts the total number of ways all designs can be formed.
func countTotalWays(patterns []string, designs []string) int {
	totalWays := 0
	for _, design := range designs {
		totalWays += countWaysToFormDesign(patterns, design)
	}
	return totalWays
}

func parseInput(r io.Reader) ([]string, []string, error) {
	scanner := bufio.NewScanner(r)

	// Read patterns
	scanner.Scan()
	patterns := strings.Split(scanner.Text(), ", ")

	// Skip blank line
	scanner.Scan()

	// Read designs
	var designs []string
	for scanner.Scan() {
		designs = append(designs, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, nil, err
	}

	return patterns, designs, nil
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	patterns, designs, err := parseInput(file)
	if err != nil {
		fmt.Println("Error parsing input")
		return
	}

	// Count possible designs
	possibleCount := countPossibleDesigns(patterns, designs)
	fmt.Printf("Number of possible designs: %d\n", possibleCount)

	// Count total ways to form designs
	totalWays := countTotalWays(patterns, designs)
	fmt.Printf("Total ways to form designs: %d\n", totalWays)
}
