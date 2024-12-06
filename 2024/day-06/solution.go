package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Direction represents the cardinal directions
type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func parseFile(r io.Reader) ([]string, error) {
	// Slice to store the lines
	var lines []string

	// Read the file line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), `" \t`) // Trim quotes, spaces, and tabs
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

// Map represents the laboratory map
type Map struct {
	grid       [][]rune
	GuardX     int
	GuardY     int
	GuardDir   Direction
	visitedPos map[string]bool
}

// NewMap creates a new Map from input strings
func NewMap(mapInput []string) *Map {
	m := &Map{
		grid:       make([][]rune, len(mapInput)),
		visitedPos: make(map[string]bool),
	}

	// Convert input to 2D rune grid and find guard's initial position
	for y, row := range mapInput {
		m.grid[y] = []rune(row)
		for x, cell := range row {
			switch cell {
			case '^':
				m.GuardX = x
				m.GuardY = y
				m.GuardDir = Up
			}
		}
	}

	// Mark starting position as visited
	m.markVisited(m.GuardX, m.GuardY)

	return m
}

// isObstacleAhead checks if there's an obstacle in the guard's current direction
func (m *Map) isObstacleAhead() bool {
	nextX, nextY := m.getNextPosition()
	return m.isOutOfBounds(nextX, nextY) || m.grid[nextY][nextX] == '#'
}

// getNextPosition calculates the next position based on current direction
func (m *Map) getNextPosition() (int, int) {
	switch m.GuardDir {
	case Up:
		return m.GuardX, m.GuardY - 1
	case Right:
		return m.GuardX + 1, m.GuardY
	case Down:
		return m.GuardX, m.GuardY + 1
	case Left:
		return m.GuardX - 1, m.GuardY
	}
	return m.GuardX, m.GuardY
}

// isOutOfBounds checks if the given coordinates are outside the map
func (m *Map) isOutOfBounds(x, y int) bool {
	return x < 0 || y < 0 || y >= len(m.grid) || x >= len(m.grid[y])
}

// markVisited adds the current position to visited positions
func (m *Map) markVisited(x, y int) {
	m.visitedPos[positionKey(x, y)] = true
}

// turnRight rotates the guard's direction 90 degrees clockwise
func (m *Map) turnRight() {
	m.GuardDir = (m.GuardDir + 1) % 4
}

// move advances the guard one step in the current direction
func (m *Map) move() {
	m.GuardX, m.GuardY = m.getNextPosition()
	m.markVisited(m.GuardX, m.GuardY)
}

// SimulateGuardPatrol runs the guard's patrol simulation
func (m *Map) SimulateGuardPatrol() int {
	for {
		// Check if the next move is possible
		nextX, nextY := m.getNextPosition()

		// Terminate if next move would be out of bounds
		if m.isOutOfBounds(nextX, nextY) {
			break
		}

		// If obstacle ahead, turn right
		if m.grid[nextY][nextX] == '#' {
			m.turnRight()
			continue
		}

		// Move forward
		m.move()
	}

	return len(m.visitedPos)
}

// Helper function to create a unique key for visited positions
func positionKey(x, y int) string {
	return string(rune(x)) + "," + string(rune(y))
}

func (m *Map) FindLoopObstructionPositions() []string {
	loopPositions := []string{}

	for y := 0; y < len(m.grid); y++ {
		row := m.grid[y] // Ensure the row exists
		for x := 0; x < len(row); x++ {

			if m.grid[y][x] == '#' {
				continue
			}

			// Temporarily place an obstruction
			originalCell := m.grid[y][x]
			m.grid[y][x] = '#'

			// Check for a loop
			if m.simulateAndDetectLoop() {
				loopPositions = append(loopPositions, positionKey(x, y))
				// fmt.Printf("Loop detected with obstruction at (%d, %d)\n", x, y)
			}

			// Restore the original map cell
			m.grid[y][x] = originalCell
		}
	}

	return loopPositions
}

func (m *Map) simulateAndDetectLoop() bool {
	visitedStates := map[string]bool{}
	m.GuardX, m.GuardY, m.GuardDir = m.findGuardStartingPosition()

	for {
		stateKey := fmt.Sprintf("%d,%d,%d", m.GuardX, m.GuardY, m.GuardDir)

		if visitedStates[stateKey] {
			return true
		}

		visitedStates[stateKey] = true

		nextX, nextY := m.getNextPosition()

		if m.isOutOfBounds(nextX, nextY) {
			break
		}

		if m.grid[nextY][nextX] == '#' {
			m.turnRight()
			continue
		}

		m.move()
	}

	return false
}

// findGuardStartingPosition resets and finds the guard's starting position and direction
func (m *Map) findGuardStartingPosition() (int, int, Direction) {
	for y := 0; y < len(m.grid); y++ {
		for x := 0; x < len(m.grid[y]); x++ {
			switch m.grid[y][x] {
			case '^':
				return x, y, Up
			}
		}
	}
	return 0, 0, Up // Default fallback (shouldn't happen with valid maps)
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	lines, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Create the map
	m := NewMap(lines)

	// Simulate the guard's patrol
	visitedPositions := m.SimulateGuardPatrol()

	// Print the results
	fmt.Printf("Total unique positions visited: %d\n", visitedPositions)
	fmt.Printf("Guard's final position: (%d, %d)\n", m.GuardX, m.GuardY)
	fmt.Printf("Guard's final direction: %v\n\n", m.GuardDir)

	// Find obstruction positions
	loopPositions := m.FindLoopObstructionPositions()

	// Print results
	fmt.Printf("Total positions to cause a loop: %d\n", len(loopPositions))
}
