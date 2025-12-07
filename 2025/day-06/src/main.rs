use std::fs;

struct Problem {
    numbers: Vec<i64>,
    operator: char,
}

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

fn parse_problems(data: &[String]) -> Vec<Problem> {
    if data.is_empty() {
        return Vec::new();
    }

    // Convert to character grid
    let grid: Vec<Vec<char>> = data.iter().map(|line| line.chars().collect()).collect();
    let num_rows = grid.len();
    let num_cols = grid.iter().map(|row| row.len()).max().unwrap_or(0);

    // Find column boundaries by identifying separator columns (all spaces)
    let mut is_separator_col = vec![true; num_cols];
    for col in 0..num_cols {
        for row in &grid {
            if col < row.len() && row[col] != ' ' {
                is_separator_col[col] = false;
                break;
            }
        }
    }

    // Group consecutive non-separator columns into problem ranges
    let mut problem_ranges: Vec<(usize, usize)> = Vec::new();
    let mut start: Option<usize> = None;

    for (col, &is_sep) in is_separator_col.iter().enumerate() {
        if !is_sep {
            if start.is_none() {
                start = Some(col);
            }
        } else if let Some(s) = start {
            problem_ranges.push((s, col));
            start = None;
        }
    }
    if let Some(s) = start {
        problem_ranges.push((s, num_cols));
    }

    // Extract problems from each column range
    let mut problems = Vec::new();
    for (start_col, end_col) in problem_ranges {
        // Get operator from last row
        let operator_row = &grid[num_rows - 1];
        let mut operator = '+';
        for col in start_col..end_col {
            if col < operator_row.len() {
                let ch = operator_row[col];
                if ch == '*' || ch == '+' {
                    operator = ch;
                    break;
                }
            }
        }

        // Get numbers from rows above the operator row
        let mut numbers = Vec::new();
        for row in grid.iter().take(num_rows - 1) {
            let slice: String = (start_col..end_col)
                .map(|col| if col < row.len() { row[col] } else { ' ' })
                .collect();
            if let Ok(num) = slice.trim().parse::<i64>() {
                numbers.push(num);
            }
        }

        problems.push(Problem { numbers, operator });
    }

    problems
}

fn solve_problem(problem: &Problem) -> i64 {
    match problem.operator {
        '*' => problem.numbers.iter().product(),
        '+' => problem.numbers.iter().sum(),
        _ => 0,
    }
}

fn part_one(data: &[String]) -> i64 {
    let problems = parse_problems(data);
    problems.iter().map(solve_problem).sum()
}

fn parse_problems_v2(data: &[String]) -> Vec<Problem> {
    if data.is_empty() {
        return Vec::new();
    }

    // Convert to character grid
    let grid: Vec<Vec<char>> = data.iter().map(|line| line.chars().collect()).collect();
    let num_rows = grid.len();
    let num_cols = grid.iter().map(|row| row.len()).max().unwrap_or(0);

    // Find column boundaries by identifying separator columns (all spaces)
    let mut is_separator_col = vec![true; num_cols];
    for col in 0..num_cols {
        for row in &grid {
            if col < row.len() && row[col] != ' ' {
                is_separator_col[col] = false;
                break;
            }
        }
    }

    // Group consecutive non-separator columns into problem ranges
    let mut problem_ranges: Vec<(usize, usize)> = Vec::new();
    let mut start: Option<usize> = None;

    for (col, &is_sep) in is_separator_col.iter().enumerate() {
        if !is_sep {
            if start.is_none() {
                start = Some(col);
            }
        } else if let Some(s) = start {
            problem_ranges.push((s, col));
            start = None;
        }
    }
    if let Some(s) = start {
        problem_ranges.push((s, num_cols));
    }

    // Extract problems from each column range
    let mut problems = Vec::new();
    for (start_col, end_col) in problem_ranges {
        // Get operator from last row
        let operator_row = &grid[num_rows - 1];
        let mut operator = '+';
        for col in start_col..end_col {
            if col < operator_row.len() {
                let ch = operator_row[col];
                if ch == '*' || ch == '+' {
                    operator = ch;
                    break;
                }
            }
        }

        // Get numbers by reading each column vertically (digits top-to-bottom)
        let mut numbers = Vec::new();
        for col in start_col..end_col {
            let mut digits = String::new();
            for row in grid.iter().take(num_rows - 1) {
                if col < row.len() {
                    let ch = row[col];
                    if ch.is_ascii_digit() {
                        digits.push(ch);
                    }
                }
            }
            if !digits.is_empty() {
                if let Ok(num) = digits.parse::<i64>() {
                    numbers.push(num);
                }
            }
        }

        problems.push(Problem { numbers, operator });
    }

    problems
}

fn part_two(data: &[String]) -> i64 {
    let problems = parse_problems_v2(data);
    problems.iter().map(solve_problem).sum()
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
123 328  51 64
 45 64  387 23
  6 98  215 314
*   +   *   +  ";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        // Expected: 123*45*6=33210, 328+64+98=490, 51*387*215=4243455, 64+23+314=401
        // Grand total: 33210 + 490 + 4243455 + 401 = 4277556
        assert_eq!(part_one(&data), 4277556);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        // Reading right-to-left, one column at a time:
        // Rightmost: 4 + 431 + 623 = 1058
        // Second from right: 175 * 581 * 32 = 3253600
        // Third from right: 8 + 248 + 369 = 625
        // Leftmost: 356 * 24 * 1 = 8544
        // Grand total: 1058 + 3253600 + 625 + 8544 = 3263827
        assert_eq!(part_two(&data), 3263827);
    }
}
