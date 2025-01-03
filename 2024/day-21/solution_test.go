package main

import (
	"reflect"
	"testing"
)

func TestNewNumericKeypad(t *testing.T) {
	keypad := NewNumericKeypad()

	expectedPositions := map[byte]Position{
		'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
		'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
		'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
		'_': {3, 0}, '0': {3, 1}, 'A': {3, 2},
	}

	if !reflect.DeepEqual(keypad.buttonPositions, expectedPositions) {
		t.Errorf("NewNumericKeypad() button positions = %v, want %v",
			keypad.buttonPositions, expectedPositions)
	}

	if keypad.emptyPosition != (Position{3, 0}) {
		t.Errorf("NewNumericKeypad() empty position = %v, want %v",
			keypad.emptyPosition, Position{3, 0})
	}

	if keypad.startPosition != (Position{3, 2}) {
		t.Errorf("NewNumericKeypad() start position = %v, want %v",
			keypad.startPosition, Position{3, 2})
	}
}

func TestNewDirectionalKeypad(t *testing.T) {
	keypad := NewDirectionalKeypad()

	expectedPositions := map[byte]Position{
		'_': {0, 0}, '^': {0, 1}, 'A': {0, 2},
		'<': {1, 0}, 'v': {1, 1}, '>': {1, 2},
	}

	if !reflect.DeepEqual(keypad.buttonPositions, expectedPositions) {
		t.Errorf("NewDirectionalKeypad() button positions = %v, want %v",
			keypad.buttonPositions, expectedPositions)
	}

	if keypad.emptyPosition != (Position{0, 0}) {
		t.Errorf("NewDirectionalKeypad() empty position = %v, want %v",
			keypad.emptyPosition, Position{0, 0})
	}

	if keypad.startPosition != (Position{0, 2}) {
		t.Errorf("NewDirectionalKeypad() start position = %v, want %v",
			keypad.startPosition, Position{0, 2})
	}
}

func TestGenerateButtonSequence(t *testing.T) {
	robotChain := NewRobotChain()
	tests := []struct {
		name       string
		code       string
		robotDepth int
		wantLength int
		wantError  bool
	}{
		{
			name:       "example code 029A with depth 2",
			code:       "029A",
			robotDepth: 2,
			wantLength: 68,
			wantError:  false,
		},
		{
			name:       "example code 980A with depth 2",
			code:       "980A",
			robotDepth: 2,
			wantLength: 60,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLength, err := robotChain.GenerateButtonSequence(tt.code, tt.robotDepth)

			if (err != nil) != tt.wantError {
				t.Errorf("GenerateButtonSequence() error = %v, wantError %v",
					err, tt.wantError)
				return
			}

			if gotLength != tt.wantLength {
				t.Errorf("GenerateButtonSequence() = %v, want %v",
					gotLength, tt.wantLength)
			}
		})
	}
}

func TestParseNumericPart(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    int
		wantErr bool
	}{
		{
			name:    "valid code 029A",
			code:    "029A",
			want:    29,
			wantErr: false,
		},
		{
			name:    "valid code 980A",
			code:    "980A",
			want:    980,
			wantErr: false,
		},
		{
			name:    "invalid code ABC",
			code:    "ABC",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNumericPart(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNumericPart() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseNumericPart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompleteExample(t *testing.T) {
	type testCase struct {
		code        string
		numericPart int
		seqLength   int
	}

	examples := []testCase{
		{
			code:        "029A",
			numericPart: 29,
			seqLength:   68,
		},
		{
			code:        "980A",
			numericPart: 980,
			seqLength:   60,
		},
		{
			code:        "179A",
			numericPart: 179,
			seqLength:   68,
		},
		{
			code:        "456A",
			numericPart: 456,
			seqLength:   64,
		},
		{
			code:        "379A",
			numericPart: 379,
			seqLength:   64,
		},
	}

	expectedPartOneSum := 126384
	robotChain := NewRobotChain()
	actualSum := 0

	for _, ex := range examples {
		t.Run("Testing "+ex.code, func(t *testing.T) {
			seqLen, err := robotChain.GenerateButtonSequence(ex.code, 2)
			if err != nil {
				t.Fatalf("GenerateButtonSequence failed: %v", err)
			}

			if seqLen != ex.seqLength {
				t.Errorf("GenerateButtonSequence(%s, 2) length = %d, want %d",
					ex.code, seqLen, ex.seqLength)
			}

			complexity := seqLen * ex.numericPart
			actualSum += complexity
		})
	}

	if actualSum != expectedPartOneSum {
		t.Errorf("Part One sum = %d, want %d", actualSum, expectedPartOneSum)
	}
}
