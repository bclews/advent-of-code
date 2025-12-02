# Advent of Code

This repository contains my solutions to the [Advent of Code](https://adventofcode.com/) challenges. Each year, the Advent of Code provides a series of programming puzzles, one for each day in December leading up to Christmas.

## Languages

This repository contains solutions in multiple languages:

- **2024**: Go - All 25 days completed
- **2025**: Rust - In progress

Each year maintains its own independent structure and build system.

## Repository Structure

The repository is organized by year, and each day's solution is stored in its own directory.

### 2024 (Go) Example:

```
2024/
└── day-01/
    ├── README.md         # Description of the problem and solution approach
    ├── go.mod            # Go module file
    ├── input.txt         # Puzzle input for the day
    ├── solution          # Compiled solution binary
    ├── solution.go       # Go source code for the solution
    └── solution_test.go  # Unit tests for the solution
```

### 2025 (Rust) Example:

```
2025/
└── day-01/
    ├── Cargo.toml        # Rust crate manifest
    ├── src/
    │   └── main.rs       # Rust source code for the solution
    ├── tests/            # Integration tests
    ├── input.txt         # Puzzle input for the day
    └── README.md         # Description of the problem and solution approach
```

## How to Use

1. Clone the repository:

   ```bash
   git clone https://github.com/bclews/advent-of-code.git
   cd advent-of-code
   ```

2. Navigate to the desired year and day:

   ```bash
   cd 2024/day-01  # For Go solutions
   # or
   cd 2025/day-01  # For Rust solutions
   ```

### Running 2024 (Go) Solutions

```bash
cd 2024/day-01
go run solution.go
go test
```

### Running 2025 (Rust) Solutions

```bash
cd 2025/day-01
cargo run              # Debug mode
cargo run --release    # Release mode (faster)
cargo test
```

## Contributing

This repository is primarily for personal learning and development. However, if you spot any errors or have suggestions, feel free to open an issue or submit a pull request.

## License

This repository is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

Happy Coding! 🎄
