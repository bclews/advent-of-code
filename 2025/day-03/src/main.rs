use std::fs;

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

fn find_max_joltage(bank: &str) -> u32 {
    let digits: Vec<u32> = bank.chars().filter_map(|c| c.to_digit(10)).collect();

    if digits.len() < 2 {
        return 0;
    }

    let mut max_joltage = 0;
    let mut max_first_digit = digits[0];

    for &digit in digits.iter().skip(1) {
        // Best joltage if current digit is the second digit
        let joltage = max_first_digit * 10 + digit;
        max_joltage = max_joltage.max(joltage);

        // Update max first digit for future iterations
        max_first_digit = max_first_digit.max(digit);
    }

    max_joltage
}

fn part_one(data: &[String]) -> u32 {
    data.iter().map(|bank| find_max_joltage(bank)).sum()
}

// Monotonic stack algorithm: Select k digits from n to form maximum k-digit number
// Example: "234234234234278" with k=12 → "434234234278"
// Strategy: Build result left-to-right, greedily keeping larger digits when possible
fn find_max_k_digit_joltage(bank: &str, k: usize) -> i64 {
    let digits: Vec<char> = bank.chars().collect();
    let n = digits.len();

    if n < k {
        return 0;
    }

    let to_skip = n - k;
    let mut stack: Vec<char> = Vec::with_capacity(k);
    let mut skips_remaining = to_skip;

    for (i, &digit) in digits.iter().enumerate() {
        let remaining = n - i - 1;

        // Pop smaller digits from stack to make room for larger ones
        // Example: stack=['2'], digit='3' → pop '2', push '3' (improves result)
        // Constraints ensure we maintain exactly k digits:
        //   - stack top < current digit (can improve)
        //   - skips_remaining > 0 (allowed to skip)
        //   - stack.len() + remaining >= k (won't run out of digits)
        while !stack.is_empty()
            && stack[stack.len() - 1] < digit
            && skips_remaining > 0
            && stack.len() + remaining >= k
        {
            stack.pop();
            skips_remaining -= 1;
        }

        stack.push(digit);
    }

    stack.truncate(k);
    stack.iter().collect::<String>().parse::<i64>().unwrap_or(0)
}

fn part_two(data: &[String]) -> i64 {
    data.iter()
        .map(|bank| find_max_k_digit_joltage(bank, 12))
        .sum()
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

    const EXAMPLE: &str = "987654321111111
811111111111119
234234234234278
818181911112111";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert_eq!(data.len(), 4);
        assert_eq!(data[0], "987654321111111");
    }

    #[test]
    fn test_max_joltage_individual_banks() {
        assert_eq!(find_max_joltage("987654321111111"), 98);
        assert_eq!(find_max_joltage("811111111111119"), 89);
        assert_eq!(find_max_joltage("234234234234278"), 78);
        assert_eq!(find_max_joltage("818181911112111"), 92);
    }

    #[test]
    fn test_total_output_joltage() {
        let data = parse_input(EXAMPLE);
        let total: u32 = data.iter().map(|bank| find_max_joltage(bank)).sum();
        assert_eq!(total, 357);
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 357);
    }

    #[test]
    fn test_find_max_k_digit_joltage_individual_banks() {
        // Each bank has 15 digits, we select 12
        assert_eq!(
            find_max_k_digit_joltage("987654321111111", 12),
            987654321111
        );
        assert_eq!(
            find_max_k_digit_joltage("811111111111119", 12),
            811111111119
        );
        assert_eq!(
            find_max_k_digit_joltage("234234234234278", 12),
            434234234278
        );
        assert_eq!(
            find_max_k_digit_joltage("818181911112111", 12),
            888911112111
        );
    }

    #[test]
    fn test_total_k_digit_joltage() {
        let data = parse_input(EXAMPLE);
        let total: i64 = data
            .iter()
            .map(|bank| find_max_k_digit_joltage(bank, 12))
            .sum();
        assert_eq!(total, 3121910778619);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 3121910778619);
    }

    #[test]
    fn test_edge_cases() {
        // Maximum possible (99)
        assert_eq!(find_max_joltage("999999"), 99);

        // Minimum possible (12)
        assert_eq!(find_max_joltage("111112"), 12);

        // Ascending order
        assert_eq!(find_max_joltage("123456789"), 89);

        // Descending order
        assert_eq!(find_max_joltage("987654321"), 98);

        // All same digit
        assert_eq!(find_max_joltage("555555"), 55);

        // Two digits only
        assert_eq!(find_max_joltage("47"), 47);

        // Best answer at end
        assert_eq!(find_max_joltage("111119"), 19);

        // Best answer scattered
        assert_eq!(find_max_joltage("21937465"), 97);

        // Empty string
        assert_eq!(find_max_joltage(""), 0);

        // Single digit
        assert_eq!(find_max_joltage("7"), 0);
    }
}
