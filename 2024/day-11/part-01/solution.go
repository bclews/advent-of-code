package main

import (
	"fmt"
	"strconv"
)

// simulateBlinks simulates the evolution of stones over a given number of blinks.
func simulateBlinks(stones []int, blinks int) []int {
	for i := 0; i < blinks; i++ {
		stones = blink(stones)
	}
	return stones
}

// blink applies the transformation rules to the stones.
func blink(stones []int) []int {
	var newStones []int
	for _, stone := range stones {
		switch {
		case stone == 0:
			newStones = append(newStones, 1)
		case len(strconv.Itoa(stone))%2 == 0:
			split := splitStone(stone)
			newStones = append(newStones, split...)
		default:
			newStones = append(newStones, stone*2024)
		}
	}
	return newStones
}

// splitStone splits a stone with an even number of digits into two stones.
func splitStone(stone int) []int {
	digits := strconv.Itoa(stone)
	half := len(digits) / 2
	left, _ := strconv.Atoi(digits[:half])
	right, _ := strconv.Atoi(digits[half:])
	return []int{left, right}
}

func main() {
	input := []int{0, 7, 6618216, 26481, 885, 42, 202642, 8791}
	blinks := 25

	stones := simulateBlinks(input, blinks)
	fmt.Println("Number of stones after", blinks, "blinks:", len(stones))
}
