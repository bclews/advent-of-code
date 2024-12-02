package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

// parseFile reads the contents of a file and returns two slices of integers
func parseFile(r io.Reader) ([]int, []int, error) {
	column1 := make([]int, 0, 1000)
	column2 := make([]int, 0, 1000)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid line format: %s", scanner.Text())
		}

		num1, err1 := strconv.Atoi(parts[0])
		num2, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("error parsing numbers: %v, %v", err1, err2)
		}

		column1 = append(column1, num1)
		column2 = append(column2, num2)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return column1, column2, nil
}

func sortColumns(col1, col2 []int) ([]int, []int) {
	sort.Ints(col1)
	sort.Ints(col2)
	return col1, col2
}

// Pair represents a pair of numbers with their distance
type Pair struct {
	Left     int
	Right    int
	Distance int
}

// calculatePairedDistance pairs sorted lists and calculates total distance
func calculatePairedDistance(col1, col2 []int) (int, []Pair) {
	// Ensure lists are the same length
	if len(col1) != len(col2) {
		return 0, nil
	}

	// Total distance to track
	totalDistance := 0

	// Store the pairs
	pairs := make([]Pair, len(col1))

	// Calculate distances between paired numbers
	for i := 0; i < len(col1); i++ {
		// Calculate absolute distance between paired numbers
		distance := abs(col1[i] - col2[i])
		totalDistance += distance

		// Store the pair
		pairs[i] = Pair{
			Left:     col1[i],
			Right:    col2[i],
			Distance: distance,
		}
	}

	return totalDistance, pairs
}

func abs(x int) int {
	mask := x >> (bits.UintSize - 1)
	return (x ^ mask) - mask
}

func calculateSimilarityScore(col1, col2 []int) int {
	// Count occurrences of each number in col2
	rightListCounts := make(map[int]int)
	for _, num := range col2 {
		rightListCounts[num]++
	}

	// Calculate similarity score directly while iterating through col1
	similarityScore := 0
	for _, num := range col1 {
		if count, exists := rightListCounts[num]; exists {
			similarityScore += num * count
		}
	}

	return similarityScore
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
	column1, column2, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Sort the columns
	sortedCol1, sortedCol2 := sortColumns(column1, column2)

	// Print the sorted lists
	fmt.Println("Sorted Column 1:", sortedCol1)
	fmt.Println("Sorted Column 2:", sortedCol2)

	// Calculate paired distances
	totalDistance, pairs := calculatePairedDistance(sortedCol1, sortedCol2)

	// Print paired results
	fmt.Println("\nPaired Numbers and Their Distances:")
	for _, pair := range pairs {
		fmt.Printf("(%d, %d) Distance: %d\n", pair.Left, pair.Right, pair.Distance)
	}
	fmt.Println("\nTotal Distance:", totalDistance)

	// Calculate similarity score
	similarityScore := calculateSimilarityScore(column1, column2)
	fmt.Println("\nSimilarity Score:", similarityScore)
}
