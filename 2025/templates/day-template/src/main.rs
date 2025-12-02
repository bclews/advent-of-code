use std::fs;

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

fn part_one(data: &[String]) -> i64 {
    // TODO: Implement part one
    0
}

fn part_two(data: &[String]) -> i64 {
    // TODO: Implement part two
    0
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
TODO: Add example input from puzzle";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 0); // TODO: Update expected value
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 0); // TODO: Update expected value
    }
}
