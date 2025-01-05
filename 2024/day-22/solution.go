package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	moduloValue    = 16777216 // Magic number from problem statement
	sequenceLength = 4        // Length of price change sequence
	secretCount    = 2000     // Number of secrets to generate per buyer
)

// Represents a monkey market buyer with their secret number sequence
type Buyer struct {
	id     int
	secret int
}

// Represents a sequence of price changes that triggers a sale
type PriceChangeSequence [sequenceLength]int

// Handles the simulation of the monkey market
type MarketSimulator struct {
	buyers []*Buyer
}

// Reads the initial secret numbers from stdin
func parseFile(r io.Reader) ([]*Buyer, error) {
	var buyers []*Buyer
	scanner := bufio.NewScanner(r)
	id := 0
	for scanner.Scan() {
		secret, _ := strconv.Atoi(scanner.Text())
		buyers = append(buyers, &Buyer{id: id, secret: secret})
		id++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return buyers, nil
}

// Computes the next secret number in the sequence
func (b *Buyer) generateNextSecret() {
	secret := b.secret
	// Step 1: Multiply by 64 (left shift 6) and mix
	secret = (secret ^ (secret << 6)) % moduloValue
	// Step 2: Divide by 32 (right shift 5) and mix
	secret = (secret ^ (secret >> 5)) % moduloValue
	// Step 3: Multiply by 2048 (left shift 11) and mix
	secret = (secret ^ (secret << 11)) % moduloValue
	b.secret = secret
}

// Returns the ones digit of the current secret
func (b *Buyer) getPrice() int {
	return b.secret % 10
}

// Generates secrets and returns the 2000th one
func (b *Buyer) simulateBuyerSequence() int {
	for i := 0; i < secretCount; i++ {
		b.generateNextSecret()
	}
	return b.secret
}

// Finds the sequence of price changes that yields maximum bananas
func (ms *MarketSimulator) findOptimalTradeSequence() int {
	// Maps sequence of price changes to the first price seen after that sequence for each buyer
	priceAfterSequence := make(map[int]map[PriceChangeSequence]int)

	// Maps sequence to the set of buyers that exhibit that sequence
	sequenceBuyers := make(map[PriceChangeSequence]map[int]struct{})

	// Simulate each buyer's price changes
	for _, buyer := range ms.buyers {
		priceAfterSequence[buyer.id] = make(map[PriceChangeSequence]int)

		var priceChanges PriceChangeSequence
		ringIndex := 0
		prevPrice := buyer.getPrice()

		for t := 1; t <= secretCount; t++ {
			buyer.generateNextSecret()
			newPrice := buyer.getPrice()
			delta := newPrice - prevPrice
			prevPrice = newPrice

			priceChanges[ringIndex] = delta
			ringIndex = (ringIndex + 1) % sequenceLength

			if t >= sequenceLength {
				sequence := ms.getCurrentSequence(priceChanges, ringIndex)

				// Record the first occurrence of this sequence for this buyer
				if _, exists := priceAfterSequence[buyer.id][sequence]; !exists {
					priceAfterSequence[buyer.id][sequence] = newPrice

					if sequenceBuyers[sequence] == nil {
						sequenceBuyers[sequence] = make(map[int]struct{})
					}
					sequenceBuyers[sequence][buyer.id] = struct{}{}
				}
			}
		}
	}

	return ms.calculateMaxBananas(priceAfterSequence, sequenceBuyers)
}

// Returns the current sequence of price changes in correct order
func (ms *MarketSimulator) getCurrentSequence(ring PriceChangeSequence, currentIndex int) PriceChangeSequence {
	var sequence PriceChangeSequence
	for i := 0; i < sequenceLength; i++ {
		sequence[i] = ring[(currentIndex+i)%sequenceLength]
	}
	return sequence
}

// Finds the sequence that yields the maximum total bananas
func (ms *MarketSimulator) calculateMaxBananas(
	priceAfterSequence map[int]map[PriceChangeSequence]int,
	sequenceBuyers map[PriceChangeSequence]map[int]struct{},
) int {
	maxBananas := 0
	for sequence, buyers := range sequenceBuyers {
		totalBananas := 0
		for buyerID := range buyers {
			totalBananas += priceAfterSequence[buyerID][sequence]
		}
		if totalBananas > maxBananas {
			maxBananas = totalBananas
		}
	}
	return maxBananas
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	buyers, err := parseFile(file)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// Part 1: Sum of 2000th secret for each buyer
	totalSecret := 0
	for _, buyer := range buyers {
		buyerCopy := *buyer // Create copy to not affect Part 2
		totalSecret += buyerCopy.simulateBuyerSequence()
	}
	fmt.Println("Part 1:", totalSecret)

	// Part 2: Find optimal trading sequence
	simulator := &MarketSimulator{buyers: buyers}
	maxBananas := simulator.findOptimalTradeSequence()
	fmt.Println("Part 2:", maxBananas)
}
