package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// represents a digital logic circuit with gates and wires
type Circuit struct {
	wireValues map[string]int
	gates      []Gate // Changed from map[string][]Gate to []Gate
}

// represents a logic gate with its connections
type Gate struct {
	input1    string
	operation string
	input2    string
	output    string
}

// creates a new Circuit from input lines
func NewCircuit(input []string) *Circuit {
	wireValues, gates := parseInput(input)
	return &Circuit{
		wireValues: wireValues,
		gates:      gates,
	}
}

// splits the input into initial wire values and gate connections
func parseInput(input []string) (map[string]int, []Gate) {
	wireValues := make(map[string]int)
	var gates []Gate

	isInitialState := true
	for _, line := range input {
		if line == "" {
			isInitialState = false
			continue
		}

		if isInitialState {
			wire, value := parseWireValue(line)
			wireValues[wire] = value
			continue
		}

		if len(line) > 3 {
			gate := parseGate(line)
			gates = append(gates, gate)
		}
	}

	return wireValues, gates
}

// extracts wire name and value from a line like "x00: 1"
func parseWireValue(line string) (string, int) {
	parts := strings.Split(line, ":")
	wire := strings.TrimSpace(parts[0])
	value, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return wire, value
}

// creates a Gate from a line like "x00 AND y00 -> z00"
func parseGate(line string) Gate {
	parts := strings.Split(line, " ")
	return Gate{
		input1:    parts[0],
		operation: parts[1],
		input2:    parts[2],
		output:    parts[4],
	}
}

// evaluates all gates in the circuit until all outputs are stable
func (c *Circuit) Simulate() {
	queue := make([]string, 0, len(c.wireValues))
	for wire := range c.wireValues {
		queue = append(queue, wire)
	}

	evaluated := make(map[Gate]bool)
	visited := make(map[string]bool)

	for len(queue) > 0 {
		wire := queue[0]
		queue = queue[1:]

		if visited[wire] || !c.hasValue(wire) {
			continue
		}

		// Find gates that use this wire as input
		for _, gate := range c.gates {
			if (gate.input1 == wire || gate.input2 == wire) && !evaluated[gate] {
				if c.evaluateGate(gate) {
					queue = append(queue, gate.output)
					evaluated[gate] = true
				}
			}
		}

		visited[wire] = true
	}
}

// checks if a wire has a value assigned
func (c *Circuit) hasValue(wire string) bool {
	_, exists := c.wireValues[wire]
	return exists
}

// applies the gate operation if both inputs are available
func (c *Circuit) evaluateGate(gate Gate) bool {
	input1Val, hasInput1 := c.wireValues[gate.input1]
	input2Val, hasInput2 := c.wireValues[gate.input2]

	if !hasInput1 || !hasInput2 {
		return false
	}

	switch gate.operation {
	case "AND":
		c.wireValues[gate.output] = input1Val & input2Val
	case "OR":
		c.wireValues[gate.output] = input1Val | input2Val
	case "XOR":
		c.wireValues[gate.output] = input1Val ^ input2Val
	}
	return true
}

// converts binary outputs starting with prefix to decimal
func (c *Circuit) GetDecimalOutput(prefix string) int64 {
	wires := c.getWiresWithPrefix(prefix)
	sort.Strings(wires)

	binary := make([]string, len(wires))
	for i := len(wires) - 1; i >= 0; i-- {
		binary = append(binary, fmt.Sprintf("%d", c.wireValues[wires[i]]))
	}

	result, _ := strconv.ParseInt(strings.Join(binary, ""), 2, 64)
	return result
}

// returns all wire names starting with the given prefix
func (c *Circuit) getWiresWithPrefix(prefix string) []string {
	var wires []string
	for wire := range c.wireValues {
		if strings.HasPrefix(wire, prefix) {
			wires = append(wires, wire)
		}
	}
	return wires
}

// checks if the circuit follows adder rules
func (c *Circuit) ValidateRippleCarryAdder() string {
	var faultyGates []string
	gateOutputs := make(map[string]bool) // Track outputs we've seen

	for _, gate := range c.gates {
		if c.isGateFaulty(gate) && !gateOutputs[gate.output] {
			faultyGates = append(faultyGates, gate.output)
			gateOutputs[gate.output] = true
		}
	}

	sort.Strings(faultyGates)
	return strings.Join(faultyGates, ",")
}

// checks if a gate violates ripple-carry adder rules
func (c *Circuit) isGateFaulty(gate Gate) bool {
	// Rule 1: z-output must use XOR except for last bit
	if strings.HasPrefix(gate.output, "z") && gate.output != "z45" && gate.operation != "XOR" {
		return true
	}

	// Rule 2: Non-z output with non-x/y inputs must not use XOR
	if !strings.HasPrefix(gate.output, "z") &&
		!isInputWire(gate.input1) && !isInputWire(gate.input2) &&
		gate.operation == "XOR" {
		return true
	}

	// Rule 3: XOR gates with x/y inputs need XOR consumer
	if gate.operation == "XOR" && isInputPair(gate.input1, gate.input2) &&
		!isFirstInputPair(gate.input1, gate.input2) { // Added check for first input pair
		if !c.hasConsumerGate(gate.output, "XOR") {
			return true
		}
	}

	// Rule 4: AND gates with x/y inputs need OR consumer
	if gate.operation == "AND" && isInputPair(gate.input1, gate.input2) &&
		!isFirstInputPair(gate.input1, gate.input2) { // Added check for first input pair
		if !c.hasConsumerGate(gate.output, "OR") {
			return true
		}
	}

	return false
}

// Helper functions for gate validation
func isInputWire(wire string) bool {
	return strings.HasPrefix(wire, "x") || strings.HasPrefix(wire, "y")
}

func isInputPair(wire1, wire2 string) bool {
	return (strings.HasPrefix(wire1, "x") && strings.HasPrefix(wire2, "y")) ||
		(strings.HasPrefix(wire1, "y") && strings.HasPrefix(wire2, "x"))
}

// isFirstInputPair checks if both inputs are x00/y00
func isFirstInputPair(wire1, wire2 string) bool {
	return (wire1 == "x00" || wire1 == "y00") && (wire2 == "x00" || wire2 == "y00")
}

// hasConsumerGate checks if a wire is used as input to a gate with the specified operation
func (c *Circuit) hasConsumerGate(wire, operation string) bool {
	for _, gate := range c.gates {
		if (gate.input1 == wire || gate.input2 == wire) && gate.operation == operation {
			return true
		}
	}
	return false
}

func main() {
	input := readLines("input.txt")
	circuit := NewCircuit(input)

	circuit.Simulate()
	part1 := circuit.GetDecimalOutput("z")
	fmt.Printf("Part 1: %d\n", part1)

	part2 := circuit.ValidateRippleCarryAdder()
	fmt.Printf("Part 2: %s\n", part2)
}

func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 512*1024)
	scanner.Buffer(buf, 512*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}
