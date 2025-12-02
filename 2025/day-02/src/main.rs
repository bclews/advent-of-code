use std::collections::HashSet;
use std::fs;

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

// Parse comma-separated ranges like "11-22,95-115" into Vec<(start, end)>
fn parse_ranges(input: &str) -> Vec<(i64, i64)> {
    input
        .split(',')
        .filter_map(|range| {
            let parts: Vec<&str> = range.trim().split('-').collect();
            if parts.len() == 2 {
                let start = parts[0].parse::<i64>().ok()?;
                let end = parts[1].parse::<i64>().ok()?;
                Some((start, end))
            } else {
                None
            }
        })
        .collect()
}

// Check if a number is an invalid ID (sequence repeated twice)
// E.g., 11 = "1"+"1", 6464 = "64"+"64", 123123 = "123"+"123"
#[cfg(test)]
fn is_invalid_id(n: i64) -> bool {
    let s = n.to_string();
    let len = s.len();

    // Must have even length to split in half
    if !len.is_multiple_of(2) {
        return false;
    }

    let mid = len / 2;
    let first_half = &s[..mid];
    let second_half = &s[mid..];

    // Check if both halves are identical
    first_half == second_half
}

// Check if a value falls within any of the given ranges
fn is_in_any_range(value: i64, ranges: &[(i64, i64)]) -> bool {
    ranges
        .iter()
        .any(|(start, end)| value >= *start && value <= *end)
}

// Generate invalid IDs using mathematical formula: S × (10^d + 1)
// where S is a d-digit sequence. This avoids brute-forcing through ranges.
// E.g., sequence 64 (2 digits) → 64 × 101 = 6464
fn generate_invalid_ids_in_bounds(min: i64, max: i64, ranges: &[(i64, i64)]) -> Vec<i64> {
    let mut invalid_ids = Vec::new();

    // Determine how many digits we need to check based on max value
    let max_digits = max.to_string().len() as u32;

    // For each possible digit count (1, 2, 3, ...)
    for d in 1..=max_digits {
        // Multiplier for d digits: 10^d + 1
        // E.g., d=2 → multiplier=101, d=3 → multiplier=1001
        let multiplier = 10_i64.pow(d) + 1;

        // Start sequence: 10^(d-1) or 1 for single digit
        // E.g., d=2 starts at 10, d=3 starts at 100
        let start_seq = if d == 1 { 1 } else { 10_i64.pow(d - 1) };

        // End sequence: limited by either 10^d - 1 or max/multiplier (to avoid overflow)
        // E.g., d=2 ends at 99, d=3 ends at 999
        let theoretical_end = 10_i64.pow(d) - 1;
        let max_safe_seq = max / multiplier; // Avoid overflow by limiting sequence range
        let end_seq = theoretical_end.min(max_safe_seq);

        // Generate invalid IDs for all sequences of this length
        for seq in start_seq..=end_seq {
            let invalid_id = seq * multiplier;

            // Early termination: if we exceed max, all subsequent will too
            if invalid_id > max {
                break;
            }

            // Check if this invalid ID falls in any of our ranges
            if invalid_id >= min && is_in_any_range(invalid_id, ranges) {
                invalid_ids.push(invalid_id);
            }
        }
    }

    invalid_ids
}

// Find all invalid IDs in a single range (helper for tests)
#[cfg(test)]
fn find_invalid_ids_in_range(start: i64, end: i64) -> Vec<i64> {
    let ranges = vec![(start, end)];
    generate_invalid_ids_in_bounds(start, end, &ranges)
}

// Check if a number is invalid according to Part Two rules (sequence repeated 2+ times)
// E.g., 111 = "1"×3, 12341234 = "1234"×2, 1212121212 = "12"×5
#[cfg(test)]
fn is_invalid_id_part_two(n: i64) -> bool {
    let s = n.to_string();
    let len = s.len();

    // Try all possible sequence lengths from 1 to len/2
    // If len is divisible by sequence length, check if it's a repeated pattern
    for seq_len in 1..=(len / 2) {
        if len.is_multiple_of(seq_len) {
            let sequence = &s[..seq_len];
            let repetitions = len / seq_len;

            // Check if the entire string is this sequence repeated
            let mut is_repeated = true;
            for i in 1..repetitions {
                let start = i * seq_len;
                let end = start + seq_len;
                if &s[start..end] != sequence {
                    is_repeated = false;
                    break;
                }
            }

            // If we found a valid repetition (2+ times), it's invalid
            if is_repeated && repetitions >= 2 {
                return true;
            }
        }
    }

    false
}

// Find all invalid IDs in a single range for Part Two (helper for tests)
#[cfg(test)]
fn find_invalid_ids_in_range_part_two(start: i64, end: i64) -> Vec<i64> {
    let ranges = vec![(start, end)];
    generate_invalid_ids_in_bounds_part_two(start, end, &ranges)
}

fn part_one(data: &[String]) -> i64 {
    // Parse the first line containing comma-separated ranges
    let ranges = parse_ranges(&data[0]);

    // Find min and max bounds across all ranges
    let min = ranges.iter().map(|(s, _)| *s).min().unwrap_or(0);
    let max = ranges.iter().map(|(_, e)| *e).max().unwrap_or(0);

    // Generate all invalid IDs that fall within any range
    let invalid_ids = generate_invalid_ids_in_bounds(min, max, &ranges);

    // Sum all invalid IDs
    invalid_ids.iter().sum()
}

// Generate invalid IDs for Part Two using mathematical formula
// Part 2: sequences repeated 2+ times (not just exactly twice)
// Formula: For d-digit sequence S repeated r times:
//   Invalid ID = S × ((10^(d×r) - 1) / (10^d - 1))
// This is a geometric series: S × (1 + 10^d + 10^(2d) + ... + 10^((r-1)d))
fn generate_invalid_ids_in_bounds_part_two(min: i64, max: i64, ranges: &[(i64, i64)]) -> Vec<i64> {
    let mut invalid_ids = HashSet::new();

    // Determine how many digits we need to check based on max value
    let max_digits = max.to_string().len() as u32;

    // For each possible digit count (1, 2, 3, ...)
    for d in 1..=max_digits {
        // Base for geometric series: 10^d
        // E.g., d=2 → base=100, d=3 → base=1000
        let base = 10_i64.pow(d);

        // Start sequence: 10^(d-1) or 1 for single digit
        let start_seq = if d == 1 { 1 } else { 10_i64.pow(d - 1) };

        // End sequence: 10^d - 1
        let theoretical_end = base - 1;

        // For each sequence S of d digits
        for seq in start_seq..=theoretical_end {
            // Quick check: smallest possible invalid ID (seq repeated 2x)
            // If this is already > max, we can skip this entire sequence
            let min_possible = match seq.checked_mul(base + 1) {
                Some(val) => val,
                None => {
                    break; // Overflow means definitely > max
                }
            };
            if min_possible > max {
                break; // All subsequent sequences will also exceed max
            }

            // Try different repetition counts (2, 3, 4, ...)
            // We need at least 2 repetitions to be invalid
            let mut repetitions = 2;

            // Use while let to cleanly handle overflow
            while let Some(base_power) = base.checked_pow(repetitions) {
                // Calculate the geometric series multiplier
                // multiplier = (base^repetitions - 1) / (base - 1)
                // E.g., for "12" (d=2, base=100) repeated 3 times:
                //   multiplier = (100^3 - 1) / (100 - 1) = 999999 / 99 = 10101
                let multiplier = (base_power - 1) / (base - 1);

                // Calculate invalid ID
                let invalid_id = match seq.checked_mul(multiplier) {
                    Some(val) => val,
                    None => break, // Overflow, stop trying more repetitions
                };

                // Early termination: if we exceed max, all subsequent reps will too
                if invalid_id > max {
                    break;
                }

                // Check if this invalid ID falls in any of our ranges
                if invalid_id >= min && is_in_any_range(invalid_id, ranges) {
                    invalid_ids.insert(invalid_id);
                }

                repetitions += 1;

                // Safety limit: stop at reasonable repetition count
                if repetitions > 20 {
                    break;
                }
            }
        }
    }

    // Convert to sorted Vec
    let mut result: Vec<i64> = invalid_ids.into_iter().collect();
    result.sort_unstable();
    result
}

fn part_two(data: &[String]) -> i64 {
    // Parse the first line containing comma-separated ranges
    let ranges = parse_ranges(&data[0]);

    // Find min and max bounds across all ranges
    let min = ranges.iter().map(|(s, _)| *s).min().unwrap_or(0);
    let max = ranges.iter().map(|(_, e)| *e).max().unwrap_or(0);

    // Generate all invalid IDs (sequences repeated 2+ times) that fall within any range
    let invalid_ids = generate_invalid_ids_in_bounds_part_two(min, max, &ranges);

    // Sum all invalid IDs
    invalid_ids.iter().sum()
}

fn main() {
    let input = fs::read_to_string("input.txt").expect("Failed to read input.txt");

    let data = parse_input(&input);

    println!("Part One: {}", part_one(&data));
    println!("Part Two: {}", part_two(&data));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_is_invalid_id() {
        // Test cases from problem: numbers that are a sequence repeated twice
        assert!(is_invalid_id(11), "11 should be invalid (1 repeated twice)");
        assert!(is_invalid_id(22), "22 should be invalid (2 repeated twice)");
        assert!(is_invalid_id(55), "55 should be invalid (5 repeated twice)");
        assert!(is_invalid_id(99), "99 should be invalid (9 repeated twice)");
        assert!(
            is_invalid_id(6464),
            "6464 should be invalid (64 repeated twice)"
        );
        assert!(
            is_invalid_id(123123),
            "123123 should be invalid (123 repeated twice)"
        );
        assert!(
            is_invalid_id(1010),
            "1010 should be invalid (10 repeated twice)"
        );
        assert!(
            is_invalid_id(222222),
            "222222 should be invalid (222 repeated twice)"
        );
        assert!(
            is_invalid_id(446446),
            "446446 should be invalid (446 repeated twice)"
        );
        assert!(
            is_invalid_id(38593859),
            "38593859 should be invalid (3859 repeated twice)"
        );
        assert!(
            is_invalid_id(1188511885),
            "1188511885 should be invalid (11885 repeated twice)"
        );

        // Test valid IDs (not repeated patterns)
        assert!(
            !is_invalid_id(101),
            "101 should be valid (not a repeated pattern)"
        );
        assert!(!is_invalid_id(1698522), "1698522 should be valid");
        assert!(!is_invalid_id(1698528), "1698528 should be valid");
    }

    #[test]
    fn test_find_invalid_ids_in_range() {
        // From problem examples:
        // 11-22 has two invalid IDs: 11 and 22
        let ids = find_invalid_ids_in_range(11, 22);
        assert_eq!(ids, vec![11, 22]);

        // 95-115 has one invalid ID: 99
        let ids = find_invalid_ids_in_range(95, 115);
        assert_eq!(ids, vec![99]);

        // 998-1012 has one invalid ID: 1010
        let ids = find_invalid_ids_in_range(998, 1012);
        assert_eq!(ids, vec![1010]);

        // 1188511880-1188511890 has one invalid ID: 1188511885
        let ids = find_invalid_ids_in_range(1188511880, 1188511890);
        assert_eq!(ids, vec![1188511885]);

        // 222220-222224 has one invalid ID: 222222
        let ids = find_invalid_ids_in_range(222220, 222224);
        assert_eq!(ids, vec![222222]);

        // 1698522-1698528 contains no invalid IDs
        let ids = find_invalid_ids_in_range(1698522, 1698528);
        assert_eq!(ids, Vec::<i64>::new());

        // 446443-446449 has one invalid ID: 446446
        let ids = find_invalid_ids_in_range(446443, 446449);
        assert_eq!(ids, vec![446446]);

        // 38593856-38593862 has one invalid ID: 38593859
        let ids = find_invalid_ids_in_range(38593856, 38593862);
        assert_eq!(ids, vec![38593859]);
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        // Sum of all invalid IDs in example = 1227775554
        assert_eq!(part_one(&data), 1227775554);
    }

    // Part Two Tests - sequences repeated at least twice (2 or more times)

    #[test]
    fn test_is_invalid_id_part_two() {
        // Part 2: Invalid IDs are sequences repeated 2+ times

        // Sequences repeated exactly twice (still invalid in part 2)
        assert!(
            is_invalid_id_part_two(11),
            "11 should be invalid (1 repeated 2 times)"
        );
        assert!(
            is_invalid_id_part_two(22),
            "22 should be invalid (2 repeated 2 times)"
        );
        assert!(
            is_invalid_id_part_two(6464),
            "6464 should be invalid (64 repeated 2 times)"
        );
        assert!(
            is_invalid_id_part_two(123123),
            "123123 should be invalid (123 repeated 2 times)"
        );
        assert!(
            is_invalid_id_part_two(12341234),
            "12341234 should be invalid (1234 repeated 2 times)"
        );

        // Sequences repeated three times
        assert!(
            is_invalid_id_part_two(111),
            "111 should be invalid (1 repeated 3 times)"
        );
        assert!(
            is_invalid_id_part_two(999),
            "999 should be invalid (9 repeated 3 times)"
        );
        assert!(
            is_invalid_id_part_two(123123123),
            "123123123 should be invalid (123 repeated 3 times)"
        );
        assert!(
            is_invalid_id_part_two(565656),
            "565656 should be invalid (56 repeated 3 times)"
        );
        assert!(
            is_invalid_id_part_two(824824824),
            "824824824 should be invalid (824 repeated 3 times)"
        );

        // Sequences repeated five times
        assert!(
            is_invalid_id_part_two(1212121212),
            "1212121212 should be invalid (12 repeated 5 times)"
        );
        assert!(
            is_invalid_id_part_two(2121212121),
            "2121212121 should be invalid (21 repeated 5 times)"
        );

        // Sequences repeated seven times
        assert!(
            is_invalid_id_part_two(1111111),
            "1111111 should be invalid (1 repeated 7 times)"
        );

        // Valid IDs (not repeated patterns)
        assert!(
            !is_invalid_id_part_two(101),
            "101 should be valid (not a repeated pattern)"
        );
        assert!(!is_invalid_id_part_two(1698522), "1698522 should be valid");
        assert!(!is_invalid_id_part_two(1698528), "1698528 should be valid");
    }

    #[test]
    fn test_find_invalid_ids_in_range_part_two() {
        // From problem part 2 examples:

        // 11-22 still has two invalid IDs: 11 and 22
        let ids = find_invalid_ids_in_range_part_two(11, 22);
        assert_eq!(ids, vec![11, 22]);

        // 95-115 now has two invalid IDs: 99 and 111
        let ids = find_invalid_ids_in_range_part_two(95, 115);
        assert_eq!(ids, vec![99, 111]);

        // 998-1012 now has two invalid IDs: 999 and 1010
        let ids = find_invalid_ids_in_range_part_two(998, 1012);
        assert_eq!(ids, vec![999, 1010]);

        // 1188511880-1188511890 still has one invalid ID: 1188511885
        let ids = find_invalid_ids_in_range_part_two(1188511880, 1188511890);
        assert_eq!(ids, vec![1188511885]);

        // 222220-222224 still has one invalid ID: 222222
        let ids = find_invalid_ids_in_range_part_two(222220, 222224);
        assert_eq!(ids, vec![222222]);

        // 1698522-1698528 still contains no invalid IDs
        let ids = find_invalid_ids_in_range_part_two(1698522, 1698528);
        assert_eq!(ids, Vec::<i64>::new());

        // 446443-446449 still has one invalid ID: 446446
        let ids = find_invalid_ids_in_range_part_two(446443, 446449);
        assert_eq!(ids, vec![446446]);

        // 38593856-38593862 still has one invalid ID: 38593859
        let ids = find_invalid_ids_in_range_part_two(38593856, 38593862);
        assert_eq!(ids, vec![38593859]);

        // 565653-565659 now has one invalid ID: 565656
        let ids = find_invalid_ids_in_range_part_two(565653, 565659);
        assert_eq!(ids, vec![565656]);

        // 824824821-824824827 now has one invalid ID: 824824824
        let ids = find_invalid_ids_in_range_part_two(824824821, 824824827);
        assert_eq!(ids, vec![824824824]);

        // 2121212118-2121212124 now has one invalid ID: 2121212121
        let ids = find_invalid_ids_in_range_part_two(2121212118, 2121212124);
        assert_eq!(ids, vec![2121212121]);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        // Sum of all invalid IDs in part 2 example = 4174379265
        assert_eq!(part_two(&data), 4174379265);
    }
}
