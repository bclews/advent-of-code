package main

import (
	"fmt"
	"strconv"
	"sync"
)

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

func (sc *stoneCache) get(number, blinks int) (int, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	result, exists := sc.cache[cacheKey{number, blinks}]
	return result, exists
}

func (sc *stoneCache) set(number, blinks, result int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.cache[cacheKey{number, blinks}] = result
}

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
