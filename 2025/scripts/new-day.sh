#!/bin/bash
# Script to create a new day's solution structure
# Usage: ./scripts/new-day.sh <day-number>
# Example: ./scripts/new-day.sh 1

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <day-number>"
  exit 1
fi

DAY=$(printf "%02d" "$1")
DAY_DIR="day-$DAY"

if [ -d "$DAY_DIR" ]; then
  echo "Error: Directory $DAY_DIR already exists"
  exit 1
fi

echo "Creating directory structure for Day $DAY..."
mkdir -p "$DAY_DIR"
cp -r templates/day-template/* "$DAY_DIR/"

# Update README with day number
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s/Day NN/Day $1/g" "$DAY_DIR/README.md"
  sed -i '' "s/day\/NN/day\/$1/g" "$DAY_DIR/README.md"
else
  # Linux
  sed -i "s/Day NN/Day $1/g" "$DAY_DIR/README.md"
  sed -i "s/day\/NN/day\/$1/g" "$DAY_DIR/README.md"
fi

echo "Day $DAY created successfully at $DAY_DIR"
echo ""
echo "Next steps:"
echo "  1. cd $DAY_DIR"
echo "  2. Add your puzzle input to input.txt"
echo "  3. Update README.md with problem description"
echo "  4. Implement solution in src/main.rs"
echo "  5. Run with: cargo run"
echo "  6. Test with: cargo test"
