package main

import (
	"strings"
	"testing"
)

const exampleInput = `kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`

func parseTestInput(input string) NetworkGraph {
	network := make(NetworkGraph)
	for _, line := range strings.Split(input, "\n") {
		computers := strings.Split(line, "-")
		network.AddConnection(computers[0], computers[1])
	}
	return network
}

func TestConnection(t *testing.T) {
	network := parseTestInput(exampleInput)

	// Test some explicit connections from the example
	testCases := []struct {
		comp1, comp2 string
		expected     bool
	}{
		{"kh", "tc", true},
		{"tc", "kh", true}, // Test bidirectional
		{"qp", "kh", true},
		{"de", "cg", true},
		{"kh", "cg", false}, // Not directly connected
	}

	for _, tc := range testCases {
		result := network.AreConnected(tc.comp1, tc.comp2)
		if result != tc.expected {
			t.Errorf("AreConnected(%q, %q) = %v; want %v",
				tc.comp1, tc.comp2, result, tc.expected)
		}
	}
}

func TestTriplets(t *testing.T) {
	network := parseTestInput(exampleInput)

	// Example triplets from the problem
	expectedTriplets := []string{
		"aq,cg,yn",
		"aq,vc,wq",
		"co,de,ka",
		"co,de,ta",
		"co,ka,ta",
		"de,ka,ta",
		"kh,qp,ub",
		"qp,td,wh",
		"tb,vc,wq",
		"tc,td,wh",
		"td,wh,yn",
		"ub,vc,wq",
	}

	// Find triplets
	singles := network.AllComputers()
	pairs := network.FindLargerGroups(singles)
	triplets := network.FindLargerGroups(pairs)

	// Convert expected triplets to a set for easier comparison
	expectedSet := make(ComputerSet)
	for _, triplet := range expectedTriplets {
		expectedSet[triplet] = struct{}{}
	}

	if len(triplets) != len(expectedSet) {
		t.Errorf("Found %d triplets; want %d", len(triplets), len(expectedSet))
	}

	for triplet := range triplets {
		if _, exists := expectedSet[triplet]; !exists {
			t.Errorf("Found unexpected triplet: %s", triplet)
		}
	}

	for triplet := range expectedSet {
		if _, exists := triplets[triplet]; !exists {
			t.Errorf("Missing expected triplet: %s", triplet)
		}
	}
}

func TestTTripletsCount(t *testing.T) {
	network := parseTestInput(exampleInput)

	// Expected triplets with 't' from the problem:
	expectedTTriplets := []string{
		"co,de,ta",
		"co,ka,ta",
		"de,ka,ta",
		"qp,td,wh",
		"tb,vc,wq",
		"tc,td,wh",
		"td,wh,yn",
	}

	// Find all triplets
	singles := network.AllComputers()
	pairs := network.FindLargerGroups(singles)
	triplets := network.FindLargerGroups(pairs)

	// Count triplets with 't'
	tCount := 0
	for computers := range triplets {
		computerList := strings.Split(computers, ",")
		for _, computer := range computerList {
			if strings.HasPrefix(computer, "t") {
				tCount++
				break
			}
		}
	}

	if tCount != len(expectedTTriplets) {
		t.Errorf("Found %d triplets with 't'; want %d", tCount, len(expectedTTriplets))
	}
}

func TestMaximumClique(t *testing.T) {
	network := parseTestInput(exampleInput)

	// Find the maximum clique
	currentGroups := network.AllComputers()
	var lastValidGroups ComputerSet

	for len(currentGroups) > 0 {
		lastValidGroups = currentGroups
		currentGroups = network.FindLargerGroups(currentGroups)
	}

	// Expected maximum clique from the problem
	expectedClique := "co,de,ka,ta"

	// Get the first (and should be only) maximum clique
	var foundClique string
	for clique := range lastValidGroups {
		foundClique = clique
		break
	}

	if foundClique != expectedClique {
		t.Errorf("Found maximum clique %q; want %q", foundClique, expectedClique)
	}

	// Verify all computers in the clique are connected to each other
	computersInClique := strings.Split(foundClique, ",")
	for i, comp1 := range computersInClique {
		for j, comp2 := range computersInClique {
			if i != j && !network.AreConnected(comp1, comp2) {
				t.Errorf("Computers %q and %q in clique are not connected", comp1, comp2)
			}
		}
	}
}
