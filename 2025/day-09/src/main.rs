use std::fs;

fn parse_input(input: &str) -> Vec<(i64, i64)> {
    input
        .lines()
        .map(|line| {
            let parts: Vec<&str> = line.split(',').collect();
            let x = parts[0].parse().unwrap();
            let y = parts[1].parse().unwrap();
            (x, y)
        })
        .collect()
}

fn part_one(tiles: &[(i64, i64)]) -> i64 {
    let mut max_area = 0;

    // Try all pairs of red tiles as opposite corners
    for i in 0..tiles.len() {
        for j in i + 1..tiles.len() {
            let (x1, y1) = tiles[i];
            let (x2, y2) = tiles[j];

            // Calculate rectangle area
            let width = (x2 - x1).abs() + 1;
            let height = (y2 - y1).abs() + 1;
            let area = width * height;

            max_area = max_area.max(area);
        }
    }

    max_area
}

// Check if a point is on a line segment
fn point_on_segment(point: (i64, i64), p1: (i64, i64), p2: (i64, i64)) -> bool {
    let (px, py) = point;
    let (x1, y1) = p1;
    let (x2, y2) = p2;

    // Check if point is on horizontal segment
    if y1 == y2 && py == y1 {
        let min_x = x1.min(x2);
        let max_x = x1.max(x2);
        return px >= min_x && px <= max_x;
    }

    // Check if point is on vertical segment
    if x1 == x2 && px == x1 {
        let min_y = y1.min(y2);
        let max_y = y1.max(y2);
        return py >= min_y && py <= max_y;
    }

    false
}

// Ray casting algorithm to check if a point is inside or on boundary of polygon
fn point_in_polygon(point: (i64, i64), polygon: &[(i64, i64)]) -> bool {
    let (px, py) = point;
    let n = polygon.len();

    // First check if point is on any edge
    for i in 0..n {
        let p1 = polygon[i];
        let p2 = polygon[(i + 1) % n];
        if point_on_segment(point, p1, p2) {
            return true; // Point is on boundary
        }
    }

    // Then do ray casting for interior check
    let mut inside = false;
    for i in 0..n {
        let (x1, y1) = polygon[i];
        let (x2, y2) = polygon[(i + 1) % n];

        // Ray casting: cast ray to the right, count intersections
        if ((y1 > py) != (y2 > py)) && (px < (x2 - x1) * (py - y1) / (y2 - y1) + x1) {
            inside = !inside;
        }
    }

    inside
}

// Check if a horizontal or vertical edge intersects with a rectangle's interior
fn edge_intersects_rect_interior(
    edge: ((i64, i64), (i64, i64)),
    rect_min: (i64, i64),
    rect_max: (i64, i64),
) -> bool {
    let ((x1, y1), (x2, y2)) = edge;
    let (rx_min, ry_min) = rect_min;
    let (rx_max, ry_max) = rect_max;

    // Horizontal edge
    if y1 == y2 {
        let y = y1;
        // Check if edge is strictly inside rectangle (not on boundary)
        if y > ry_min && y < ry_max {
            let edge_min_x = x1.min(x2);
            let edge_max_x = x1.max(x2);
            // Check if edge overlaps with rectangle horizontally
            if edge_max_x > rx_min && edge_min_x < rx_max {
                return true;
            }
        }
    }
    // Vertical edge
    else if x1 == x2 {
        let x = x1;
        // Check if edge is strictly inside rectangle (not on boundary)
        if x > rx_min && x < rx_max {
            let edge_min_y = y1.min(y2);
            let edge_max_y = y1.max(y2);
            // Check if edge overlaps with rectangle vertically
            if edge_max_y > ry_min && edge_min_y < ry_max {
                return true;
            }
        }
    }

    false
}

// Check if polygon contains the entire rectangle
fn polygon_contains_rect(
    polygon: &[(i64, i64)],
    rect_min: (i64, i64),
    rect_max: (i64, i64),
) -> bool {
    let (x_min, y_min) = rect_min;
    let (x_max, y_max) = rect_max;

    // Check all 4 corners are inside polygon
    let corners = [
        (x_min, y_min),
        (x_min, y_max),
        (x_max, y_min),
        (x_max, y_max),
    ];

    for &corner in &corners {
        if !point_in_polygon(corner, polygon) {
            return false;
        }
    }

    // Check no polygon edges intersect rectangle interior
    for i in 0..polygon.len() {
        let edge = (polygon[i], polygon[(i + 1) % polygon.len()]);
        if edge_intersects_rect_interior(edge, rect_min, rect_max) {
            return false;
        }
    }

    true
}

fn part_two(tiles: &[(i64, i64)]) -> i64 {
    // Generate all rectangle candidates with their areas
    let mut candidates = Vec::new();
    for i in 0..tiles.len() {
        for j in i + 1..tiles.len() {
            let (x1, y1) = tiles[i];
            let (x2, y2) = tiles[j];
            let area = ((x2 - x1).abs() + 1) * ((y2 - y1).abs() + 1);
            candidates.push((area, (x1, y1), (x2, y2)));
        }
    }

    // Sort by area descending - check largest first
    candidates.sort_by(|a, b| b.0.cmp(&a.0));

    let mut max_area = 0;

    // Check rectangles in order of decreasing area
    for (area, (x1, y1), (x2, y2)) in candidates {
        // Skip if this area can't beat current max
        if area <= max_area {
            break; // All remaining will be smaller
        }

        // Define rectangle bounds
        let rect_min = (x1.min(x2), y1.min(y2));
        let rect_max = (x1.max(x2), y1.max(y2));

        // Check if polygon contains the rectangle
        if polygon_contains_rect(tiles, rect_min, rect_max) {
            max_area = area;
        }
    }

    max_area
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
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3";

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert_eq!(data.len(), 8); // 8 red tile coordinates in example
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_one(&data), 50); // Largest rectangle area from example
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        assert_eq!(part_two(&data), 24); // Largest rectangle using only red/green tiles
    }
}
