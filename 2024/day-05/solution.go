package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PageOrderChecker struct {
	rules map[int]map[int]bool
}

func NewPageOrderChecker() *PageOrderChecker {
	return &PageOrderChecker{
		rules: make(map[int]map[int]bool),
	}
}

func (poc *PageOrderChecker) AddRule(before, after int) {
	if poc.rules[before] == nil {
		poc.rules[before] = make(map[int]bool)
	}
	poc.rules[before][after] = true
}

func (poc *PageOrderChecker) IsValidOrder(pages []int) bool {
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			// Check if the later page must come before the earlier page
			if poc.rules[pages[j]] != nil && poc.rules[pages[j]][pages[i]] {
				return false
			}
		}
	}
	return true
}

func (poc *PageOrderChecker) ReorderUpdate(pages []int) []int {
	// Create a copy of the pages slice to avoid modifying the original
	updatedPages := make([]int, len(pages))
	copy(updatedPages, pages)

	// Continue attempting to reorder until the order is valid
	for !poc.IsValidOrder(updatedPages) {
		// Find the first pair of pages that violate the ordering rules
		violatingPair := poc.findViolatingPair(updatedPages)
		if violatingPair == nil {
			break // No clear way to resolve conflicts
		}

		// Swap the pages to attempt to resolve the violation
		updatedPages[violatingPair[0]], updatedPages[violatingPair[1]] = updatedPages[violatingPair[1]], updatedPages[violatingPair[0]]
	}

	return updatedPages
}

func (poc *PageOrderChecker) findViolatingPair(pages []int) []int {
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			// Check if the later page must come before the earlier page
			if poc.rules[pages[j]] != nil && poc.rules[pages[j]][pages[i]] {
				return []int{i, j}
			}
		}
	}
	return nil
}

func parseFile(r io.Reader) ([]string, [][]int, error) {
	scanner := bufio.NewScanner(r)

	rules, err := extractRules(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract rules: %v", err)
	}

	updates, err := extractUpdates(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract updates: %v", err)
	}

	return rules, updates, nil
}

func extractRules(scanner *bufio.Scanner) ([]string, error) {
	var rules []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		rules = append(rules, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return rules, nil
}

func extractUpdates(scanner *bufio.Scanner) ([][]int, error) {
	var updates [][]int
	for scanner.Scan() {
		line := scanner.Text()
		pageStrings := strings.Split(line, ",")
		pageInts, err := convertStringsToInts(pageStrings)
		if err != nil {
			return nil, fmt.Errorf("error parsing updates: %v", err)
		}
		updates = append(updates, pageInts)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return updates, nil
}

func convertStringsToInts(strings []string) ([]int, error) {
	ints := make([]int, len(strings))
	for i, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		ints[i] = num
	}
	return ints, nil
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
	rules, updates, err := parseFile(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	poc := NewPageOrderChecker()

	// Add rules
	for _, rule := range rules {
		parts := strings.Split(rule, "|")
		before, _ := strconv.Atoi(parts[0])
		after, _ := strconv.Atoi(parts[1])
		poc.AddRule(before, after)
	}

	// Check updates and sum middle pages
	var middlePages []int
	for _, update := range updates {
		// If not already in valid order, try to reorder
		if !poc.IsValidOrder(update) {
			reorderedUpdate := poc.ReorderUpdate(update)

			// Only consider if reordering was successful
			if poc.IsValidOrder(reorderedUpdate) {
				middlePage := reorderedUpdate[len(reorderedUpdate)/2]
				middlePages = append(middlePages, middlePage)
			}
		}
	}

	// Calculate sum of middle pages
	sum := 0
	for _, page := range middlePages {
		sum += page
	}
	fmt.Println("Sum of middle pages:", sum)
}
