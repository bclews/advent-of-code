package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Represents a coordinate on a keypad
type Position struct {
	row, col int
}

// Represents the layout of buttons on either a numeric or directional keypad
type Keypad struct {
	buttonPositions map[byte]Position
	emptyPosition   Position
	startPosition   Position
}

// Creates a keypad with the numeric layout (789/456/123/_0A)
func NewNumericKeypad() *Keypad {
	return &Keypad{
		buttonPositions: map[byte]Position{
			'7': {0, 0},
			'8': {0, 1},
			'9': {0, 2},

			'4': {1, 0},
			'5': {1, 1},
			'6': {1, 2},

			'1': {2, 0},
			'2': {2, 1},
			'3': {2, 2},

			'_': {3, 0},
			'0': {3, 1},
			'A': {3, 2},
		},
		emptyPosition: Position{3, 0},
		startPosition: Position{3, 2}, // 'A' position
	}
}

// NewDirectionalKeypad creates a keypad with the directional layout (_^A/<v>)
func NewDirectionalKeypad() *Keypad {
	return &Keypad{
		buttonPositions: map[byte]Position{
			'_': {0, 0},
			'^': {0, 1},
			'A': {0, 2},

			'<': {1, 0},
			'v': {1, 1},
			'>': {1, 2},
		},
		emptyPosition: Position{0, 0},
		startPosition: Position{0, 2}, // 'A' position
	}
}

// Represents the unique identifier for memoization
type cacheKey struct {
	sequence string
	depth    int
}

// Represents a chain of robots pressing buttons on keypads
type RobotChain struct {
	numericKeypad     *Keypad
	directionalKeypad *Keypad
	memoCache         map[cacheKey]int
}

// Creates a new instance of RobotChain
func NewRobotChain() *RobotChain {
	return &RobotChain{
		numericKeypad:     NewNumericKeypad(),
		directionalKeypad: NewDirectionalKeypad(),
		memoCache:         make(map[cacheKey]int),
	}
}

// Generates a sequence of moves for a specific keypad
func (r *RobotChain) generateKeypadSequence(input string, keypad *Keypad) string {
	var sequence strings.Builder
	currentPos := keypad.startPosition

	for _, char := range input {
		targetPos := keypad.buttonPositions[byte(char)]

		horizontal, vertical := r.calculateMoves(currentPos, targetPos)
		moves := r.optimizeMoveOrder(currentPos, targetPos, keypad.emptyPosition, horizontal, vertical)
		sequence.WriteString(moves)
		sequence.WriteByte('A')

		currentPos = targetPos
	}

	return sequence.String()
}

// Determines horizontal and vertical moves needed
func (r *RobotChain) calculateMoves(from, to Position) (string, string) {
	var horizontal, vertical string

	// Horizontal moves
	deltaCol := to.col - from.col
	for deltaCol > 0 {
		horizontal += ">"
		deltaCol--
	}
	for deltaCol < 0 {
		horizontal += "<"
		deltaCol++
	}

	// Vertical moves
	deltaRow := to.row - from.row
	for deltaRow > 0 {
		vertical += "v"
		deltaRow--
	}
	for deltaRow < 0 {
		vertical += "^"
		deltaRow++
	}

	return horizontal, vertical
}

// Determines the best order of moves to avoid the empty position
func (r *RobotChain) optimizeMoveOrder(from, to, empty Position, horizontal, vertical string) string {
	switch {
	case from.row == empty.row && to.col == empty.col:
		return vertical + horizontal
	case from.col == empty.col && to.row == empty.row:
		return horizontal + vertical
	case strings.Contains(horizontal, "<"):
		return horizontal + vertical
	default:
		return vertical + horizontal
	}
}

// Generates the complete sequence through the robot chain
func (r *RobotChain) GenerateButtonSequence(code string, robotDepth int) (int, error) {
	firstSequence := r.generateKeypadSequence(code, r.numericKeypad)
	return r.calculateSequenceLength(firstSequence, robotDepth)
}

// Calculates the total length accounting for robot chain depth
func (r *RobotChain) calculateSequenceLength(sequence string, depth int) (int, error) {
	if depth == 0 {
		return len(sequence), nil
	}

	key := cacheKey{sequence, depth}
	if cachedLen, exists := r.memoCache[key]; exists {
		return cachedLen, nil
	}

	totalLen := 0
	parts := strings.SplitAfter(sequence, "A")
	for _, part := range parts {
		if part == "" {
			continue
		}

		dirSequence := r.generateKeypadSequence(part, r.directionalKeypad)
		length, err := r.calculateSequenceLength(dirSequence, depth-1)
		if err != nil {
			return 0, err
		}
		totalLen += length
	}

	r.memoCache[key] = totalLen
	return totalLen, nil
}

// Extracts and parses the numeric portion of a code
func ParseNumericPart(code string) (int, error) {
	numericPart := strings.TrimRight(code, "A")
	return strconv.Atoi(numericPart)
}

func calculateComplexitySum(codes []string, robotChain *RobotChain, depth int) int {
	sum := 0
	for _, code := range codes {
		seqLength, err := robotChain.GenerateButtonSequence(code, depth)
		if err != nil {
			fmt.Printf("Error processing code %s: %v\n", code, err)
			continue
		}

		numericPart, err := ParseNumericPart(code)
		if err != nil {
			fmt.Printf("Error parsing numeric part of code %s: %v\n", code, err)
			continue
		}

		sum += seqLength * numericPart
	}
	return sum
}

// parseFile reads the contents of a file and returns two slices of integers
func parseFile(r io.Reader) ([]string, error) {
	var codes []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return nil, err
	}
	return codes, nil
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	codes, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	robotChain := NewRobotChain()

	fmt.Println("Part One:", calculateComplexitySum(codes, robotChain, 2))
	fmt.Println("Part Two:", calculateComplexitySum(codes, robotChain, 25))
}
