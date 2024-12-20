package main

import (
	"strings"
	"testing"
)

func TestFindLowestScore(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedLowestScore   int
		expectedBestPathTiles int
	}{
		{
			name: "Example 1",
			input: `###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`,
			expectedLowestScore:   7036,
			expectedBestPathTiles: 45,
		},
		{
			name: "Example 2",
			input: `#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`,
			expectedLowestScore:   11048,
			expectedBestPathTiles: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			maze, err := ParseMaze(reader)
			if err != nil {
				t.Fatalf("ParseMaze() error = %v", err)
			}

			lowestScore, bestPathTiles := FindLowestScoreWithPaths(maze)
			if lowestScore != tt.expectedLowestScore {
				t.Errorf("FindLowestScore() = %v, want %v", lowestScore, tt.expectedLowestScore)
			}

			if len(bestPathTiles) != tt.expectedBestPathTiles {
				t.Errorf("TraceBestPaths() = %v, want %v", len(bestPathTiles), tt.expectedBestPathTiles)
			}
		})
	}
}
