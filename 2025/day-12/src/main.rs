use std::collections::HashSet;
use std::fs;

type Coord = (i32, i32);
type Shape = Vec<Coord>;

#[derive(Debug, Clone)]
struct Present {
    cells: Shape,
}

#[derive(Debug)]
struct Region {
    width: usize,
    height: usize,
    required: Vec<usize>,
}

struct Grid {
    width: usize,
    height: usize,
    cells: Vec<Vec<bool>>,
}

#[derive(Debug, Clone)]
struct Transformation {
    cells: Shape,
}

impl Grid {
    fn new(width: usize, height: usize) -> Self {
        Grid {
            width,
            height,
            cells: vec![vec![false; width]; height],
        }
    }

    fn can_place(&self, shape: &[Coord], x: usize, y: usize) -> bool {
        for (dx, dy) in shape {
            let nx = x as i32 + dx;
            let ny = y as i32 + dy;

            if nx < 0 || ny < 0 || nx >= self.width as i32 || ny >= self.height as i32 {
                return false;
            }

            if self.cells[ny as usize][nx as usize] {
                return false;
            }
        }
        true
    }

    fn place(&mut self, shape: &[Coord], x: usize, y: usize) {
        for (dx, dy) in shape {
            let nx = (x as i32 + dx) as usize;
            let ny = (y as i32 + dy) as usize;
            self.cells[ny][nx] = true;
        }
    }

    fn remove(&mut self, shape: &[Coord], x: usize, y: usize) {
        for (dx, dy) in shape {
            let nx = (x as i32 + dx) as usize;
            let ny = (y as i32 + dy) as usize;
            self.cells[ny][nx] = false;
        }
    }
}

fn normalize_shape(cells: &[Coord]) -> Shape {
    if cells.is_empty() {
        return Vec::new();
    }

    let min_x = cells.iter().map(|(x, _)| x).min().unwrap();
    let min_y = cells.iter().map(|(_, y)| y).min().unwrap();

    let mut normalized: Vec<Coord> = cells.iter().map(|(x, y)| (x - min_x, y - min_y)).collect();

    normalized.sort();
    normalized
}

fn rotate_90(cells: &[Coord]) -> Shape {
    let rotated: Vec<Coord> = cells.iter().map(|(x, y)| (-y, *x)).collect();
    normalize_shape(&rotated)
}

fn flip_horizontal(cells: &[Coord]) -> Shape {
    let flipped: Vec<Coord> = cells.iter().map(|(x, y)| (-x, *y)).collect();
    normalize_shape(&flipped)
}

fn add_unique_transformation(
    shape: Shape,
    unique: &mut HashSet<Shape>,
    transformations: &mut Vec<Transformation>,
) {
    if unique.insert(shape.clone()) {
        transformations.push(Transformation { cells: shape });
    }
}

fn generate_transformations(present: &Present) -> Vec<Transformation> {
    let mut unique = HashSet::new();
    let mut transformations = Vec::new();
    let mut current = present.cells.clone();

    for _ in 0..4 {
        add_unique_transformation(current.clone(), &mut unique, &mut transformations);

        let flipped = flip_horizontal(&current);
        add_unique_transformation(flipped, &mut unique, &mut transformations);

        current = rotate_90(&current);
    }

    transformations
}

fn try_place_transformation(
    grid: &mut Grid,
    transformation: &Transformation,
    presents: &[usize],
    present_idx: usize,
    transformations: &[Vec<Transformation>],
) -> bool {
    for y in 0..grid.height {
        for x in 0..grid.width {
            if grid.can_place(&transformation.cells, x, y) {
                grid.place(&transformation.cells, x, y);

                if backtrack(grid, presents, present_idx + 1, transformations) {
                    return true;
                }

                grid.remove(&transformation.cells, x, y);
            }
        }
    }
    false
}

fn backtrack(
    grid: &mut Grid,
    presents: &[usize],
    present_idx: usize,
    transformations: &[Vec<Transformation>],
) -> bool {
    if present_idx >= presents.len() {
        return true;
    }

    let shape_idx = presents[present_idx];

    for transformation in &transformations[shape_idx] {
        if try_place_transformation(grid, transformation, presents, present_idx, transformations) {
            return true;
        }
    }

    false
}

fn expand_required_presents(required: &[usize]) -> Vec<usize> {
    let mut presents = Vec::new();
    for (shape_idx, &count) in required.iter().enumerate() {
        for _ in 0..count {
            presents.push(shape_idx);
        }
    }
    presents
}

fn calculate_total_cells(presents: &[usize], shapes: &[Present]) -> usize {
    presents.iter().map(|&idx| shapes[idx].cells.len()).sum()
}

fn solve_region(
    shapes: &[Present],
    region: &Region,
    transformations: &[Vec<Transformation>],
) -> bool {
    let presents_to_place = expand_required_presents(&region.required);
    let total_cells = calculate_total_cells(&presents_to_place, shapes);

    if total_cells > region.width * region.height {
        return false;
    }

    let mut grid = Grid::new(region.width, region.height);
    backtrack(&mut grid, &presents_to_place, 0, transformations)
}

fn parse_shape_grid(lines: &[&str]) -> Shape {
    let mut cells = Vec::new();
    for (y, line) in lines.iter().enumerate() {
        for (x, ch) in line.chars().enumerate() {
            if ch == '#' {
                cells.push((x as i32, y as i32));
            }
        }
    }
    normalize_shape(&cells)
}

fn parse_shapes(sections: &[&str]) -> (Vec<Present>, usize) {
    let mut shapes = Vec::new();
    let mut i = 0;

    while i < sections.len() {
        let section = sections[i].trim();
        if section.is_empty() {
            i += 1;
            continue;
        }

        if section.contains(':') && !section.contains('x') {
            let lines: Vec<&str> = section.lines().collect();
            let cells = parse_shape_grid(&lines[1..]);
            shapes.push(Present { cells });
            i += 1;
        } else {
            break;
        }
    }

    (shapes, i)
}

fn parse_region_line(line: &str) -> Option<Region> {
    if !line.contains(':') || !line.contains('x') {
        return None;
    }

    let parts: Vec<&str> = line.split(':').collect();
    let dims: Vec<&str> = parts[0].split('x').collect();
    let width = dims[0].trim().parse().ok()?;
    let height = dims[1].trim().parse().ok()?;

    let required: Vec<usize> = parts[1]
        .split_whitespace()
        .filter_map(|s| s.parse().ok())
        .collect();

    Some(Region {
        width,
        height,
        required,
    })
}

fn parse_regions(sections: &[&str]) -> Vec<Region> {
    sections
        .iter()
        .flat_map(|section| section.lines())
        .filter_map(|line| parse_region_line(line.trim()))
        .collect()
}

fn parse_input(input: &str) -> (Vec<Present>, Vec<Region>) {
    let sections: Vec<&str> = input.split("\n\n").collect();
    let (shapes, shape_end_idx) = parse_shapes(&sections);
    let regions = parse_regions(&sections[shape_end_idx..]);
    (shapes, regions)
}

fn part_one(shapes: &[Present], regions: &[Region]) -> usize {
    let transformations: Vec<Vec<Transformation>> =
        shapes.iter().map(generate_transformations).collect();

    regions
        .iter()
        .filter(|region| solve_region(shapes, region, &transformations))
        .count()
}

fn main() {
    let input = fs::read_to_string("input.txt").expect("Failed to read input.txt");

    let (shapes, regions) = parse_input(&input);

    println!("Part One: {}", part_one(&shapes, &regions));
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = "\
0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2";

    #[test]
    fn test_parse_input() {
        let (shapes, regions) = parse_input(EXAMPLE);
        assert_eq!(shapes.len(), 6);
        assert_eq!(regions.len(), 3);
        assert_eq!(shapes[4].cells.len(), 7);
    }

    #[test]
    fn test_normalize_shape() {
        let cells = vec![(5, 3), (6, 3), (5, 4)];
        let normalized = normalize_shape(&cells);
        assert_eq!(normalized[0], (0, 0));
        assert!(normalized.contains(&(0, 1)));
        assert!(normalized.contains(&(1, 0)));
    }

    #[test]
    fn test_rotation() {
        let cells = vec![(0, 0), (1, 0), (0, 1)];
        let rotated = rotate_90(&cells);
        assert!(rotated.contains(&(0, 0)));
    }

    #[test]
    fn test_can_place() {
        let mut grid = Grid::new(4, 4);
        let shape = vec![(0, 0), (1, 0), (2, 0)];

        assert!(grid.can_place(&shape, 0, 0));
        assert!(grid.can_place(&shape, 1, 0));
        assert!(!grid.can_place(&shape, 2, 0));

        grid.place(&shape, 0, 0);
        assert!(!grid.can_place(&shape, 0, 0));
    }

    #[test]
    fn test_part_one() {
        let (shapes, regions) = parse_input(EXAMPLE);
        assert_eq!(part_one(&shapes, &regions), 2);
    }
}
