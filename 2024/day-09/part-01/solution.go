package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func parseFile(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// parseLine processes a string of digits and returns a slice of integers based on specific rules.
func parseLine(line string) []int {
	totalSize := calculateTotalSize(line)
	result := make([]int, totalSize)

	populateRegister(line, result)
	return result
}

// calculateTotalSize computes the total size of the output slice.
func calculateTotalSize(line string) int {
	totalSize := 0
	for _, char := range line {
		totalSize += parseDigit(char)
	}
	return totalSize
}

// parseDigit converts a rune representing a digit to an integer.
func parseDigit(char rune) int {
	return int(char - '0')
}

// populateRegister fills the output slice based on the input string's rules.
func populateRegister(line string, result []int) {
	id := 0
	ptr := 0

	for i, char := range line {
		size := parseDigit(char)

		if i%2 == 0 {
			// Even index: assign `id`.
			fillSlice(result[ptr:ptr+size], id)
			id++
		} else {
			// Odd index: assign `-1`.
			fillSlice(result[ptr:ptr+size], -1)
		}

		ptr += size
	}
}

// fillSlice assigns a given value to all elements of a slice.
func fillSlice(slice []int, value int) {
	for i := range slice {
		slice[i] = value
	}
}

// checksum calculates the weighted sum of the values in the register.
func checksum(register []int) int {
	sum := 0
	for index, value := range register {
		sum += index * value
	}
	return sum
}

func solve(input []string) (output string) {
	register := parseLine(input[0])

	// Start pointer at the last element of the register
	ptr := len(register) - 1

	// Iterate through the register
	for i, id := range register {
		// If current element is -1 (needs to be replaced)
		if id == -1 {
			// Find the last non-negative element from the end
			for ptr > i && register[ptr] == -1 {
				ptr--
			}

			// If no more elements to swap, exit the loop
			if ptr <= i {
				break
			}

			// Swap the current -1 element with the last non-negative element
			register[i], register[ptr] = register[ptr], register[i]
		}
	}

	// Trim the register to remove trailing elements
	register = register[:ptr]

	// Convert the checksum of the modified register to a string
	return strconv.Itoa(checksum(register))
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// parse the file
	input, err := parseFile(file)

	solution := solve(input)
	fmt.Println(solution)
}
