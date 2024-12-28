package main

import (
	"strings"
	"testing"
)

const exampleMaze = `###############
#...#...#.....#
#.#.#.#.#.###.#
#S#...#.#.#...#
#######.#.#.###
#######.#.#...#
#######.#.###.#
###..E#...#...#
###.#######.###
#...###...#...#
#.#####.#.###.#
#.#...#.#.#...#
#.#.#.#.#.#.###
#...#...#...###
###############`

func TestParseMaze(t *testing.T) {
	reader := strings.NewReader(exampleMaze)
	maze, err := ParseMaze(reader)
	if err != nil {
		t.Fatalf("ParseMaze returned unexpected error: %v", err)
	}

	// Check start position (3,1) and finish position (7,5)
	expectedStart := Position{3, 1}
	expectedFinish := Position{7, 5}

	if maze.start != expectedStart {
		t.Errorf("wrong start position: got %v, want %v", maze.start, expectedStart)
	}
	if maze.finish != expectedFinish {
		t.Errorf("wrong finish position: got %v, want %v", maze.finish, expectedFinish)
	}

	// Verify some wall positions
	wallTests := []struct {
		pos      Position
		expected bool
	}{
		{Position{0, 0}, true},   // corner wall
		{Position{1, 1}, false},  // path
		{Position{3, 1}, false},  // start position
		{Position{7, 5}, false},  // end position
		{Position{14, 14}, true}, // bottom-right corner
	}

	for _, tt := range wallTests {
		_, isWall := maze.walls[tt.pos]
		if isWall != tt.expected {
			t.Errorf("position %v: got wall=%v, want wall=%v", tt.pos, isWall, tt.expected)
		}
	}
}

func TestFindShortestPath(t *testing.T) {
	reader := strings.NewReader(exampleMaze)
	maze, err := ParseMaze(reader)
	if err != nil {
		t.Fatalf("ParseMaze returned unexpected error: %v", err)
	}

	path := maze.FindShortestPath()
	if path == nil {
		t.Fatal("FindShortestPath returned nil, expected valid path")
	}

	// The problem states that the shortest path takes 84 picoseconds
	// Each move takes 1 picosecond, so the path length should be 85 (including start position)
	expectedLength := 85
	if len(path) != expectedLength {
		t.Errorf("wrong path length: got %d, want %d", len(path), expectedLength)
	}
}

func TestCountValidCheats(t *testing.T) {
	reader := strings.NewReader(exampleMaze)
	maze, err := ParseMaze(reader)
	if err != nil {
		t.Fatalf("ParseMaze returned unexpected error: %v", err)
	}

	path := maze.FindShortestPath()
	if path == nil {
		t.Fatal("FindShortestPath returned nil, expected valid path")
	}

	// Test cases from the problem examples
	tests := []struct {
		name       string
		params     CheatParams
		wantCheats int
		savedTimes map[int]int // map of time saved to number of cheats
	}{
		{
			name: "Part 1 example cheats",
			params: CheatParams{
				maxDistance:  2,
				minTimeSaved: 0, // Set to 0 to count all cheats
			},
			savedTimes: map[int]int{
				2:  14, // 14 cheats that save 2 picoseconds
				4:  14, // 14 cheats that save 4 picoseconds
				6:  2,  // 2 cheats that save 6 picoseconds
				8:  4,  // 4 cheats that save 8 picoseconds
				10: 2,  // 2 cheats that save 10 picoseconds
				12: 3,  // 3 cheats that save 12 picoseconds
				20: 1,  // 1 cheat that saves 20 picoseconds
				36: 1,  // 1 cheat that saves 36 picoseconds
				38: 1,  // 1 cheat that saves 38 picoseconds
				40: 1,  // 1 cheat that saves 40 picoseconds
				64: 1,  // 1 cheat that saves 64 picoseconds
			},
		},
		{
			name: "Part 2 example cheats ≥ 50",
			params: CheatParams{
				maxDistance:  20,
				minTimeSaved: 50,
			},
			savedTimes: map[int]int{
				50: 32, // 32 cheats that save 50 picoseconds
				52: 31, // 31 cheats that save 52 picoseconds
				54: 29, // 29 cheats that save 54 picoseconds
				56: 39, // 39 cheats that save 56 picoseconds
				58: 25, // 25 cheats that save 58 picoseconds
				60: 23, // 23 cheats that save 60 picoseconds
				62: 20, // 20 cheats that save 62 picoseconds
				64: 19, // 19 cheats that save 64 picoseconds
				66: 12, // 12 cheats that save 66 picoseconds
				68: 14, // 14 cheats that save 68 picoseconds
				70: 12, // 12 cheats that save 70 picoseconds
				72: 22, // 22 cheats that save 72 picoseconds
				74: 4,  // 4 cheats that save 74 picoseconds
				76: 3,  // 3 cheats that save 76 picoseconds
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cheatsByTimeSaved := make(map[int]int)

			// Count cheats for each amount of time saved
			for i := 0; i < len(path)-1; i++ {
				for j := i + 1; j < len(path); j++ {
					dist := ManhattanDistance(path[i], path[j])
					if dist > 0 && dist <= tt.params.maxDistance {
						timeSaved := j - i - dist
						if timeSaved >= tt.params.minTimeSaved {
							cheatsByTimeSaved[timeSaved]++
						}
					}
				}
			}

			// Verify the number of cheats for each time saved matches the example
			for timeSaved, expectedCount := range tt.savedTimes {
				gotCount := cheatsByTimeSaved[timeSaved]
				if gotCount != expectedCount {
					t.Errorf("for time saved %d: got %d cheats, want %d cheats",
						timeSaved, gotCount, expectedCount)
				}
			}
		})
	}
}

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		p1, p2 Position
		want   int
	}{
		{Position{0, 0}, Position{3, 0}, 3},
		{Position{0, 0}, Position{0, 5}, 5},
		{Position{1, 1}, Position{4, 5}, 7},
		{Position{3, 2}, Position{0, 0}, 5},
		{Position{2, 2}, Position{2, 2}, 0},
	}

	for _, tt := range tests {
		got := ManhattanDistance(tt.p1, tt.p2)
		if got != tt.want {
			t.Errorf("ManhattanDistance(%v, %v) = %d, want %d",
				tt.p1, tt.p2, got, tt.want)
		}
	}
}
