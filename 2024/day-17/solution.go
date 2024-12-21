package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// Register names for the 3-bit computer
const (
	RegisterA = 4
	RegisterB = 5
	RegisterC = 6
)

// Operation codes for the 3-bit computer
const (
	OpADV = iota // Divide A by 2^operand
	OpBXL        // XOR B with literal operand
	OpBST        // Set B to combo operand mod 8
	OpJNZ        // Jump if A is not zero
	OpBXC        // XOR B with C
	OpOUT        // Output combo operand mod 8
	OpBDV        // Divide A by 2^operand, store in B
	OpCDV        // Divide A by 2^operand, store in C
)

// Computer represents the state of our 3-bit computer
type Computer struct {
	A, B, C int   // Registers
	Code    []int // Program
	IP      int   // Instruction Pointer
	Output  []int // Program output
}

// NewComputer creates a new computer with initial register values and program
func NewComputer(a, b, c int, code []int) *Computer {
	return &Computer{
		A:      a,
		B:      b,
		C:      c,
		Code:   code,
		IP:     0,
		Output: make([]int, 0),
	}
}

// getComboValue returns the value of a combo operand based on the rules
func (c *Computer) getComboValue(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case RegisterA:
		return c.A
	case RegisterB:
		return c.B
	case RegisterC:
		return c.C
	default:
		return 0 // Handle reserved value 7
	}
}

// executeInstruction executes a single instruction and returns whether to continue execution
func (c *Computer) executeInstruction() bool {
	if c.IP >= len(c.Code)-1 { // Check if we can read both opcode and operand
		return false
	}

	opcode := c.Code[c.IP]
	operand := c.Code[c.IP+1]
	comboValue := c.getComboValue(operand)

	switch opcode {
	case OpADV:
		c.A = c.A >> comboValue
	case OpBXL:
		c.B ^= operand
	case OpBST:
		c.B = comboValue % 8
	case OpJNZ:
		if c.A != 0 {
			c.IP = operand
			return true
		}
	case OpBXC:
		c.B ^= c.C
	case OpOUT:
		c.Output = append(c.Output, comboValue%8)
	case OpBDV:
		c.B = c.A >> comboValue
	case OpCDV:
		c.C = c.A >> comboValue
	}

	c.IP += 2
	return true
}

// Run executes the program until completion and returns the output
func (c *Computer) Run() []int {
	for c.executeInstruction() {
	}
	return c.Output
}

// ProgramInput represents the parsed input file
type ProgramInput struct {
	A, B, C int
	Code    []int
}

// parseInput reads and parses the input file
func parseInput(r io.Reader) (*ProgramInput, error) {
	scanner := bufio.NewScanner(r)
	registers := make([]int, 3)

	// Parse register values
	for i := 0; i < 3; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("failed to read register %d", i)
		}
		val, err := strconv.Atoi(strings.Fields(scanner.Text())[2])
		if err != nil {
			return nil, fmt.Errorf("invalid register value: %v", err)
		}
		registers[i] = val
	}

	// Skip empty line
	scanner.Scan()

	// Parse program code
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read program code")
	}
	var code []int
	for _, s := range strings.Split(strings.Fields(scanner.Text())[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid program code: %v", err)
		}
		code = append(code, n)
	}

	return &ProgramInput{
		A:    registers[0],
		B:    registers[1],
		C:    registers[2],
		Code: code,
	}, scanner.Err()
}

// findSelfReplicatingValue finds the lowest value of A that makes the program output itself
func findSelfReplicatingValue(b, c int, code []int) int {
	a := 0
	for i := len(code) - 1; i >= 0; i-- {
		// make room for 3 bits by shifting to the left
		a <<= 3

		// Incrementally try values until we find the right one
		for {
			comp := NewComputer(a, b, c, code)
			if slices.Equal(comp.Run(), code[i:]) {
				break
			}
			a++
		}
	}
	return a
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	input, err := parseInput(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	// Part 1: Run the program with initial values
	computer := NewComputer(input.A, input.B, input.C, input.Code)
	output := computer.Run()
	fmt.Printf("Part 1: %s\n", strings.Join(formatOutput(output), ","))

	// Part 2: Find the self-replicating value
	result := findSelfReplicatingValue(input.B, input.C, input.Code)
	fmt.Printf("Part 2: %d\n", result)
}

// formatOutput converts a slice of ints to a slice of strings
func formatOutput(nums []int) []string {
	result := make([]string, len(nums))
	for i, num := range nums {
		result[i] = strconv.Itoa(num)
	}
	return result
}
