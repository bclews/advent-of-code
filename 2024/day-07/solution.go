package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// Helper function to evaluate a sequence with given operators
func evaluate(numbers []int, operators []rune) int {
	result := numbers[0]
	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case '+':
			result += numbers[i+1]
		case '*':
			result *= numbers[i+1]
		case '|': // Concatenation operator
			result, _ = strconv.Atoi(fmt.Sprintf("%d%d", result, numbers[i+1]))

		}
	}
	return result
}

func generateOperatorCombinations(n int) [][]rune {
	combinations := [][]rune{}
	operators := []rune{'+', '*', '|'}    // Added concatenation operator
	total := int(math.Pow(3, float64(n))) // 3^n combinations
	for i := 0; i < total; i++ {
		combination := []rune{}
		temp := i
		for j := 0; j < n; j++ {
			opIndex := temp % 3
			combination = append(combination, operators[opIndex])
			temp /= 3
		}
		combinations = append(combinations, combination)
	}
	return combinations
}

// Function to determine solvable equations and calculate the total calibration result
func calculateCalibrationResult(equations map[int][]int) int {
	totalSum := 0

	for testValue, numbers := range equations {
		n := len(numbers) - 1 // Number of operator positions
		operators := generateOperatorCombinations(n)
		isSolvable := false

		// Check all operator combinations
		for _, ops := range operators {
			if evaluate(numbers, ops) == testValue {
				isSolvable = true
				break
			}
		}

		// Add to total sum if solvable
		if isSolvable {
			totalSum += testValue
		}
	}

	return totalSum
}

func parseInputFile(r io.Reader) (map[int][]int, error) {
	equations := make(map[int][]int)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into test value and numbers
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		// Parse the test value
		testValue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid test value: %s", parts[0])
		}

		// Parse the numbers
		numberStrings := strings.Fields(strings.TrimSpace(parts[1]))
		numbers := make([]int, len(numberStrings))
		for i, numStr := range numberStrings {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", numStr)
			}
			numbers[i] = num
		}

		// Add to the map
		equations[testValue] = numbers
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return equations, nil
}

// Main function
func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the input file
	equations, err := parseInputFile(file)
	if err != nil {
		fmt.Println("Error parsing input file:", err)
		return
	}

	result := calculateCalibrationResult(equations)
	fmt.Printf("Total Calibration Result: %d\n", result)
}
