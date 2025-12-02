use std::fs;

const DIAL_START: i32 = 50;
const DIAL_SIZE: i32 = 100;

fn parse_input(input: &str) -> Result<Vec<(char, i32)>, Box<dyn std::error::Error>> {
    input
        .lines()
        .filter(|line| !line.is_empty())
        .map(|line| {
            let direction = line.chars().next()
                .ok_or("Empty line after filtering")?;
            let distance = line.get(1..)
                .ok_or("Line too short")?
                .parse::<i32>()?;
            Ok((direction, distance))
        })
        .collect()
}

fn part_one(rotations: &[(char, i32)]) -> i64 {
    let mut position: i32 = DIAL_START;
    let mut count = 0;

    for &(direction, distance) in rotations {
        match direction {
            'L' => {
                // Left means toward lower numbers (subtract)
                position = (position - distance).rem_euclid(DIAL_SIZE);
            }
            'R' => {
                // Right means toward higher numbers (add)
                position = (position + distance).rem_euclid(DIAL_SIZE);
            }
            _ => panic!("Invalid direction: {}", direction),
        }

        if position == 0 {
            count += 1;
        }
    }

    count
}

fn part_two(rotations: &[(char, i32)]) -> i64 {
    let mut position: i32 = DIAL_START;
    let mut count = 0;

    for &(direction, distance) in rotations {
        // Simulate each click one at a time
        for _ in 0..distance {
            match direction {
                'L' => {
                    position = (position - 1).rem_euclid(DIAL_SIZE);
                }
                'R' => {
                    position = (position + 1).rem_euclid(DIAL_SIZE);
                }
                _ => panic!("Invalid direction: {}", direction),
            }

            if position == 0 {
                count += 1;
            }
        }
    }

    count
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let input = fs::read_to_string("input.txt")?;
    let data = parse_input(&input)?;

    println!("Part One: {}", part_one(&data));
    println!("Part Two: {}", part_two(&data));

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
L68
L30
R48
L5
R60
L55
L1
L99
R14
L82";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE).unwrap();
        assert_eq!(data.len(), 10);
        assert_eq!(data[0], ('L', 68));
        assert_eq!(data[2], ('R', 48));
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE).unwrap();
        // Expected answer: 3 (dial points at 0 after rotations R48, L55, and L99)
        assert_eq!(part_one(&data), 3);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE).unwrap();
        // Expected answer: 6 (3 times at end of rotation + 3 times during rotation)
        // During rotations that cross 0:
        // - L68 from 50 -> 82 crosses 0 once
        // - R60 from 95 -> 55 crosses 0 once
        // - L82 from 14 -> 32 crosses 0 once
        // Plus the 3 times from part one (ending at 0)
        assert_eq!(part_two(&data), 6);
    }
}
