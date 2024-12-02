package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// parseFile reads integers from the input file into a 2D slice
func parseFile(r io.Reader) ([][]int, error) {
	var result [][]int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		numberStrings := strings.Fields(line)

		var lineInts []int

		// Convert each string to an integer
		for _, numStr := range numberStrings {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing integer: %v", err)
			}
			lineInts = append(lineInts, num)
		}

		result = append(result, lineInts)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return result, nil
}

func IsReportSafe(report []int) bool {
	// Handle edge cases
	if len(report) < 2 {
		return false
	}

	// If first two elements are equal, it's unsafe
	if report[0] == report[1] {
		return false
	}

	// Determine initial direction
	isIncreasing := report[0] < report[1]
	isDecreasing := report[0] > report[1]

	// Check each adjacent pair
	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]

		// Check direction consistency
		if isIncreasing && diff <= 0 {
			return false
		}
		if isDecreasing && diff >= 0 {
			return false
		}

		// Check difference range
		absDiff := abs(diff)
		if absDiff < 1 || absDiff > 3 {
			return false
		}
	}

	return true
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CountSafeReports(reports [][]int) int {
	safeReportCount := 0

	for _, report := range reports {
		if IsReportSafe(report) {
			safeReportCount++
		}
	}

	return safeReportCount
}

// Optional main function for demonstration
func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the file
	unusualData, err := parseFile(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Count the number of safe reports
	safeReportCount := CountSafeReports(unusualData)
	fmt.Println("Number of safe reports:", safeReportCount)
}
