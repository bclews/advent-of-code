package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestPageOrderChecker(t *testing.T) {
	// Test cases for IsValidOrder
	testCases := []struct {
		name    string
		rules   []int
		pages   []int
		isValid bool
	}{
		{
			name:    "Valid order: 75, 47, 61, 53, 29",
			rules:   []int{75, 47, 61, 53, 29},
			pages:   []int{75, 47, 61, 53, 29},
			isValid: true,
		},
		{
			name:    "Valid order: 97, 61, 53, 29, 13",
			rules:   []int{97, 61, 53, 29, 13},
			pages:   []int{97, 61, 53, 29, 13},
			isValid: true,
		},
		{
			name:    "Valid order: 75, 29, 13",
			rules:   []int{75, 29, 13},
			pages:   []int{75, 29, 13},
			isValid: true,
		},
		{
			name:    "Invalid order: 75, 97, 47, 61, 53",
			rules:   []int{97, 75, 47, 61, 53},
			pages:   []int{75, 97, 47, 61, 53},
			isValid: false,
		},
		{
			name:    "Invalid order: 61, 13, 29",
			rules:   []int{29, 13},
			pages:   []int{61, 13, 29},
			isValid: false,
		},
		{
			name:    "Invalid order: 97, 13, 75, 29, 47",
			rules:   []int{97, 13, 75, 29, 47},
			pages:   []int{97, 13, 75, 29, 47},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new page order checker for each test
			poc := NewPageOrderChecker()

			// Add rules for this test case
			// We'll create rules based on the input and basic constraints
			ruleSet := []string{
				"75|47", "75|61", "75|53", "75|29",
				"47|61", "47|53", "47|29",
				"61|53", "61|29",
				"53|29",
				"97|61", "97|53", "97|29", "97|13",
				"97|75", "97|47",
				"29|13",
			}

			// Add the rules to the checker
			for _, rule := range ruleSet {
				parts := strings.Split(rule, "|")
				before, _ := strconv.Atoi(parts[0])
				after, _ := strconv.Atoi(parts[1])
				poc.AddRule(before, after)
			}

			// Check the order
			isValid := poc.IsValidOrder(tc.pages)
			if isValid != tc.isValid {
				t.Errorf("Expected valid order to be %v, but got %v", tc.isValid, isValid)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// Test parsing a complete input
	input := `75|47
75|61
75|53
75|29

75,47,61,53,29
97,61,53,29,13`

	reader := strings.NewReader(input)
	rules, updates, err := parseFile(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check rules
	expectedRules := []string{"75|47", "75|61", "75|53", "75|29"}
	if len(rules) != len(expectedRules) {
		t.Errorf("Expected %d rules, got %d", len(expectedRules), len(rules))
	}
	for i, rule := range expectedRules {
		if rules[i] != rule {
			t.Errorf("Expected rule %s, got %s", rule, rules[i])
		}
	}

	// Check updates
	expectedUpdates := [][]int{
		{75, 47, 61, 53, 29},
		{97, 61, 53, 29, 13},
	}
	if len(updates) != len(expectedUpdates) {
		t.Errorf("Expected %d updates, got %d", len(expectedUpdates), len(updates))
	}
	for i, update := range expectedUpdates {
		for j, page := range update {
			if updates[i][j] != page {
				t.Errorf("Expected update %v, got %v", update, updates[i])
				break
			}
		}
	}
}

func TestConvertStringsToInts(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []int
		hasError bool
	}{
		{
			input:    []string{"1", "2", "3"},
			expected: []int{1, 2, 3},
			hasError: false,
		},
		{
			input:    []string{"75", "47", "61"},
			expected: []int{75, 47, 61},
			hasError: false,
		},
		{
			input:    []string{"1", "a", "3"},
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		result, err := convertStringsToInts(tc.input)

		if tc.hasError {
			if err == nil {
				t.Errorf("Expected an error for input %v, but got none", tc.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input %v: %v", tc.input, err)
			continue
		}

		if len(result) != len(tc.expected) {
			t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
			continue
		}

		for i, val := range tc.expected {
			if result[i] != val {
				t.Errorf("Expected %d at index %d, got %d", val, i, result[i])
			}
		}
	}
}
