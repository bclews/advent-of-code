package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// GameState represents the current state of the puzzle
type GameState struct {
	Grid    Grid
	Robot   Position
	MaxRows int
	MaxCols int
}

// Position represents a point on the grid
type Position struct {
	Row, Col int
}

// Direction represents a movement vector
type Direction struct {
	DRow, DCol int
}

// Grid is a wrapper around the map representation of the game board
type Grid map[Position]rune

// Constants for game elements
const (
	Empty = '.'
	Wall  = '#'
	BoxL  = '[' // Left side of box
	BoxR  = ']' // Right side of box
	Robot = '@'
)

// Movement directions
var directions = map[byte]Direction{
	'<': {0, -1},
	'>': {0, 1},
	'^': {-1, 0},
	'v': {1, 0},
}

// NewGameState creates and initializes a new game state
func NewGameState(r io.Reader) (*GameState, []byte, error) {
	grid, moves, robotPos, err := parseInput(r)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing input: %w", err)
	}

	maxRows, maxCols := calculateGridDimensions(grid)

	return &GameState{
		Grid:    grid,
		Robot:   robotPos,
		MaxRows: maxRows,
		MaxCols: maxCols,
	}, moves, nil
}

// parseInput reads the puzzle input and returns the initial game state
func parseInput(r io.Reader) (Grid, []byte, Position, error) {
	scanner := bufio.NewScanner(r)
	tempGrid := make(Grid)
	var robotPos Position
	var maxCol int

	// Parse the grid
	for row := 0; scanner.Scan(); row++ {
		line := scanner.Text()
		if line == "" {
			break
		}

		if len(line) > maxCol {
			maxCol = len(line)
		}

		for col, char := range line {
			pos := Position{row, col}
			if char == Robot {
				robotPos = pos
				tempGrid[pos] = Empty
			} else {
				tempGrid[pos] = char
			}
		}
	}

	// Parse the moves
	var moves []byte
	for scanner.Scan() {
		moves = append(moves, []byte(scanner.Text())...)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, Position{}, fmt.Errorf("scanning input: %w", err)
	}

	expandedGrid := expandGrid(tempGrid)
	robotPos.Col *= 2 // Adjust for expanded grid

	return expandedGrid, moves, robotPos, nil
}

// expandGrid doubles the grid horizontally to handle box movements
func expandGrid(input Grid) Grid {
	expanded := make(Grid)
	maxRow, maxCol := calculateGridDimensions(input)

	for row := 0; row <= maxRow; row++ {
		for col := 0; col <= maxCol*2+1; col++ {
			pos := Position{row, col}
			originalCol := col / 2
			originalPos := Position{row, originalCol}

			if char, exists := input[originalPos]; exists {
				expanded[pos] = expandCharacter(char, col%2 == 0)
			} else {
				expanded[pos] = Empty
			}
		}
	}
	return expanded
}

// expandCharacter handles the conversion of characters during grid expansion
func expandCharacter(char rune, isEven bool) rune {
	switch char {
	case Wall:
		return Wall
	case 'O': // Original box character
		if isEven {
			return BoxL
		}
		return BoxR
	default:
		return Empty
	}
}

// calculateGridDimensions returns the maximum row and column values in the grid
func calculateGridDimensions(grid Grid) (maxRow, maxCol int) {
	for pos := range grid {
		if pos.Row > maxRow {
			maxRow = pos.Row
		}
		if pos.Col > maxCol {
			maxCol = pos.Col
		}
	}
	return
}

// Print displays the current state of the game
func (g *GameState) Print() {
	for row := 0; row <= g.MaxRows; row++ {
		for col := 0; col <= g.MaxCols; col++ {
			pos := Position{row, col}
			if pos == g.Robot {
				fmt.Print(string(Robot))
			} else {
				char, exists := g.Grid[pos]
				if !exists {
					char = Empty
				}
				fmt.Print(string(char))
			}
		}
		fmt.Println()
	}
}

// Move attempts to move the robot in the specified direction
func (g *GameState) Move(move byte) {
	dir, exists := directions[move]
	if !exists {
		return
	}

	if move == '<' || move == '>' {
		g.handleHorizontalMove(dir)
	} else {
		g.handleVerticalMove(dir)
	}
}

// handleHorizontalMove processes horizontal robot movements
func (g *GameState) handleHorizontalMove(dir Direction) {
	nextEmpty, ok := g.findNextEmpty(dir)
	if !ok {
		return
	}

	// Move boxes
	for curr := nextEmpty; curr != g.Robot; {
		prev := Position{curr.Row, curr.Col - dir.DCol}
		g.Grid[curr], g.Grid[prev] = g.Grid[prev], g.Grid[curr]
		curr = prev
	}

	// Move robot
	g.Robot.Col += dir.DCol
}

// handleVerticalMove processes vertical robot movements
func (g *GameState) handleVerticalMove(dir Direction) {
	affected, maxLevel, ok := g.findAffectedColumns(dir.DRow)
	if !ok {
		return
	}

	// Move boxes
	for row := maxLevel; row != g.Robot.Row; row -= dir.DRow {
		for col := range affected[row] {
			pos1 := Position{row + dir.DRow, col}
			pos2 := Position{row, col}
			g.Grid[pos1], g.Grid[pos2] = g.Grid[pos2], g.Grid[pos1]
		}
	}

	// Move robot
	g.Robot.Row += dir.DRow
}

// findNextEmpty finds the next empty position in the specified direction
func (g *GameState) findNextEmpty(dir Direction) (Position, bool) {
	pos := g.Robot
	for {
		pos = Position{pos.Row + dir.DRow, pos.Col + dir.DCol}
		switch g.Grid[pos] {
		case Empty:
			return pos, true
		case Wall:
			return Position{}, false
		}
	}
}

// findAffectedColumns determines which columns are affected by a vertical move
func (g *GameState) findAffectedColumns(deltaRow int) (map[int]map[int]struct{}, int, bool) {
	affected := map[int]map[int]struct{}{
		g.Robot.Row: {g.Robot.Col: struct{}{}},
	}

	for currRow := g.Robot.Row; ; currRow += deltaRow {
		newCols, ok := g.findNewColumns(currRow+deltaRow, affected[currRow])
		if !ok {
			return nil, 0, false
		}

		if len(newCols) == 0 {
			return affected, currRow, true
		}

		affected[currRow+deltaRow] = newCols
	}
}

// findNewColumns identifies new columns affected by box movement
func (g *GameState) findNewColumns(nextRow int, columns map[int]struct{}) (map[int]struct{}, bool) {
	newCols := make(map[int]struct{})
	for col := range columns {
		switch g.Grid[Position{nextRow, col}] {
		case Wall:
			return nil, false
		case BoxL:
			newCols[col] = struct{}{}
			newCols[col+1] = struct{}{}
		case BoxR:
			newCols[col] = struct{}{}
			newCols[col-1] = struct{}{}
		}
	}
	return newCols, true
}

// CalculateScore computes the final score based on box positions
func (g *GameState) CalculateScore() int {
	score := 0
	for pos, char := range g.Grid {
		if char == BoxL {
			score += 100*pos.Row + pos.Col
		}
	}
	return score
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	game, moves, err := NewGameState(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing game: %v\n", err)
		os.Exit(1)
	}

	// Execute all moves
	for _, move := range moves {
		game.Move(move)
	}

	// Print final state and score
	game.Print()
	fmt.Printf("Final score: %d\n", game.CalculateScore())
}
