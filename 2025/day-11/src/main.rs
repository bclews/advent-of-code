use std::collections::HashMap;
use std::fs;

// --- TYPE ALIASES ---

type Graph = HashMap<String, Vec<String>>;
type MemoKey = (String, bool, bool);
type MemoCache = HashMap<MemoKey, i64>;

// --- PARSING ---

fn parse_input(input: &str) -> Graph {
    input
        .lines()
        .filter(|line| !line.is_empty())
        .map(|line| {
            let parts: Vec<&str> = line.split(':').collect();
            let device = parts[0].trim().to_string();
            let outputs = parts[1].split_whitespace().map(|s| s.to_string()).collect();
            (device, outputs)
        })
        .collect()
}

// --- PART ONE: SIMPLE PATH COUNTING ---

/// Counts all paths from current node to target using DFS
fn count_paths(graph: &Graph, current: &str, target: &str) -> i64 {
    if current == target {
        return 1;
    }

    graph
        .get(current)
        .map(|neighbors| {
            neighbors
                .iter()
                .map(|next| count_paths(graph, next, target))
                .sum()
        })
        .unwrap_or(0)
}

// --- PART TWO: PATH COUNTING WITH REQUIRED NODE VISITS ---

/// Updates visited flags based on current node
#[inline]
fn update_visited_flags(current: &str, visited_dac: bool, visited_fft: bool) -> (bool, bool) {
    (
        visited_dac || current == "dac",
        visited_fft || current == "fft",
    )
}

/// Creates a cache key from current state
#[inline]
fn make_cache_key(node: &str, visited_dac: bool, visited_fft: bool) -> MemoKey {
    (node.to_string(), visited_dac, visited_fft)
}

/// Checks if both required nodes have been visited
#[inline]
fn both_required_visited(visited_dac: bool, visited_fft: bool) -> bool {
    visited_dac && visited_fft
}

/// Counts paths with memoization, visiting both "dac" and "fft"
fn count_paths_memoized(
    graph: &Graph,
    current: &str,
    target: &str,
    visited_dac: bool,
    visited_fft: bool,
    memo: &mut MemoCache,
) -> i64 {
    let (new_visited_dac, new_visited_fft) =
        update_visited_flags(current, visited_dac, visited_fft);

    // Base case: reached target, only count if both required nodes visited
    if current == target {
        return if both_required_visited(new_visited_dac, new_visited_fft) {
            1
        } else {
            0
        };
    }

    // Check cache
    let cache_key = make_cache_key(current, new_visited_dac, new_visited_fft);
    if let Some(&cached_result) = memo.get(&cache_key) {
        return cached_result;
    }

    // Compute: sum paths through all neighbors
    let result = graph
        .get(current)
        .map(|neighbors| {
            neighbors
                .iter()
                .map(|next| {
                    count_paths_memoized(
                        graph,
                        next,
                        target,
                        new_visited_dac,
                        new_visited_fft,
                        memo,
                    )
                })
                .sum()
        })
        .unwrap_or(0);

    // Cache and return
    memo.insert(cache_key, result);
    result
}

/// Wrapper to initialize memoization and count paths with required node visits
fn count_paths_with_required(graph: &Graph, start: &str, target: &str) -> i64 {
    let mut memo = MemoCache::new();
    count_paths_memoized(graph, start, target, false, false, &mut memo)
}

// --- SOLUTION ---

fn part_one(data: &Graph) -> i64 {
    count_paths(data, "you", "out")
}

fn part_two(data: &Graph) -> i64 {
    count_paths_with_required(data, "svr", "out")
}

// --- MAIN ---

fn main() {
    let input = fs::read_to_string("input.txt").expect("Failed to read input.txt");

    let data = parse_input(&input);

    println!("Part One: {}", part_one(&data));
    println!("Part Two: {}", part_two(&data));
}

// --- TESTS ---

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out";

    const EXAMPLE_PART_TWO: &str = "\
svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert_eq!(data.len(), 10); // 10 devices in example
        assert_eq!(data.get("you").unwrap(), &vec!["bbb", "ccc"]);
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 5);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE_PART_TWO);
        assert_eq!(part_two(&data), 2);
    }
}
