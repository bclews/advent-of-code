package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Coord represents a 2D coordinate with x and y values
type Coord struct {
	x, y int
}

// Grid represents the 2D grid of input data
type Grid map[Coord]uint8

// parseFile reads the input file and creates a grid representation
func parseFile(r io.Reader) (Grid, int, error) {
	scanner := bufio.NewScanner(r)
	grid := make(Grid)
	height := 0

	for scanner.Scan() {
		line := scanner.Text()
		for x, char := range line {
			grid[Coord{height, x}] = uint8(char)
		}
		height++
	}

	if err := scanner.Err(); err != nil {
		return nil, -1, fmt.Errorf("error scanning file: %w", err)
	}

	return grid, height, nil
}

// exploreRegion explores a connected region in the grid
type regionExplorer struct {
	grid     Grid
	visited  map[Coord]bool
	regionID uint8
}

// explore calculates area, perimeter, and corners of a region
// Got the idea of exploring corners from:
// https://www.reddit.com/r/adventofcode/comments/1hcf16m/2024_day_12_everyone_must_be_hating_today_so_here
func (re *regionExplorer) explore(start Coord) (area, perimeter, corners int) {
	if re.visited[start] {
		return 0, 0, 0
	}

	re.regionID = re.grid[start]
	re.visited[start] = true

	var dfs func(Coord)
	dfs = func(current Coord) {
		area++

		// Cardinal directions
		cardinalOffsets := []Coord{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, offset := range cardinalOffsets {
			next := Coord{current.x + offset.x, current.y + offset.y}
			val, exists := re.grid[next]
			if exists && val == re.regionID && !re.visited[next] {
				re.visited[next] = true
				dfs(next)
			} else if !exists || val != re.regionID {
				perimeter++
			}
		}

		// Corner exploration
		cornerOffsets := []Coord{{-1, -1}, {1, 1}, {1, -1}, {-1, 1}}
		for _, corner := range cornerOffsets {
			// Convex corner check
			if re.grid[Coord{current.x + corner.x, current.y}] != re.regionID &&
				re.grid[Coord{current.x, current.y + corner.y}] != re.regionID {
				corners++
			}

			// Concave corner check
			if re.grid[Coord{current.x + corner.x, current.y}] == re.regionID &&
				re.grid[Coord{current.x, current.y + corner.y}] == re.regionID &&
				re.grid[Coord{current.x + corner.x, current.y + corner.y}] != re.regionID {
				corners++
			}
		}
	}

	dfs(start)
	return area, perimeter, corners
}

func calculatePrice(grid Grid, height int) (price, discountedPrice int) {
	explorer := &regionExplorer{
		grid:    grid,
		visited: make(map[Coord]bool),
	}

	for i := 0; i < height; i++ {
		for j := 0; j < height; j++ {
			current := Coord{i, j}
			if explorer.visited[current] {
				continue
			}

			area, perimeter, corners := explorer.explore(current)
			price += area * perimeter
			discountedPrice += area * corners
		}
	}
	return price, discountedPrice
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Parse the file
	grid, height, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	price, discountedPrice := calculatePrice(grid, height)
	fmt.Println("Price:", price)
	fmt.Println("Discounted Price:", discountedPrice)
}
