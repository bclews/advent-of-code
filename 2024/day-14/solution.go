package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Robot struct {
	position Point
	velocity Point
}

func parseFile(r io.Reader) ([]Robot, error) {
	var robots []Robot

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		robot, err := parseInput(scanner.Text())
		if err != nil {
			return nil, err
		}
		robots = append(robots, robot)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return robots, nil
}

// parseInput parses a line like "p=0,4 v=3,-3" into a Robot
func parseInput(line string) (Robot, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return Robot{}, fmt.Errorf("invalid input format: %s", line)
	}

	position, err := parseCoordinates(strings.TrimPrefix(parts[0], "p="))
	if err != nil {
		return Robot{}, fmt.Errorf("invalid position format: %s", err)
	}

	velocity, err := parseCoordinates(strings.TrimPrefix(parts[1], "v="))
	if err != nil {
		return Robot{}, fmt.Errorf("invalid velocity format: %s", err)
	}

	return Robot{
		position: position,
		velocity: velocity,
	}, nil
}

// parseCoordinates parses a coordinate string like "0,4" into a Point
func parseCoordinates(coord string) (Point, error) {
	coords := strings.Split(coord, ",")
	if len(coords) != 2 {
		return Point{}, fmt.Errorf("invalid coordinate format: %s", coord)
	}

	x, err := strconv.Atoi(coords[0])
	if err != nil {
		return Point{}, err
	}

	y, err := strconv.Atoi(coords[1])
	if err != nil {
		return Point{}, err
	}

	return Point{x: x, y: y}, nil
}

// calculatePosition determines the final position after n seconds
func (r Robot) calculatePosition(n, width, height int) Point {
	// Calculate new position using modular arithmetic for wrap-around
	newX := ((r.position.x+r.velocity.x*n)%width + width) % width
	newY := ((r.position.y+r.velocity.y*n)%height + height) % height
	return Point{x: newX, y: newY}
}

func getQuadrant(p Point, width, height int) int {
	midX := width / 2
	midY := height / 2

	switch {
	case p.x == midX || p.y == midY:
		return 0 // On the Midline
	case p.x < midX && p.y < midY:
		return 1 // Top-left
	case p.x < midX && p.y >= midY:
		return 3 // Bottom-left
	case p.x >= midX && p.y < midY:
		return 2 // Top-right
	default:
		return 4 // Bottom-right
	}
}

func calculateSafetyFactor(robots []Robot, width, height, seconds int) int {
	// Count robots in each quadrant
	quadrantCounts := make(map[int]int)

	// Initialize all quadrants to 0
	for i := 1; i <= 4; i++ {
		quadrantCounts[i] = 0
	}

	for _, robot := range robots {
		finalPos := robot.calculatePosition(seconds, width, height)
		quadrant := getQuadrant(finalPos, width, height)
		if quadrant != 0 { // Don't count robots on midlines
			quadrantCounts[quadrant]++
		}
	}

	// If any quadrant is empty, return 0
	for i := 1; i <= 4; i++ {
		if quadrantCounts[i] == 0 {
			return 0
		}
	}

	return quadrantCounts[1] * quadrantCounts[2] * quadrantCounts[3] * quadrantCounts[4]
}

func findMaxConsecutiveRobots(robots []Robot, width, height, seconds int) (int, int) {
	maxConsecutive, maxTime := 0, 0

	for time := 0; time < seconds; time++ {
		positions := trackRobotPositions(robots, time, width, height)

		for y := 0; y < height; y++ {
			for startX := 0; startX < width; startX++ {
				consecutive := countConsecutiveRobots(positions, startX, y, width)
				if consecutive > maxConsecutive {
					maxConsecutive = consecutive
					maxTime = time
				}
			}
		}
	}

	return maxConsecutive, maxTime
}

func trackRobotPositions(robots []Robot, time, width, height int) map[Point]bool {
	positions := make(map[Point]bool)
	for _, robot := range robots {
		pos := robot.calculatePosition(time, width, height)
		positions[pos] = true
	}
	return positions
}

func countConsecutiveRobots(positions map[Point]bool, startX, y, width int) int {
	consecutive := 0
	for x := startX; x < width; x++ {
		if !positions[Point{x: x, y: y}] {
			break
		}
		consecutive++
	}
	return consecutive
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	robots, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	width := 101
	height := 103
	seconds := 100

	// Calculate the safety factor
	safetyFactor := calculateSafetyFactor(robots, width, height, seconds)
	fmt.Println("Part one:")
	fmt.Println("\tSafety factor:", safetyFactor)

	// Find the maximum consecutive robots on a line
	seconds = 10000
	maxConsecutive, maxTime := findMaxConsecutiveRobots(robots, width, height, seconds)
	fmt.Println("Part two:")
	fmt.Println("\tMaximum consecutive robots: ", maxConsecutive)
	fmt.Println("\tFound at time: ", maxTime)
}
