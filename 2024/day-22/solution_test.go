package main

import (
	"reflect"
	"testing"
)

func TestGenerateNextSecret(t *testing.T) {
	tests := []struct {
		name           string
		initialSecret  int
		expectedSeq    []int
		sequenceLength int
	}{
		{
			name:          "Example sequence from problem",
			initialSecret: 123,
			expectedSeq: []int{
				15887950,
				16495136,
				527345,
				704524,
				1553684,
				12683156,
				11100544,
				12249484,
				7753432,
				5908254,
			},
			sequenceLength: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyer := &Buyer{secret: tt.initialSecret}
			var sequence []int

			for i := 0; i < tt.sequenceLength; i++ {
				buyer.generateNextSecret()
				sequence = append(sequence, buyer.secret)
			}

			if !reflect.DeepEqual(sequence, tt.expectedSeq) {
				t.Errorf("Expected sequence %v, got %v", tt.expectedSeq, sequence)
			}
		})
	}
}

func TestGetPrice(t *testing.T) {
	tests := []struct {
		name          string
		initialSecret int
		expectedSeq   []int // sequence of prices (ones digits)
	}{
		{
			name:          "Example price sequence from problem",
			initialSecret: 123,
			expectedSeq:   []int{3, 0, 6, 5, 4, 4, 6, 4, 4, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyer := &Buyer{secret: tt.initialSecret}
			prices := []int{buyer.getPrice()}

			for i := 0; i < len(tt.expectedSeq)-1; i++ {
				buyer.generateNextSecret()
				prices = append(prices, buyer.getPrice())
			}

			if !reflect.DeepEqual(prices, tt.expectedSeq) {
				t.Errorf("Expected price sequence %v, got %v", tt.expectedSeq, prices)
			}
		})
	}
}

func TestSimulateBuyerSequence(t *testing.T) {
	tests := []struct {
		name           string
		initialSecret  int
		expected2000th int
	}{
		{name: "Initial secret 1", initialSecret: 1, expected2000th: 8685429},
		{name: "Initial secret 10", initialSecret: 10, expected2000th: 4700978},
		{name: "Initial secret 100", initialSecret: 100, expected2000th: 15273692},
		{name: "Initial secret 2024", initialSecret: 2024, expected2000th: 8667524},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buyer := &Buyer{secret: tt.initialSecret}
			result := buyer.simulateBuyerSequence()
			if result != tt.expected2000th {
				t.Errorf("Expected 2000th secret to be %d, got %d", tt.expected2000th, result)
			}
		})
	}
}

func TestFindOptimalTradeSequence(t *testing.T) {
	tests := []struct {
		name            string
		initialSecrets  []int
		expectedBananas int
	}{
		{
			name:            "Example from problem part 2",
			initialSecrets:  []int{1, 2, 3, 2024},
			expectedBananas: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create buyers from initial secrets
			var buyers []*Buyer
			for i, secret := range tt.initialSecrets {
				buyers = append(buyers, &Buyer{id: i, secret: secret})
			}

			simulator := &MarketSimulator{buyers: buyers}
			maxBananas := simulator.findOptimalTradeSequence()

			if maxBananas != tt.expectedBananas {
				t.Errorf("Expected maximum bananas to be %d, got %d", tt.expectedBananas, maxBananas)
			}
		})
	}
}

