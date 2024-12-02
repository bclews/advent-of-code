# Advent of Code

This repository contains my solutions to the [Advent of Code](https://adventofcode.com/) challenges. Each year, the Advent of Code provides a series of programming puzzles, one for each day in December leading up to Christmas.

## Repository Structure

The repository is organized by year, and each day's solution is stored in its own directory. For example:

```
.
├── 2024
│   └── day-01
│       ├── README.md         # Description of the problem and solution approach
│       ├── go.mod            # Go module file
│       ├── input.txt         # Puzzle input for the day
│       ├── solution          # Compiled solution binary
│       ├── solution.go       # Go source code for the solution
│       └── solution_test.go  # Unit tests for the solution
```

## How to Use

1. Clone the repository:

   ```bash
   git clone https://github.com/bclews/advent-of-code.git
   cd advent-of-code
   ```

2. Navigate to the desired year and day:

   ```bash
   cd 2024/day-01
   ```

3. Run the solution:

   ```bash
   go run solution.go
   ```

4. Run the tests:

   ```bash
   go test
   ```

## Contributing

This repository is primarily for personal learning and development. However, if you spot any errors or have suggestions, feel free to open an issue or submit a pull request.

## License

This repository is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

Happy Coding! 🎄
