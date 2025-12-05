use std::fs;

/// The 8 possible directions to check for adjacent paper rolls.
/// Ordered as: top-left, top, top-right, left, right, bottom-left, bottom, bottom-right.
const DIRECTIONS: [(i32, i32); 8] = [
    (-1, -1), (-1, 0), (-1, 1),
    (0, -1),           (0, 1),
    (1, -1),  (1, 0),  (1, 1),
];

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

/// Counts the number of paper rolls (@) adjacent to the given position.
/// Uses early termination - stops counting once 4 neighbors are found.
/// Returns a count in range [0, 4] since we only care if count < 4.
fn count_neighbors(grid: &[String], row: usize, col: usize) -> usize {
    DIRECTIONS
        .iter()
        .filter_map(|&(dr, dc)| {
            let nr = (row as i32 + dr) as usize;
            let nc = (col as i32 + dc) as usize;

            // Safe bounds checking with .get() returning Option
            grid.get(nr)?
                .chars()
                .nth(nc)
                .filter(|&ch| ch == '@')
        })
        .take(4)  // Early termination: only need to know if < 4
        .count()
}

/// Counts the number of paper rolls (@) adjacent to the given position in a 2D grid.
/// Uses early termination - stops counting once 4 neighbors are found.
/// Returns a count in range [0, 4] since we only care if count < 4.
fn count_neighbors_2d(grid: &[Vec<char>], row: usize, col: usize) -> usize {
    DIRECTIONS
        .iter()
        .filter_map(|&(dr, dc)| {
            let nr = (row as i32 + dr) as usize;
            let nc = (col as i32 + dc) as usize;

            // Safe bounds checking with .get() returning Option
            let ch = grid.get(nr)?.get(nc)?;
            (*ch == '@').then_some(ch)
        })
        .take(4)  // Early termination: only need to know if < 4
        .count()
}

/// Counts the number of paper rolls that forklifts can access.
/// A roll is accessible if it has fewer than 4 adjacent rolls.
fn part_one(data: &[String]) -> i64 {
    data.iter()
        .enumerate()
        .flat_map(|(row, line)| {
            line.chars()
                .enumerate()
                .map(move |(col, ch)| (row, col, ch))
        })
        .filter(|(row, col, ch)| {
            *ch == '@' && count_neighbors(data, *row, *col) < 4
        })
        .count() as i64
}

/// Counts the total number of paper rolls that can be removed through iterative processing.
/// Repeatedly removes accessible rolls (< 4 neighbors) until no more can be removed.
/// After each removal iteration, previously inaccessible rolls may become accessible.
fn part_two(data: &[String]) -> i64 {
    // Convert to mutable 2D array for efficient in-place mutation
    let mut grid: Vec<Vec<char>> = data
        .iter()
        .map(|line| line.chars().collect())
        .collect();

    let mut total_removed = 0;

    // Repeat until no more rolls can be removed
    loop {
        // Collect all currently accessible positions
        let mut removable = Vec::new();
        for (row, line) in grid.iter().enumerate() {
            for (col, &ch) in line.iter().enumerate() {
                if ch == '@' && count_neighbors_2d(&grid, row, col) < 4 {
                    removable.push((row, col));
                }
            }
        }

        // Exit condition: no more accessible rolls
        if removable.is_empty() {
            break;
        }

        // Remove all accessible rolls in this iteration
        for (row, col) in removable.iter() {
            grid[*row][*col] = '.';
        }

        total_removed += removable.len();
    }

    total_removed as i64
}

fn main() {
    let input = fs::read_to_string("input.txt")
        .expect("Failed to read input.txt");

    let data = parse_input(&input);

    println!("Part One: {}", part_one(&data));
    println!("Part Two: {}", part_two(&data));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 13);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 43);
    }
}
