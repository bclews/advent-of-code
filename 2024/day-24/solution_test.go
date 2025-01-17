package main

import (
	"strings"
	"testing"
)

func TestSmallExample(t *testing.T) {
	input := `x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02`

	circuit := NewCircuit(strings.Split(input, "\n"))
	circuit.Simulate()
	result := circuit.GetDecimalOutput("z")

	expected := int64(4) // Binary 100
	if result != expected {
		t.Errorf("Small example failed: got %d, want %d", result, expected)
	}
}

func TestLargerExample(t *testing.T) {
	input := `x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj`

	circuit := NewCircuit(strings.Split(input, "\n"))
	circuit.Simulate()
	result := circuit.GetDecimalOutput("z")

	expected := int64(2024) // Binary 0011111101000
	if result != expected {
		t.Errorf("Larger example failed: got %d, want %d", result, expected)
	}
}

// Test helper functions
func TestIsInputWire(t *testing.T) {
	tests := []struct {
		wire     string
		expected bool
	}{
		{"x00", true},
		{"y01", true},
		{"z00", false},
		{"abc", false},
	}

	for _, test := range tests {
		result := isInputWire(test.wire)
		if result != test.expected {
			t.Errorf("isInputWire(%s) = %v, want %v", test.wire, result, test.expected)
		}
	}
}

func TestIsInputPair(t *testing.T) {
	tests := []struct {
		wire1    string
		wire2    string
		expected bool
	}{
		{"x00", "y00", true},
		{"y01", "x01", true},
		{"x00", "x01", false},
		{"y00", "y01", false},
		{"z00", "x00", false},
	}

	for _, test := range tests {
		result := isInputPair(test.wire1, test.wire2)
		if result != test.expected {
			t.Errorf("isInputPair(%s, %s) = %v, want %v",
				test.wire1, test.wire2, result, test.expected)
		}
	}
}

func TestParseGate(t *testing.T) {
	tests := []struct {
		input string
		want  Gate
	}{
		{
			"x00 AND y00 -> z00",
			Gate{input1: "x00", operation: "AND", input2: "y00", output: "z00"},
		},
		{
			"x01 XOR y01 -> z01",
			Gate{input1: "x01", operation: "XOR", input2: "y01", output: "z01"},
		},
		{
			"x02 OR y02 -> z02",
			Gate{input1: "x02", operation: "OR", input2: "y02", output: "z02"},
		},
	}

	for _, test := range tests {
		got := parseGate(test.input)
		if got != test.want {
			t.Errorf("parseGate(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

