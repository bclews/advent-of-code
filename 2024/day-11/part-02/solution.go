package main

import (
	"fmt"
	"strconv"
	"sync"
)

/*
	Optimisation Technique:
    Realise the order of the numbers doesn’t matter and recursively score the numbers.

  Memoization:
    Memoize the recursive score function to avoid redundant calculations.

  Cache Implementation:
    Use a sparse array to store the results for efficient lookups.
*/

// Store the results of previous calculations to avoid redundant computations.
// It uses a `sync.RWMutex` to ensure thread-safe access to the cache
type stoneCache struct {
	cache map[cacheKey]int
	mu    sync.RWMutex
}

type cacheKey struct {
	number int
	blinks int
}

func newStoneCache() *stoneCache {
	return &stoneCache{
		cache: make(map[cacheKey]int),
	}
}

// Uses a read lock (`RLock`) to allow concurrent reads.
func (sc *stoneCache) get(number, blinks int) (int, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	result, exists := sc.cache[cacheKey{number, blinks}]
	return result, exists
}

// Uses a write lock (`Lock`) to ensure exclusive access during updates.
func (sc *stoneCache) set(number, blinks, result int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache[cacheKey{number, blinks}] = result
}

/*
This recursive function simulates the transformation of a stone based on its
current state (`number`) and the remaining number of blinks.

It first checks if the result is already cached.

If not, it applies one of three rules based on the stone's current state:
  - If the number is zero, it transforms to one.
  - If the number has an even number of digits, it splits the number into two halves and recursively processes each half.
  - Otherwise, it multiplies the number by 2024 and recursively processes the result.

The function caches the result before returning it.
*/
func simulateStoneCount(number int, blinks int, cache *stoneCache) int {
	// Check cache first
	if cachedResult, found := cache.get(number, blinks); found {
		return cachedResult
	}

	// Base case: no more blinks
	if blinks == 0 {
		cache.set(number, blinks, 1)
		return 1
	}

	var stoneCount int
	switch {
	case number == 0:
		// Rule 1: 0 becomes 1
		stoneCount = simulateStoneCount(1, blinks-1, cache)
	case len(strconv.Itoa(number))%2 == 0:
		// Rule 2: Even digit count - split the number
		digits := strconv.Itoa(number)
		half := len(digits) / 2
		left, _ := strconv.Atoi(digits[:half])
		right, _ := strconv.Atoi(digits[half:])

		// Recursively count stones for both halves
		leftCount := simulateStoneCount(left, blinks-1, cache)
		rightCount := simulateStoneCount(right, blinks-1, cache)

		stoneCount = leftCount + rightCount
	default:
		// Rule 3: Multiply by 2024
		stoneCount = simulateStoneCount(number*2024, blinks-1, cache)
	}

	// Cache and return the result
	cache.set(number, blinks, stoneCount)
	return stoneCount
}

func simulateBlinks(stones []int, blinks int) int {
	cache := newStoneCache()
	totalStones := 0

	// Process each stone independently and sum stone counts
	for _, stone := range stones {
		totalStones += simulateStoneCount(stone, blinks, cache)
	}

	return totalStones
}

func main() {
	input := []int{0, 7, 6618216, 26481, 885, 42, 202642, 8791}
	blink_runs := []int{25, 75}

	for _, blinks := range blink_runs {
		stoneCount := simulateBlinks(input, blinks)
		fmt.Println("Number of stones after", blinks, "blinks:", stoneCount)
	}
}
