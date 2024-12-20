package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
)

// Direction represents the cardinal directions
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Position represents a point in the maze
type Position struct {
	x, y int
}

type State struct {
	pos      Position
	dir      Direction
	cost     int
	index    int
	previous *State
}

// PreviousState stores information about how we reached a state
type PreviousState struct {
	pos Position
	dir Direction
}

// PriorityQueue implementation
type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*State)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// PathInfo stores information about a state in the optimal path
type PathInfo struct {
	next map[string]bool // Stores next possible states in optimal paths
	cost int
}

// StateKey creates a unique key for visited states
func stateKey(pos Position, dir Direction) string {
	return fmt.Sprintf("%d,%d,%d", pos.x, pos.y, dir)
}

// GetNextPosition returns the next position given current position and direction
func getNextPosition(pos Position, dir Direction) Position {
	switch dir {
	case North:
		return Position{pos.x, pos.y - 1}
	case East:
		return Position{pos.x + 1, pos.y}
	case South:
		return Position{pos.x, pos.y + 1}
	case West:
		return Position{pos.x - 1, pos.y}
	}
	return pos
}

// RotateClockwise returns the new direction after rotating 90° clockwise
func rotateClockwise(dir Direction) Direction {
	return (dir + 1) % 4
}

// RotateCounterClockwise returns the new direction after rotating 90° counterclockwise
func rotateCounterClockwise(dir Direction) Direction {
	return (dir + 3) % 4
}

func FindLowestScoreWithPaths(maze []string) (int, map[Position]bool) {
	// Find start and end positions
	var start, end Position
	for y := range maze {
		for x := range maze[y] {
			if maze[y][x] == 'S' {
				start = Position{x, y}
			}
			if maze[y][x] == 'E' {
				end = Position{x, y}
			}
		}
	}

	// First pass: Find the lowest cost to reach the end
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &State{pos: start, dir: East, cost: 0})

	distTo := make(map[string]int)
	pathInfo := make(map[string]*PathInfo)
	lowestEndCost := -1

	// Run Dijkstra's algorithm
	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*State)
		currentKey := stateKey(current.pos, current.dir)

		// Skip if we've found a better path to this state
		if cost, exists := distTo[currentKey]; exists && cost < current.cost {
			continue
		}

		// Update lowest end cost if we've reached the end
		if current.pos == end {
			if lowestEndCost == -1 || current.cost < lowestEndCost {
				lowestEndCost = current.cost
			}
			continue
		}

		// Initialize or update path info
		if _, exists := pathInfo[currentKey]; !exists {
			pathInfo[currentKey] = &PathInfo{
				cost: current.cost,
				next: make(map[string]bool),
			}
		}

		// Try all possible moves
		moves := []struct {
			pos  Position
			dir  Direction
			cost int
		}{
			{getNextPosition(current.pos, current.dir), current.dir, 1},
			{current.pos, rotateClockwise(current.dir), 1000},
			{current.pos, rotateCounterClockwise(current.dir), 1000},
		}

		for _, move := range moves {
			if move.pos.y < 0 || move.pos.y >= len(maze) ||
				move.pos.x < 0 || move.pos.x >= len(maze[0]) ||
				maze[move.pos.y][move.pos.x] == '#' {
				continue
			}

			nextCost := current.cost + move.cost
			nextKey := stateKey(move.pos, move.dir)

			if nextCost <= lowestEndCost || lowestEndCost == -1 {
				if cost, exists := distTo[nextKey]; !exists || nextCost < cost {
					distTo[nextKey] = nextCost
					pathInfo[currentKey].next[nextKey] = true
					heap.Push(&pq, &State{
						pos:  move.pos,
						dir:  move.dir,
						cost: nextCost,
					})
				} else if nextCost == cost {
					// Add this as an alternative optimal path
					pathInfo[currentKey].next[nextKey] = true
				}
			}
		}
	}

	// Second pass: Find all cells that are part of paths reaching the end with optimal cost
	optimalCells := make(map[Position]bool)
	seen := make(map[string]bool)

	var dfs func(pos Position, dir Direction) bool
	dfs = func(pos Position, dir Direction) bool {
		if pos == end {
			return true
		}

		currentKey := stateKey(pos, dir)
		if seen[currentKey] {
			return false
		}
		seen[currentKey] = true
		defer func() { seen[currentKey] = false }()

		info := pathInfo[currentKey]
		if info == nil {
			return false
		}

		reachesEnd := false
		for nextKey := range info.next {
			var nextPos Position
			var nextDir Direction
			fmt.Sscanf(nextKey, "%d,%d,%d", &nextPos.x, &nextPos.y, &nextDir)

			if distTo[nextKey] <= lowestEndCost && dfs(nextPos, nextDir) {
				reachesEnd = true
			}
		}

		if reachesEnd {
			optimalCells[pos] = true
		}
		return reachesEnd
	}

	// Start DFS from the starting position
	dfs(start, East)

	// Ensure end position is marked
	optimalCells[end] = true

	return lowestEndCost, optimalCells
}

func ParseMaze(r io.Reader) ([]string, error) {
	var maze []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			maze = append(maze, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading maze: %w", err)
	}

	return maze, nil
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	maze, err := ParseMaze(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	lowestScore, optimalCells := FindLowestScoreWithPaths(maze)
	fmt.Println("Lowest score:", lowestScore)
	fmt.Println("Number of cells on optimal paths:", len(optimalCells))
}
