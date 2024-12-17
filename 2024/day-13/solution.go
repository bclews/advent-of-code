package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Pair struct{ x, y int }

type ClawMachine struct {
	ButtonA, ButtonB, Prize Pair
}

func parseFile(r io.Reader) ([]ClawMachine, error) {
	scanner := bufio.NewScanner(r)

	clawMachines := []ClawMachine{}
	clawMachine := ClawMachine{}
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			clawMachine = ClawMachine{}
		}

		var x, y int
		switch {
		case strings.HasPrefix(line, "Button A"):
			fmt.Sscanf(line, "Button A: X+%d, Y+%d", &x, &y)
			clawMachine.ButtonA = Pair{x, y}

		case strings.HasPrefix(line, "Button B"):
			fmt.Sscanf(line, "Button B: X+%d, Y+%d", &x, &y)
			clawMachine.ButtonB = Pair{x, y}

		case strings.HasPrefix(line, "Prize"):
			fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y)
			clawMachine.Prize = Pair{x, y}
			clawMachines = append(clawMachines, clawMachine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}
	return clawMachines, nil
}

// Solve the claw machine puzzle by calculating the minimum number of tokens
// required to win the prize. It uses a system of linear Diophantine equations
// to determine the number of times each button (A and B) needs to be pressed
// to align the claw exactly with the prize's position.
//
// Returns the total number of tokens required if a solution exists,
// otherwise returns 0 if no valid solution is found.
func solve(clawMachine ClawMachine) int {
	// We have the following system of equations derived from the problem:
	// Equation 1: nA * ButtonA.x + nB * ButtonB.x = Prize.x
	// Equation 2: nA * ButtonA.y + nB * ButtonB.y = Prize.y
	//
	//  	// To further expand on the above equations, we have:
	//   - Button A moves (aX, aY) units
	//   - Button B moves (bX, bY) units.
	//   - nA and nB are the number of times buttons A and B are pressed, respectively.
	//   - The claw starts at (0, 0) and needs to reach (targetX, targetY).
	//   - You must determine integers nA and nB such that:
	//         nA⋅aX+nB⋅bX=targetX (Equation 1)
	//         nA⋅aY+nB⋅bY=targetY (Equation 2)
	//   - Additionally, the cost equation is determined from the:
	//     problem statement:
	//         "It costs 3 tokens to push the A button and 1 token to push the B button."
	//     So, the total cost is calculated as:
	//         cost=3⋅nA+nB (Equation 3)

	// First, we calculate the value of 'b' (the number of times the B button is pressed),
	// by rearranging the first equation for nB:
	// nB = (Prize.x * ButtonA.y - Prize.y * ButtonA.x) / (ButtonA.x * ButtonB.y - ButtonA.y * ButtonB.x)
	// This formula is derived from solving the system of linear equations.
	// It ensures that the result will satisfy both the X and Y positions of the prize.
	b := (clawMachine.ButtonA.x*clawMachine.Prize.y - clawMachine.ButtonA.y*clawMachine.Prize.x) /
		(clawMachine.ButtonA.x*clawMachine.ButtonB.y - clawMachine.ButtonA.y*clawMachine.ButtonB.x)

	// Now we can calculate the value of 'a' (the number of times the A button is pressed),
	// using the first equation (rearranged):
	// nA = (Prize.x - b * ButtonB.x) / ButtonA.x
	// This ensures that the calculated nA and nB together will reach the prize's X position.
	a := (clawMachine.Prize.x - b*clawMachine.ButtonB.x) / clawMachine.ButtonA.x

	// Now we check if the values for nA and nB satisfy both the X and Y equations.
	// We substitute 'a' and 'b' into both equations to ensure correctness:
	// If both X and Y equations are satisfied, it means we have found a valid solution.
	if a*clawMachine.ButtonA.x+b*clawMachine.ButtonB.x == clawMachine.Prize.x &&
		a*clawMachine.ButtonA.y+b*clawMachine.ButtonB.y == clawMachine.Prize.y {
		// If a valid solution is found, calculate the total cost in tokens.
		// Button A costs 3 tokens per press, and Button B costs 1 token per press.
		// The total cost is 3 * nA + 1 * nB.
		return 3*a + b
	}

	// If no valid solution is found, return 0 to indicate that the prize cannot be won.
	return 0
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
	clawMachines, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		os.Exit(1)
	}

	// Part one
	tokens := 0
	for _, clawMachine := range clawMachines {
		tokens += solve(clawMachine)
	}
	fmt.Println("The number of tokens required to win all prizes for part one is:", tokens)

	// Part two
	// Iterate over `clawMachine` and add 10000000000000 to `Prize.x` and `Prize.y`
	tokens = 0
	for _, clawMachine := range clawMachines {
		clawMachine.Prize.x += 10000000000000
		clawMachine.Prize.y += 10000000000000
		tokens += solve(clawMachine)
	}
	fmt.Println("The number of tokens required to win all prizes for part two is:", tokens)
}
