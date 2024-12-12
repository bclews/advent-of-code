package main

import (
	"testing"
)

func TestExplore(t *testing.T) {
	tests := []struct {
		name              string
		grid              Grid
		start             Coord
		expectedArea      int
		expectedPerimeter int
		expectedCorners   int
	}{
		{
			/*
			   +-+-+-+-+
			   |A A A A|
			   +-+-+-+-+
			*/
			name: "horizontal line region",
			grid: Grid{
				{0, 0}: 'A', {0, 1}: 'A', {0, 2}: 'A', {0, 3}: 'A',
			},
			start:             Coord{0, 0},
			expectedArea:      4,
			expectedPerimeter: 10,
			expectedCorners:   4,
		},
		{
			/*
				+-+
				|D|
				+-+
			*/
			name: "single cell D region",
			grid: Grid{
				{0, 0}: 'D',
			},
			start:             Coord{0, 0},
			expectedArea:      1,
			expectedPerimeter: 4,
			expectedCorners:   4,
		},
		{
			/*
				+-+-+
				|B B|
				+   +
				|B B|
				+-+-+
			*/
			name: "2x2 B region",
			grid: Grid{
				{0, 0}: 'B', {0, 1}: 'B',
				{1, 0}: 'B', {1, 1}: 'B',
			},
			start:             Coord{0, 0},
			expectedArea:      4,
			expectedPerimeter: 8,
			expectedCorners:   4,
		},
		{
			/*
			 +-+
			 |C|
			 + +-+
			 |C C|
			 +-+ +
			   |C|
			   +-+
			*/
			name: "C-shaped region",
			grid: Grid{
				{0, 0}: 'C',
				{1, 0}: 'C', {1, 1}: 'C',
				{2, 1}: 'C',
			},
			start:             Coord{0, 0},
			expectedArea:      4,
			expectedPerimeter: 10,
			expectedCorners:   8,
		},
		{
			/*
				+-+-+-+
				|E E E|
				+-+-+-+
			*/

			name: "horizontal line E region",
			grid: Grid{
				{0, 0}: 'E', {0, 1}: 'E', {0, 2}: 'E',
			},
			start:             Coord{0, 0},
			expectedArea:      3,
			expectedPerimeter: 8,
			expectedCorners:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := &regionExplorer{
				grid:    tt.grid,
				visited: make(map[Coord]bool),
			}
			area, perimeter, corners := re.explore(tt.start)
			if area != tt.expectedArea {
				t.Errorf("expected area %d, got %d", tt.expectedArea, area)
			}
			if perimeter != tt.expectedPerimeter {
				t.Errorf("expected perimeter %d, got %d", tt.expectedPerimeter, perimeter)
			}
			if corners != tt.expectedCorners {
				t.Errorf("expected corners %d, got %d", tt.expectedCorners, corners)
			}
		})
	}
}

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		grid             Grid
		height           int
		expectedPrice    int
		expectedDiscount int
	}{
		{
			/*
				AAAA
				BBCD
				BBCC
				EEEC
			*/
			grid: Grid{
				{0, 0}: 1, {0, 1}: 1, {0, 2}: 1, {0, 3}: 1,
				{1, 0}: 2, {1, 1}: 2, {1, 2}: 3, {1, 3}: 4,
				{2, 0}: 2, {2, 1}: 2, {2, 2}: 3, {2, 3}: 3,
				{3, 0}: 5, {3, 1}: 5, {3, 2}: 5, {3, 3}: 3,
			},
			height:           4,
			expectedPrice:    140,
			expectedDiscount: 80,
		},
		{
			/*
				OOOOO
				OXOXO
				OOOOO
				OXOXO
				OOOOO
			*/
			grid: Grid{
				{0, 0}: 1, {0, 1}: 1, {0, 2}: 1, {0, 3}: 1, {0, 4}: 1,
				{1, 0}: 1, {1, 1}: 0, {1, 2}: 1, {1, 3}: 0, {1, 4}: 1,
				{2, 0}: 1, {2, 1}: 1, {2, 2}: 1, {2, 3}: 1, {2, 4}: 1,
				{3, 0}: 1, {3, 1}: 0, {3, 2}: 1, {3, 3}: 0, {3, 4}: 1,
				{4, 0}: 1, {4, 1}: 1, {4, 2}: 1, {4, 3}: 1, {4, 4}: 1,
			},
			height:           5,
			expectedPrice:    772,
			expectedDiscount: 436,
		},
		{
			/*
				RRRRIICCFF
				RRRRIICCCF
				VVRRRCCFFF
				VVRCCCJFFF
				VVVVCJJCFE
				VVIVCCJJEE
				VVIIICJJEE
				MIIIIIJJEE
				MIIISIJEEE
				MMMISSJEEE
			*/
			grid: Grid{
				{0, 0}: 82, {0, 1}: 82, {0, 2}: 82, {0, 3}: 82, {0, 4}: 73,
				{0, 5}: 73, {0, 6}: 67, {0, 7}: 67, {0, 8}: 70, {0, 9}: 70,
				{1, 0}: 82, {1, 1}: 82, {1, 2}: 82, {1, 3}: 82, {1, 4}: 73,
				{1, 5}: 73, {1, 6}: 67, {1, 7}: 67, {1, 8}: 67, {1, 9}: 70,
				{2, 0}: 86, {2, 1}: 86, {2, 2}: 82, {2, 3}: 82, {2, 4}: 82,
				{2, 5}: 67, {2, 6}: 67, {2, 7}: 70, {2, 8}: 70, {2, 9}: 70,
				{3, 0}: 86, {3, 1}: 86, {3, 2}: 82, {3, 3}: 67, {3, 4}: 67,
				{3, 5}: 67, {3, 6}: 74, {3, 7}: 70, {3, 8}: 70, {3, 9}: 70,
				{4, 0}: 86, {4, 1}: 86, {4, 2}: 86, {4, 3}: 86, {4, 4}: 67,
				{4, 5}: 74, {4, 6}: 74, {4, 7}: 67, {4, 8}: 70, {4, 9}: 69,
				{5, 0}: 86, {5, 1}: 86, {5, 2}: 73, {5, 3}: 86, {5, 4}: 67,
				{5, 5}: 67, {5, 6}: 74, {5, 7}: 74, {5, 8}: 69, {5, 9}: 69,
				{6, 0}: 86, {6, 1}: 86, {6, 2}: 73, {6, 3}: 73, {6, 4}: 73,
				{6, 5}: 67, {6, 6}: 74, {6, 7}: 74, {6, 8}: 69, {6, 9}: 69,
				{7, 0}: 77, {7, 1}: 73, {7, 2}: 73, {7, 3}: 73, {7, 4}: 73,
				{7, 5}: 73, {7, 6}: 74, {7, 7}: 74, {7, 8}: 69, {7, 9}: 69,
				{8, 0}: 77, {8, 1}: 73, {8, 2}: 73, {8, 3}: 73, {8, 4}: 83,
				{8, 5}: 73, {8, 6}: 74, {8, 7}: 69, {8, 8}: 69, {8, 9}: 69,
				{9, 0}: 77, {9, 1}: 77, {9, 2}: 77, {9, 3}: 73, {9, 4}: 83,
				{9, 5}: 83, {9, 6}: 74, {9, 7}: 69, {9, 8}: 69, {9, 9}: 69,
			},
			height:           10,
			expectedPrice:    1930,
			expectedDiscount: 1206,
		},
	}

	for _, tt := range tests {
		price, discountedPrice := calculatePrice(tt.grid, tt.height)
		if price != tt.expectedPrice {
			t.Errorf("expected price %d, got %d", tt.expectedPrice, price)
		}
		if discountedPrice != tt.expectedDiscount {
			t.Errorf("expected discounted price %d, got %d", tt.expectedDiscount, discountedPrice)
		}
	}
}
