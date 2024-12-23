package main

import "testing"

func TestLinenLayout(t *testing.T) {
	patterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}
	designs := []string{
		"brwrr",
		"bggr",
		"gbbr",
		"rrbgbr",
		"ubwu",
		"bwurrg",
		"brgr",
		"bbrgwb",
	}
	expected := []bool{
		true,  // brwrr
		true,  // bggr
		true,  // gbbr
		true,  // rrbgbr
		false, // ubwu
		true,  // bwurrg
		true,  // brgr
		false, // bbrgwb
	}

	for i, design := range designs {
		result := isDesignPossible(patterns, design)
		if result != expected[i] {
			t.Errorf("For design %s, expected %v but got %v", design, expected[i], result)
		}
	}
}

func TestCountPossibleDesigns(t *testing.T) {
	patterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}
	designs := []string{
		"brwrr",
		"bggr",
		"gbbr",
		"rrbgbr",
		"ubwu",
		"bwurrg",
		"brgr",
		"bbrgwb",
	}
	expectedCount := 6

	if count := countPossibleDesigns(patterns, designs); count != expectedCount {
		t.Errorf("Expected %d possible designs, but got %d", expectedCount, count)
	}
}

func TestCountWaysToFormDesign(t *testing.T) {
	patterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}
	designs := map[string]int{
		"brwrr":  2,
		"bggr":   1,
		"gbbr":   4,
		"rrbgbr": 6,
		"bwurrg": 1,
		"brgr":   2,
		"ubwu":   0,
		"bbrgwb": 0,
	}

	for design, expectedWays := range designs {
		result := countWaysToFormDesign(patterns, design)
		if result != expectedWays {
			t.Errorf("For design %s, expected %d ways but got %d", design, expectedWays, result)
		}
	}
}

func TestCountTotalWays(t *testing.T) {
	patterns := []string{"r", "wr", "b", "g", "bwu", "rb", "gb", "br"}
	designs := []string{
		"brwrr",
		"bggr",
		"gbbr",
		"rrbgbr",
		"ubwu",
		"bwurrg",
		"brgr",
		"bbrgwb",
	}
	expectedTotalWays := 16

	if totalWays := countTotalWays(patterns, designs); totalWays != expectedTotalWays {
		t.Errorf("Expected total ways %d but got %d", expectedTotalWays, totalWays)
	}
}
