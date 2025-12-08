use std::fs;
use std::collections::{VecDeque, HashSet, HashMap};

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

fn find_start(grid: &[Vec<char>]) -> Option<(usize, usize)> {
    for (row, line) in grid.iter().enumerate() {
        for (col, &ch) in line.iter().enumerate() {
            if ch == 'S' {
                return Some((row, col));
            }
        }
    }
    None
}

fn simulate_beams(grid: &[Vec<char>], start: (usize, usize)) -> i64 {
    let rows = grid.len() as i32;
    let cols = grid[0].len() as i32;

    let mut queue: VecDeque<(i32, i32)> = VecDeque::new();
    let mut visited: HashSet<(i32, i32)> = HashSet::new();
    let mut split_count = 0;

    // Start with initial beam at starting position
    queue.push_back((start.0 as i32, start.1 as i32));

    while let Some((row, col)) = queue.pop_front() {
        // Check bounds
        if row < 0 || row >= rows || col < 0 || col >= cols {
            continue;
        }

        // Skip if already visited
        if visited.contains(&(row, col)) {
            continue;
        }

        // Mark as visited
        visited.insert((row, col));

        // Get current cell
        let current = grid[row as usize][col as usize];

        // Process based on cell type
        match current {
            '.' | 'S' => {
                // Continue beam downward
                queue.push_back((row + 1, col));
            }
            '^' => {
                // Splitter: increment count and create two beams
                split_count += 1;
                // Left beam continues downward
                queue.push_back((row + 1, col - 1));
                // Right beam continues downward
                queue.push_back((row + 1, col + 1));
            }
            _ => {}
        }
    }

    split_count
}

fn part_one(data: &[String]) -> i64 {
    // Convert to grid
    let grid: Vec<Vec<char>> = data.iter().map(|line| line.chars().collect()).collect();

    // Find starting position
    if let Some(start) = find_start(&grid) {
        simulate_beams(&grid, start)
    } else {
        0
    }
}

fn count_timelines_memo(
    grid: &[Vec<char>],
    row: i32,
    col: i32,
    rows: i32,
    cols: i32,
    memo: &mut HashMap<(i32, i32), i64>
) -> i64 {
    // Base case: exited the grid (complete path)
    if row >= rows || col < 0 || col >= cols {
        return 1;  // This path is complete, count it as one timeline
    }

    // Check memoization cache
    if let Some(&cached) = memo.get(&(row, col)) {
        return cached;
    }

    // Get current cell
    let current = grid[row as usize][col as usize];

    // Calculate result based on cell type
    let result = match current {
        '.' | 'S' => {
            // Continue downward
            count_timelines_memo(grid, row + 1, col, rows, cols, memo)
        }
        '^' => {
            // Splitter: quantum particle takes BOTH paths
            let left_timelines = count_timelines_memo(grid, row + 1, col - 1, rows, cols, memo);
            let right_timelines = count_timelines_memo(grid, row + 1, col + 1, rows, cols, memo);
            left_timelines + right_timelines
        }
        _ => 0
    };

    // Cache the result
    memo.insert((row, col), result);
    result
}

fn part_two(data: &[String]) -> i64 {
    // Convert to grid
    let grid: Vec<Vec<char>> = data.iter().map(|line| line.chars().collect()).collect();

    // Find starting position
    if let Some(start) = find_start(&grid) {
        let rows = grid.len() as i32;
        let cols = grid[0].len() as i32;
        let mut memo = HashMap::new();
        count_timelines_memo(&grid, start.0 as i32, start.1 as i32, rows, cols, &mut memo)
    } else {
        0
    }
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
.......S.......
...............
.......^.......
...............
......^.^......
...............
.....^.^.^.....
...............
....^.^...^....
...............
...^.^...^.^...
...............
..^...^.....^..
...............
.^.^.^.^.^...^.
...............";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 21); // 21 beam splits in the example
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 40); // 40 different timelines in the example
    }
}
