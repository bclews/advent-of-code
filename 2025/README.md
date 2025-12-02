# Advent of Code 2025 - Rust Solutions

This directory contains my solutions for Advent of Code 2025, implemented in Rust.

## Structure

Each day is an independent Rust crate with its own `Cargo.toml`:

```
2025/
├── day-01/             # Each day follows this structure
│   ├── Cargo.toml
│   ├── src/main.rs     # Solution with inline tests
│   ├── input.txt
│   └── README.md
├── scripts/
│   └── new-day.sh      # Creates new day from template
└── templates/
    └── day-template/   # Template for new days
```

## Quick Start

### Running a Solution

```bash
cd day-01
cargo run              # Debug mode
cargo run --release    # Release mode (faster)
```

### Running Tests

```bash
cd day-01
cargo test             # Run all tests
cargo test -- --nocapture  # With output
```

### Creating a New Day

```bash
./scripts/new-day.sh 1    # Creates day-01/
```

## Development Pattern

Each solution follows a consistent pattern:

1. **Parse Input**: `parse_input()` reads and structures the data
2. **Part One**: `part_one()` solves the first puzzle
3. **Part Two**: `part_two()` solves the second puzzle
4. **Tests**: Inline tests in `src/main.rs` using example data from the puzzle
5. **Main**: Reads `input.txt` and prints both answers

Solutions are self-contained with no shared dependencies - each day handles its own file reading and logic.
