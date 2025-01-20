package main

import (
	"strings"
	"testing"
)

func TestAreCompatible(t *testing.T) {
	tests := []struct {
		key      PinHeights
		lock     PinHeights
		expected bool
	}{
		{
			key:      PinHeights{5, 0, 2, 1, 3},
			lock:     PinHeights{0, 5, 3, 4, 3},
			expected: false,
		},
		{
			key:      PinHeights{4, 3, 4, 0, 2},
			lock:     PinHeights{0, 5, 3, 4, 3},
			expected: false,
		},
		{
			key:      PinHeights{3, 0, 2, 0, 1},
			lock:     PinHeights{0, 5, 3, 4, 3},
			expected: true,
		},
		{
			key:      PinHeights{5, 0, 2, 1, 3},
			lock:     PinHeights{1, 2, 0, 5, 3},
			expected: false,
		},
		{
			key:      PinHeights{4, 3, 4, 0, 2},
			lock:     PinHeights{1, 2, 0, 5, 3},
			expected: true,
		},
		{
			key:      PinHeights{3, 0, 2, 0, 1},
			lock:     PinHeights{1, 2, 0, 5, 3},
			expected: true,
		},
	}

	for _, test := range tests {
		result := areCompatible(test.key, test.lock)
		if result != test.expected {
			t.Errorf("areCompatible(%v, %v) = %v; want %v", test.key, test.lock, result, test.expected)
		}
	}
}

func TestParseAndCountCompatiblePairs(t *testing.T) {
	input := `#####
.####
.####
.####
.#.#.
.#...
.....

#####
##.##
.#.##
...##
...#.
...#.
.....

.....
#....
#....
#...#
#.#.#
#.###
#####

.....
.....
#.#..
###..
###.#
###.#
#####

.....
.....
.....
#....
#.#..
#.#.#
#####`

	r := strings.NewReader(input)
	locks, keys, err := parseFile(r)
	if err != nil {
		t.Fatalf("parseFile returned an error: %v", err)
	}

	if len(locks) != 2 {
		t.Errorf("Expected 2 locks, got %d", len(locks))
	}
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	expectedCount := 3
	result := countCompatiblePairs(locks, keys)
	if result != expectedCount {
		t.Errorf("countCompatiblePairs = %d; want %d", result, expectedCount)
	}
}
