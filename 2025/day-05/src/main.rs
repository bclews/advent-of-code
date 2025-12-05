use std::fs;

struct CafeteriaData {
    ranges: Vec<(i64, i64)>,  // Fresh ingredient ID ranges
    ingredient_ids: Vec<i64>, // Available ingredients to check
}

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

fn parse_ranges_section(lines: &[String]) -> Vec<(i64, i64)> {
    lines
        .iter()
        .filter_map(|line| {
            let parts: Vec<&str> = line.trim().split('-').collect();
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

fn parse_cafeteria_data(data: &[String]) -> CafeteriaData {
    // Find the blank line that separates the two sections
    let blank_pos = data.iter().position(|s| s.is_empty()).unwrap_or(data.len());

    // Parse ranges from the first section
    let ranges = parse_ranges_section(&data[..blank_pos]);

    // Parse ingredient IDs from the second section
    let ingredient_ids: Vec<i64> = data[blank_pos + 1..]
        .iter()
        .filter_map(|line| line.trim().parse::<i64>().ok())
        .collect();

    CafeteriaData {
        ranges,
        ingredient_ids,
    }
}

fn is_fresh(id: i64, ranges: &[(i64, i64)]) -> bool {
    ranges.iter().any(|(start, end)| id >= *start && id <= *end)
}

fn merge_ranges(ranges: &[(i64, i64)]) -> Vec<(i64, i64)> {
    if ranges.is_empty() {
        return Vec::new();
    }

    // Sort ranges by start position
    let mut sorted_ranges = ranges.to_vec();
    sorted_ranges.sort_by_key(|(start, _)| *start);

    let mut merged = Vec::new();
    let mut current = sorted_ranges[0];

    for &(start, end) in &sorted_ranges[1..] {
        // Check if current range overlaps or is adjacent to the next range
        if start <= current.1 + 1 {
            // Merge: extend current range's end to the maximum of both ends
            current.1 = current.1.max(end);
        } else {
            // No overlap: add current range to result and start a new one
            merged.push(current);
            current = (start, end);
        }
    }

    // Add the last range
    merged.push(current);
    merged
}

fn part_one(data: &[String]) -> i64 {
    let cafeteria = parse_cafeteria_data(data);

    cafeteria
        .ingredient_ids
        .iter()
        .filter(|&id| is_fresh(*id, &cafeteria.ranges))
        .count() as i64
}

fn part_two(data: &[String]) -> i64 {
    let cafeteria = parse_cafeteria_data(data);

    // Merge overlapping ranges first
    let merged = merge_ranges(&cafeteria.ranges);

    // Calculate total count by summing the sizes of merged ranges
    merged.iter().map(|(start, end)| end - start + 1).sum()
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

    const EXAMPLE: &str = "\
3-5
10-14
16-20
12-18
        w 

1
5
8
11
17
32";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert!(!data.is_empty());
    }

    #[test]
    fn test_parse_cafeteria_data() {
        let data = parse_input(EXAMPLE);
        let cafeteria = parse_cafeteria_data(&data);
        assert_eq!(cafeteria.ranges, vec![(3, 5), (10, 14), (16, 20), (12, 18)]);
        assert_eq!(cafeteria.ingredient_ids, vec![1, 5, 8, 11, 17, 32]);
    }

    #[test]
    fn test_is_fresh() {
        let ranges = vec![(3, 5), (10, 14), (16, 20), (12, 18)];
        // Test fresh ingredients
        assert!(is_fresh(5, &ranges));
        assert!(is_fresh(11, &ranges));
        assert!(is_fresh(17, &ranges));
        // Test spoiled ingredients
        assert!(!is_fresh(1, &ranges));
        assert!(!is_fresh(8, &ranges));
        assert!(!is_fresh(32, &ranges));
        // Test edge cases (range boundaries)
        assert!(is_fresh(3, &ranges));
        assert!(is_fresh(14, &ranges));
        assert!(is_fresh(16, &ranges));
        assert!(is_fresh(20, &ranges));
    }

    #[test]
    fn test_merge_ranges_overlapping() {
        // Test overlapping ranges
        let ranges = vec![(3, 5), (10, 14), (16, 20), (12, 18)];
        let merged = merge_ranges(&ranges);
        assert_eq!(merged, vec![(3, 5), (10, 20)]);
    }

    #[test]
    fn test_merge_ranges_adjacent() {
        // Test adjacent ranges that should merge
        let ranges = vec![(1, 5), (6, 10), (11, 15)];
        let merged = merge_ranges(&ranges);
        assert_eq!(merged, vec![(1, 15)]);
    }

    #[test]
    fn test_merge_ranges_non_overlapping() {
        // Test non-overlapping ranges that should stay separate
        let ranges = vec![(1, 5), (10, 15), (20, 25)];
        let merged = merge_ranges(&ranges);
        assert_eq!(merged, vec![(1, 5), (10, 15), (20, 25)]);
    }

    #[test]
    fn test_merge_ranges_fully_contained() {
        // Test fully contained ranges
        let ranges = vec![(1, 20), (5, 10), (8, 12)];
        let merged = merge_ranges(&ranges);
        assert_eq!(merged, vec![(1, 20)]);
    }

    #[test]
    fn test_merge_ranges_empty() {
        // Test empty input
        let ranges: Vec<(i64, i64)> = vec![];
        let merged = merge_ranges(&ranges);
        assert_eq!(merged, vec![]);
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 3); // 3 fresh ingredients: 5, 11, 17
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 14); // 14 unique fresh IDs: 3,4,5,10,11,12,13,14,15,16,17,18,19,20
    }
}
