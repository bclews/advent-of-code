package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ComputerSet represents a set of connected computers
type ComputerSet map[string]struct{}

// NewComputerSet creates a new empty set of computers
func NewComputerSet() ComputerSet {
	return make(ComputerSet)
}

// ToSortedString converts the set to a sorted, comma-separated string
func (s ComputerSet) ToSortedString(additional ...string) string {
	computers := make([]string, 0, len(s)+len(additional))
	for computer := range s {
		computers = append(computers, computer)
	}
	computers = append(computers, additional...)
	sort.Strings(computers)
	return strings.Join(computers, ",")
}

// FromString populates the set from a comma-separated string
func (s *ComputerSet) FromString(computers string) {
	newSet := NewComputerSet()
	for _, computer := range strings.Split(computers, ",") {
		newSet[computer] = struct{}{}
	}
	*s = newSet
}

// NetworkGraph represents connections between computers
type NetworkGraph map[string]ComputerSet

// AddConnection adds a bidirectional connection between two computers
func (g NetworkGraph) AddConnection(comp1, comp2 string) {
	if _, exists := g[comp1]; !exists {
		g[comp1] = NewComputerSet()
	}
	if _, exists := g[comp2]; !exists {
		g[comp2] = NewComputerSet()
	}
	g[comp1][comp2] = struct{}{}
	g[comp2][comp1] = struct{}{}
}

// AreConnected checks if two computers are directly connected
func (g NetworkGraph) AreConnected(comp1, comp2 string) bool {
	connections, exists := g[comp1]
	if !exists {
		return false
	}
	_, connected := connections[comp2]
	return connected
}

// AllComputers returns a set of all computers in the network
func (g NetworkGraph) AllComputers() ComputerSet {
	computers := NewComputerSet()
	for computer := range g {
		computers[computer] = struct{}{}
	}
	return computers
}

// FindLargerGroups finds all possible groups with one more computer
func (g NetworkGraph) FindLargerGroups(groups ComputerSet) ComputerSet {
	newGroups := NewComputerSet()

	for groupStr := range groups {
		currentGroup := NewComputerSet()
		currentGroup.FromString(groupStr)

		// Try adding each connected computer to the group
		baseComputer := strings.Split(groupStr, ",")[0]
		for candidate := range g[baseComputer] {
			if _, alreadyInGroup := currentGroup[candidate]; alreadyInGroup {
				continue
			}

			if g.isConnectedToAllComputers(candidate, currentGroup) {
				newGroups[currentGroup.ToSortedString(candidate)] = struct{}{}
			}
		}
	}
	return newGroups
}

// isConnectedToAllComputers checks if a computer is connected to all computers in a group
func (g NetworkGraph) isConnectedToAllComputers(computer string, group ComputerSet) bool {
	for existingComputer := range group {
		if !g.AreConnected(existingComputer, computer) {
			return false
		}
	}
	return true
}

// parseFile reads the network configuration from stdin
func parseFile(r io.Reader) (NetworkGraph, error) {
	network := make(NetworkGraph)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		computers := strings.Split(scanner.Text(), "-")
		network.AddConnection(computers[0], computers[1])
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return network, nil
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	network, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Part 1: Find triplets with at least one 't' computer
	groups := network.AllComputers()            // Start with single computers
	pairs := network.FindLargerGroups(groups)   // Find pairs
	triplets := network.FindLargerGroups(pairs) // Find triplets

	tComputerCount := 0
	for computers := range triplets {
		computerList := strings.Split(computers, ",")
		// Check if any computer in the triplet starts with 't'
		for _, computer := range computerList {
			if strings.HasPrefix(computer, "t") {
				tComputerCount++
				break
			}
		}
	}
	fmt.Println("Triplets with 't' computers:", tComputerCount)

	// Part 2: Find the largest fully connected group (maximum clique)
	currentGroups := network.AllComputers()
	var lastValidGroups ComputerSet

	for len(currentGroups) > 0 {
		lastValidGroups = currentGroups
		currentGroups = network.FindLargerGroups(currentGroups)
	}

	// Print the first (and only) maximum clique found
	for password := range lastValidGroups {
		fmt.Println("LAN Party password:", password)
		return
	}
}
