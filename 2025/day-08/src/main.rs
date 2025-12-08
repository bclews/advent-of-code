use std::collections::HashMap;
use std::fs;

/// Day 8: Playground - Minimum Spanning Tree (MST) problem
/// Connect junction boxes in 3D space using Kruskal's algorithm with Union-Find
///
/// Represents a junction box position in 3D space
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Point3D {
    x: i32,
    y: i32,
    z: i32,
}

impl Point3D {
    /// Calculate 3D Euclidean distance to another point
    fn distance_to(&self, other: &Point3D) -> f64 {
        let dx = (self.x - other.x) as f64;
        let dy = (self.y - other.y) as f64;
        let dz = (self.z - other.z) as f64;
        (dx * dx + dy * dy + dz * dz).sqrt()
    }
}

/// Union-Find (Disjoint Set Union) data structure for tracking connected components
/// Uses path compression and union by rank for efficient operations
struct UnionFind {
    parent: Vec<usize>,
    rank: Vec<usize>,
    component_count: usize, // Track number of separate components
}

impl UnionFind {
    /// Create a new Union-Find structure with `size` elements, each in its own set
    fn new(size: usize) -> Self {
        UnionFind {
            parent: (0..size).collect(), // Each node is its own parent initially
            rank: vec![0; size],
            component_count: size, // Start with N separate components
        }
    }

    /// Find the root of the set containing `x`, with path compression
    fn find(&mut self, x: usize) -> usize {
        if self.parent[x] != x {
            self.parent[x] = self.find(self.parent[x]); // Path compression
        }
        self.parent[x]
    }

    /// Unite the sets containing `x` and `y`. Returns true if they were in different sets.
    fn union(&mut self, x: usize, y: usize) -> bool {
        let root_x = self.find(x);
        let root_y = self.find(y);

        if root_x == root_y {
            return false; // Already in same component - would create cycle
        }

        // Union by rank - attach smaller tree under larger tree
        if self.rank[root_x] < self.rank[root_y] {
            self.parent[root_x] = root_y;
        } else if self.rank[root_x] > self.rank[root_y] {
            self.parent[root_y] = root_x;
        } else {
            self.parent[root_y] = root_x;
            self.rank[root_x] += 1;
        }

        self.component_count -= 1; // Two components merged into one
        true
    }

    /// Get the sizes of all connected components
    fn get_component_sizes(&mut self) -> Vec<usize> {
        let mut sizes: HashMap<usize, usize> = HashMap::new();

        for i in 0..self.parent.len() {
            let root = self.find(i);
            *sizes.entry(root).or_insert(0) += 1;
        }

        sizes.values().copied().collect()
    }

    /// Get the current number of separate components
    fn component_count(&self) -> usize {
        self.component_count
    }
}

/// Represents an edge between two junction boxes with distance
struct Edge {
    from: usize,
    to: usize,
    distance: f64,
}

/// Parse input into a vector of 3D points (junction box positions)
fn parse_input(input: &str) -> Vec<Point3D> {
    input
        .lines()
        .map(|line| {
            let parts: Vec<i32> = line.split(',').map(|s| s.trim().parse().unwrap()).collect();
            Point3D {
                x: parts[0],
                y: parts[1],
                z: parts[2],
            }
        })
        .collect()
}

/// Solve part one: process N edges and return product of 3 largest circuit sizes
fn solve(data: &[Point3D], num_connections: usize) -> i64 {
    solve_debug(data, num_connections, false)
}

/// Helper function with optional debug output for part one
fn solve_debug(data: &[Point3D], num_connections: usize, debug: bool) -> i64 {
    // Generate all possible edges between junction boxes
    let mut edges: Vec<Edge> = Vec::new();
    for i in 0..data.len() {
        for j in i + 1..data.len() {
            edges.push(Edge {
                from: i,
                to: j,
                distance: data[i].distance_to(&data[j]),
            });
        }
    }

    // Sort edges by distance (Kruskal's algorithm requires sorted edges)
    edges.sort_by(|a, b| a.distance.partial_cmp(&b.distance).unwrap());

    // Apply Kruskal's algorithm: process the first N edges
    // Note: We process N edges, not N successful unions (some may be skipped)
    let mut uf = UnionFind::new(data.len());
    let edges_to_process = num_connections.min(edges.len());

    for (i, edge) in edges.iter().enumerate().take(edges_to_process) {
        let was_union = uf.union(edge.from, edge.to);
        if debug {
            if was_union {
                println!(
                    "Edge {}: {} <-> {} (distance: {:.2}) - CONNECTED",
                    i + 1,
                    edge.from,
                    edge.to,
                    edge.distance
                );
            } else {
                println!(
                    "Edge {}: {} <-> {} (distance: {:.2}) - SKIPPED (already connected)",
                    i + 1,
                    edge.from,
                    edge.to,
                    edge.distance
                );
            }
        }
    }

    // Find the three largest circuits and multiply their sizes
    let mut sizes = uf.get_component_sizes();
    sizes.sort_unstable_by(|a, b| b.cmp(a));

    if debug {
        println!("Component sizes (sorted): {:?}", sizes);
        println!(
            "Top 3: {} * {} * {} = {}",
            sizes[0],
            sizes[1],
            sizes[2],
            sizes[0] * sizes[1] * sizes[2]
        );
    }

    (sizes[0] * sizes[1] * sizes[2]) as i64
}

/// Part One: Process 1000 edges, return product of 3 largest circuit sizes
fn part_one(data: &[Point3D]) -> i64 {
    solve(data, 1000)
}

/// Part Two: Complete the MST, return product of X coordinates of final edge
fn part_two(data: &[Point3D]) -> i64 {
    // Generate all possible edges
    let mut edges: Vec<Edge> = Vec::new();
    for i in 0..data.len() {
        for j in i + 1..data.len() {
            edges.push(Edge {
                from: i,
                to: j,
                distance: data[i].distance_to(&data[j]),
            });
        }
    }

    // Sort edges by distance
    edges.sort_by(|a, b| a.distance.partial_cmp(&b.distance).unwrap());

    // Apply Kruskal's algorithm until all nodes are in one component
    // An MST with N nodes requires exactly N-1 edges
    let mut uf = UnionFind::new(data.len());

    for edge in edges {
        if uf.union(edge.from, edge.to) {
            // Check if this union completed the MST
            if uf.component_count() == 1 {
                // This is the final edge - return product of X coordinates
                let x1 = data[edge.from].x as i64;
                let x2 = data[edge.to].x as i64;
                return x1 * x2;
            }
        }
    }

    panic!("Failed to connect all nodes into one component");
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
162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert_eq!(data.len(), 20);
        assert_eq!(
            data[0],
            Point3D {
                x: 162,
                y: 817,
                z: 812
            }
        );
        assert_eq!(
            data[19],
            Point3D {
                x: 425,
                y: 690,
                z: 689
            }
        );
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        // Example uses 10 edges: results in circuits [5, 4, 2, 2, 1, 1, 1, 1, 1, 1]
        // Answer: 5 * 4 * 2 = 40
        assert_eq!(solve(&data, 10), 40);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        // Final edge connects (216,146,977) and (117,168,530)
        // Answer: 216 * 117 = 25272
        assert_eq!(part_two(&data), 25272);
    }
}
