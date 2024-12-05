package main

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestPageOrderChecker(t *testing.T) {
	// Test cases for IsValidOrder and ReorderUpdate
	testCases := []struct {
		name    string
		rules   []string
		pages   []int
		isValid bool
	}{
		{
			name:    "Valid order: 75, 47, 61, 53, 29",
			rules:   []string{"75|47", "75|61", "75|53", "75|29", "47|61", "47|53", "47|29", "61|53", "61|29", "53|29"},
			pages:   []int{75, 47, 61, 53, 29},
			isValid: true,
		},
		{
			name:    "Reorder: 75, 97, 47, 61, 53",
			rules:   []string{"75|47", "75|61", "75|53", "97|75", "97|47", "97|61", "97|53"},
			pages:   []int{75, 97, 47, 61, 53},
			isValid: false,
		},
		{
			name:    "Reorder: 61, 13, 29",
			rules:   []string{"29|13", "61|29"},
			pages:   []int{61, 13, 29},
			isValid: false,
		},
		{
			name:    "Reorder: 97, 13, 75, 29, 47",
			rules:   []string{"97|13", "97|75", "97|29", "97|47", "75|47", "75|29", "29|13"},
			pages:   []int{97, 13, 75, 29, 47},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new page order checker for each test
			poc := NewPageOrderChecker()

			// Add rules for this test case
			for _, rule := range tc.rules {
				parts := strings.Split(rule, "|")
				before, _ := strconv.Atoi(parts[0])
				after, _ := strconv.Atoi(parts[1])
				poc.AddRule(before, after)
			}

			// Check initial order
			isValid := poc.IsValidOrder(tc.pages)
			if isValid != tc.isValid {
				t.Errorf("Expected initial valid order to be %v, but got %v", tc.isValid, isValid)
			}

			// Attempt to reorder
			reorderedPages := poc.ReorderUpdate(tc.pages)

			// Check if reordered pages are valid
			if !poc.IsValidOrder(reorderedPages) {
				t.Errorf("Reordered pages %v are not valid", reorderedPages)
			}

			// Check that reordered pages contain the same elements as original
			if !slicesHaveSameElements(tc.pages, reorderedPages) {
				t.Errorf("Reordered pages %v do not contain the same elements as original %v",
					reorderedPages, tc.pages)
			}
		})
	}
}

// Helper function to check if two slices have the same elements (regardless of order)
func slicesHaveSameElements(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	// Create copies to avoid modifying original slices
	copyA := make([]int, len(a))
	copyB := make([]int, len(b))
	copy(copyA, a)
	copy(copyB, b)

	// Sort both slices
	sort.Ints(copyA)
	sort.Ints(copyB)

	// Compare sorted slices
	for i := range copyA {
		if copyA[i] != copyB[i] {
			return false
		}
	}

	return true
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
