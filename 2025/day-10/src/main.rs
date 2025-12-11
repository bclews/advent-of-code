use std::fs;

// --- DATA STRUCTURES ---

/// Represents a factory machine with indicator lights and buttons
#[derive(Debug, Clone)]
struct Machine {
    num_lights: usize,
    target: Vec<bool>,
    buttons: Vec<Vec<usize>>,
    joltage: Vec<i64>,
}

/// Matrix over GF(2) for Gaussian elimination
#[derive(Debug, Clone)]
struct GF2Matrix {
    rows: usize,
    cols: usize,
    data: Vec<Vec<bool>>,
}

/// Solution to the linear system over GF(2)
#[derive(Debug, Clone)]
enum Solution {
    Inconsistent,
    Unique(Vec<bool>),
    Infinite {
        particular: Vec<bool>,
        free_vars: Vec<usize>,
        null_space: Vec<Vec<bool>>,
    },
}

// --- PARSING FUNCTIONS ---

/// Parses the indicator light diagram from a line, e.g., "[.##.]"
fn parse_diagram(line: &str) -> Option<Vec<bool>> {
    let start = line.find('[')? + 1;
    let end = line.find(']')?;
    let diagram_str = &line[start..end];
    Some(diagram_str.chars().map(|c| c == '#').collect())
}

/// Parses all button definitions from a line, e.g., "(3) (1,3) (2)"
fn parse_buttons_from_str(line: &str) -> Vec<Vec<usize>> {
    let button_strs: Vec<&str> = line
        .split('(')
        .skip(1) // Skip empty first element
        .filter_map(|s| {
            let end = s.find(')')?;
            Some(&s[..end])
        })
        .filter(|s| !s.starts_with('{') && !s.is_empty())
        .collect();

    parse_buttons_from_strs(&button_strs)
}

/// Parses the joltage requirements from a line, e.g., "{3,5,4,7}"
fn parse_joltage(line: &str) -> Vec<i64> {
    if let Some(start) = line.find('{') {
        let end = line.find('}').unwrap();
        let jolt_str = &line[start + 1..end];
        jolt_str
            .split(',')
            .map(|s| s.trim().parse::<i64>().unwrap())
            .collect()
    } else {
        vec![]
    }
}

/// Parse a single machine line from input
/// Format: [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
fn parse_machine(line: &str) -> Machine {
    let target = parse_diagram(line).expect("Failed to parse diagram");
    let num_lights = target.len();
    let buttons = parse_buttons_from_str(line);
    let joltage = parse_joltage(line);

    Machine {
        num_lights,
        target,
        buttons,
        joltage,
    }
}

/// Parse button definitions from a slice of strings like ["3", "1,3", "2"]
/// Returns Vec<Vec<usize>> where each inner vec is light indices
fn parse_buttons_from_strs(button_strs: &[&str]) -> Vec<Vec<usize>> {
    button_strs
        .iter()
        .map(|s| {
            s.split(',')
                .map(|num| num.trim().parse::<usize>().unwrap())
                .collect()
        })
        .collect()
}

fn parse_input(input: &str) -> Vec<String> {
    input.lines().map(|s| s.to_string()).collect()
}

// --- MATRIX CONSTRUCTION ---

/// Build augmented matrix A|b from a machine
/// A[i][j] = true if button j toggles light i
/// b[i] = target state of light i
fn build_augmented_matrix(machine: &Machine) -> GF2Matrix {
    let rows = machine.num_lights;
    let cols = machine.buttons.len(); // Not counting augmented column
    let mut data = vec![vec![false; cols + 1]; rows];

    // Fill matrix A
    for (button_idx, button) in machine.buttons.iter().enumerate() {
        for &light_idx in button {
            data[light_idx][button_idx] = true;
        }
    }

    // Fill augmented column b
    for (light_idx, &target_state) in machine.target.iter().enumerate() {
        data[light_idx][cols] = target_state;
    }

    GF2Matrix { rows, cols, data }
}

// --- GF(2) LINEAR ALGEBRA ---

impl GF2Matrix {
    /// Get value at (row, col), where col can be the augmented column
    fn get(&self, row: usize, col: usize) -> bool {
        self.data[row][col]
    }

    /// Set value at (row, col)
    fn set(&mut self, row: usize, col: usize, val: bool) {
        self.data[row][col] = val;
    }

    /// Find first row starting from start_row that has a 1 in column col
    fn find_pivot(&self, start_row: usize, col: usize) -> Option<usize> {
        (start_row..self.rows).find(|&row| self.get(row, col))
    }

    /// Swap two rows in the matrix
    fn swap_rows(&mut self, row1: usize, row2: usize) {
        if row1 != row2 {
            self.data.swap(row1, row2);
        }
    }

    /// XOR source_row into target_row (mod 2 addition)
    fn xor_rows(&mut self, target_row: usize, source_row: usize) {
        for col in 0..=self.cols {
            let val = self.get(target_row, col) ^ self.get(source_row, col);
            self.set(target_row, col, val);
        }
    }

    /// Perform Gaussian elimination over GF(2) to bring the matrix to RREF.
    /// Returns pivot positions for each column (None if no pivot).
    fn reduce_to_rref(&mut self) -> Vec<Option<usize>> {
        let mut pivots = vec![None; self.cols];
        let mut current_row = 0;

        for (col_idx, pivot_slot) in pivots.iter_mut().enumerate() {
            if current_row >= self.rows {
                break;
            }

            if let Some(pivot_row) = self.find_pivot(current_row, col_idx) {
                self.swap_rows(current_row, pivot_row);
                *pivot_slot = Some(current_row);

                for row in 0..self.rows {
                    if row != current_row && self.get(row, col_idx) {
                        self.xor_rows(row, current_row);
                    }
                }
                current_row += 1;
            }
        }
        pivots
    }
}

// --- SOLUTION EXTRACTION ---

/// For an infinite solution space, find the basis for the null space.
fn find_null_space_basis(
    matrix: &GF2Matrix,
    pivots: &[Option<usize>],
    free_vars: &[usize],
) -> Vec<Vec<bool>> {
    let mut null_space = Vec::new();
    for &free_col in free_vars {
        let mut basis_vec = vec![false; matrix.cols];
        basis_vec[free_col] = true;

        // Set pivot variables to satisfy homogeneous equations (Ax = 0)
        for (col, &pivot_row) in pivots.iter().enumerate() {
            if let Some(row) = pivot_row {
                basis_vec[col] = matrix.get(row, free_col);
            }
        }
        null_space.push(basis_vec);
    }
    null_space
}

/// Extract solution from RREF matrix
fn extract_solution_gf2(matrix: &GF2Matrix, pivots: &[Option<usize>]) -> Solution {
    // Check for inconsistency: a row like [0 0 ... 0 | 1]
    for row in 0..matrix.rows {
        let is_all_zero = (0..matrix.cols).all(|col| !matrix.get(row, col));
        if is_all_zero && matrix.get(row, matrix.cols) {
            return Solution::Inconsistent;
        }
    }

    let free_vars: Vec<usize> = pivots
        .iter()
        .enumerate()
        .filter_map(|(col, &p)| if p.is_none() { Some(col) } else { None })
        .collect();

    if free_vars.is_empty() {
        // Unique solution
        let mut solution = vec![false; matrix.cols];
        for (col, &pivot_row) in pivots.iter().enumerate() {
            if let Some(row) = pivot_row {
                solution[col] = matrix.get(row, matrix.cols);
            }
        }
        Solution::Unique(solution)
    } else {
        // Infinite solutions: find a particular solution and the null space basis
        let mut particular = vec![false; matrix.cols];
        for (col, &pivot_row) in pivots.iter().enumerate() {
            if let Some(row) = pivot_row {
                particular[col] = matrix.get(row, matrix.cols);
            }
        }

        let null_space = find_null_space_basis(matrix, pivots, &free_vars);

        Solution::Infinite {
            particular,
            free_vars,
            null_space,
        }
    }
}

/// Count the number of true values (Hamming weight)
fn hamming_weight(vec: &[bool]) -> usize {
    vec.iter().filter(|&&x| x).count()
}

/// Find solution with minimum Hamming weight from a solution space
fn find_minimum_solution(solution: &Solution) -> Option<Vec<bool>> {
    match solution {
        Solution::Inconsistent => None,
        Solution::Unique(sol) => Some(sol.clone()),
        Solution::Infinite {
            particular,
            free_vars,
            null_space,
        } => {
            let num_free = free_vars.len();

            // Guard against exploring a huge solution space
            if num_free > 20 {
                return Some(particular.clone());
            }

            let mut best_solution = particular.clone();
            let mut best_weight = hamming_weight(particular);

            // Iterate through all 2^n combinations of free variables
            for i in 0..(1 << num_free) {
                let mut candidate = particular.clone();

                // Create a linear combination of null space vectors based on the mask `i`
                for (j, basis_vec) in null_space.iter().enumerate() {
                    if (i >> j) & 1 == 1 {
                        // XOR this basis vector into the candidate
                        for (k, val) in candidate.iter_mut().enumerate() {
                            *val ^= basis_vec[k];
                        }
                    }
                }

                let weight = hamming_weight(&candidate);
                if weight < best_weight {
                    best_weight = weight;
                    best_solution = candidate;
                }
            }

            Some(best_solution)
        }
    }
}

// --- MACHINE SOLVING ---

/// Solve a single machine for part one: return minimum button presses needed
fn solve_machine(machine: &Machine) -> Option<usize> {
    let mut matrix = build_augmented_matrix(machine);
    let pivots = matrix.reduce_to_rref();
    let solution = extract_solution_gf2(&matrix, &pivots);

    find_minimum_solution(&solution).map(|sol| hamming_weight(&sol))
}

// --- PARTS ---

fn part_one(data: &[String]) -> i64 {
    data.iter()
        .map(|line| {
            let machine = parse_machine(line);
            solve_machine(&machine).unwrap_or(0) as i64
        })
        .sum()
}

// --- PART TWO: INTEGER LINEAR PROGRAMMING ---

const INF: f64 = f64::INFINITY;
const EPS: f64 = 1e-9;

struct SimplexSolver {
    m: usize,
    n: usize,
    tableau: Vec<Vec<f64>>,
    basic_vars: Vec<i32>,
    non_basic_vars: Vec<i32>,
}

impl SimplexSolver {
    fn new(constraints: &[Vec<f64>], obj_coeffs: &[f64]) -> Self {
        let m = constraints.len();
        let n = constraints[0].len() - 1;

        let mut non_basic_vars: Vec<i32> = (0..n as i32).collect();
        non_basic_vars.push(-1);
        let basic_vars: Vec<i32> = (n as i32..(n + m) as i32).collect();
        let mut tableau = vec![vec![0.0; n + 2]; m + 2];

        for (tableau_row, lhs_row) in tableau.iter_mut().zip(constraints.iter()) {
            tableau_row[..=n].copy_from_slice(lhs_row);
            tableau_row[n + 1] = -1.0;
        }

        for row in tableau.iter_mut().take(m) {
            row.swap(n, n + 1);
        }

        tableau[m][..n].copy_from_slice(&obj_coeffs[..n]);
        tableau[m + 1][n] = 1.0;

        Self {
            m,
            n,
            tableau,
            basic_vars,
            non_basic_vars,
        }
    }

    fn pivot(&mut self, r: usize, s: usize) {
        let k = 1.0 / self.tableau[r][s];
        for i in 0..self.m + 2 {
            if i == r {
                continue;
            }
            for j in 0..self.n + 2 {
                if j != s {
                    self.tableau[i][j] -= self.tableau[r][j] * self.tableau[i][s] * k;
                }
            }
        }
        for val in self.tableau[r].iter_mut() {
            *val *= k;
        }
        for row in self.tableau.iter_mut() {
            row[s] *= -k;
        }
        self.tableau[r][s] = k;
        std::mem::swap(&mut self.basic_vars[r], &mut self.non_basic_vars[s]);
    }

    fn iterate(&mut self, p_idx: usize) -> bool {
        loop {
            let mut best_s = usize::MAX;
            let mut best_val = (INF, i32::MAX);

            for i in 0..=self.n {
                if p_idx != 0 || self.non_basic_vars[i] != -1 {
                    let val = self.tableau[self.m + p_idx][i];
                    let key = (val, self.non_basic_vars[i]);
                    if best_s == usize::MAX
                        || key.0 < best_val.0 - EPS
                        || ((key.0 - best_val.0).abs() <= EPS && key.1 < best_val.1)
                    {
                        best_s = i;
                        best_val = key;
                    }
                }
            }
            let s = best_s;

            if self.tableau[self.m + p_idx][s] > -EPS {
                return true;
            }

            let mut best_r = usize::MAX;
            let mut best_r_key = (INF, i32::MAX);

            for i in 0..self.m {
                if self.tableau[i][s] > EPS {
                    let ratio = self.tableau[i][self.n + 1] / self.tableau[i][s];
                    let key = (ratio, self.basic_vars[i]);
                    if best_r == usize::MAX
                        || key.0 < best_r_key.0 - EPS
                        || ((key.0 - best_r_key.0).abs() <= EPS && key.1 < best_r_key.1)
                    {
                        best_r = i;
                        best_r_key = key;
                    }
                }
            }
            let r = best_r;

            if r == usize::MAX {
                return false;
            }
            self.pivot(r, s);
        }
    }

    fn prepare_initial_solution(&mut self) -> bool {
        let mut split_r = 0;
        let mut min_val = self.tableau[0][self.n + 1];
        for (i, row) in self.tableau.iter().enumerate().take(self.m).skip(1) {
            if row[self.n + 1] < min_val {
                min_val = row[self.n + 1];
                split_r = i;
            }
        }

        if self.tableau[split_r][self.n + 1] < -EPS {
            self.pivot(split_r, self.n);
            if !self.iterate(1) || self.tableau[self.m + 1][self.n + 1] < -EPS {
                return false; // Infeasible
            }
            for i in 0..self.m {
                if self.basic_vars[i] == -1 {
                    let mut best_s = 0;
                    let mut best_key = (self.tableau[i][0], self.non_basic_vars[0]);
                    for j in 1..self.n {
                        let key = (self.tableau[i][j], self.non_basic_vars[j]);
                        if key.0 < best_key.0 - EPS
                            || ((key.0 - best_key.0).abs() <= EPS && key.1 < best_key.1)
                        {
                            best_s = j;
                            best_key = key;
                        }
                    }
                    self.pivot(i, best_s);
                }
            }
        }
        true
    }

    fn get_solution(self, obj_coeffs: &[f64]) -> (f64, Option<Vec<f64>>) {
        let mut x = vec![0.0; self.n];
        for i in 0..self.m {
            if self.basic_vars[i] >= 0 && (self.basic_vars[i] as usize) < self.n {
                x[self.basic_vars[i] as usize] = self.tableau[i][self.n + 1];
            }
        }
        let mut sum_val = 0.0;
        for i in 0..self.n {
            sum_val += obj_coeffs[i] * x[i];
        }
        (sum_val, Some(x))
    }
}

fn simplex(constraints: &[Vec<f64>], obj_coeffs: &[f64]) -> (f64, Option<Vec<f64>>) {
    let mut solver = SimplexSolver::new(constraints, obj_coeffs);

    if !solver.prepare_initial_solution() {
        return (-INF, None); // Infeasible
    }

    if solver.iterate(0) {
        solver.get_solution(obj_coeffs)
    } else {
        (-INF, None) // Unbounded
    }
}

struct BranchAndBoundSolver<'a> {
    stack: Vec<Vec<Vec<f64>>>,
    obj_coeffs: &'a [f64],
    best_val: f64,
}

impl<'a> BranchAndBoundSolver<'a> {
    fn new(initial_problem: Vec<Vec<f64>>, obj_coeffs: &'a [f64]) -> Self {
        Self {
            stack: vec![initial_problem],
            obj_coeffs,
            best_val: INF,
        }
    }

    fn solve(&mut self) -> i64 {
        while let Some(current_problem_constraints) = self.stack.pop() {
            self.process_node(current_problem_constraints);
        }

        if self.best_val == INF {
            0
        } else {
            self.best_val.round() as i64
        }
    }

    fn process_node(&mut self, current_problem_constraints: Vec<Vec<f64>>) {
        let (val, x_opt) = simplex(&current_problem_constraints, self.obj_coeffs);

        // Bounding step: if this node can't produce a better solution, prune it.
        if val == -INF || val >= self.best_val - EPS {
            return;
        }

        if let Some(x) = x_opt {
            // Check if the solution is all integers
            if let Some((idx, val)) = Self::find_fractional_variable(&x) {
                // Branching step
                self.branch(current_problem_constraints, idx, val);
            } else {
                // Update best solution found so far
                if val < self.best_val {
                    self.best_val = val;
                }
            }
        }
    }

    fn find_fractional_variable(solution: &[f64]) -> Option<(usize, f64)> {
        solution
            .iter()
            .enumerate()
            .find(|(_, &xv)| (xv - xv.round()).abs() > EPS)
            .map(|(i, &xv)| (i, xv))
    }

    fn branch(
        &mut self,
        current_problem_constraints: Vec<Vec<f64>>,
        fractional_idx: usize,
        fractional_val: f64,
    ) {
        let n_cols = current_problem_constraints[0].len();

        // Branch 1: Add constraint x_i <= floor(v)
        let floor_v = fractional_val.floor();
        let mut row1 = vec![0.0; n_cols];
        row1[fractional_idx] = 1.0;
        row1[n_cols - 1] = floor_v;
        let mut branch1_problem = current_problem_constraints.clone();
        branch1_problem.push(row1);
        self.stack.push(branch1_problem);

        // Branch 2: Add constraint x_i >= ceil(v)  (i.e., -x_i <= -ceil(v))
        let ceil_v = fractional_val.ceil();
        let mut row2 = vec![0.0; n_cols];
        row2[fractional_idx] = -1.0;
        row2[n_cols - 1] = -ceil_v;
        let mut branch2_problem = current_problem_constraints; // Reuse allocation
        branch2_problem.push(row2);
        self.stack.push(branch2_problem);
    }
}

fn solve_ilp_bnb(initial_problem: Vec<Vec<f64>>, obj_coeffs: &[f64]) -> i64 {
    let mut solver = BranchAndBoundSolver::new(initial_problem, obj_coeffs);
    solver.solve()
}

fn solve_machine_joltage(machine: &Machine) -> Option<i64> {
    let n_counters = machine.joltage.len();
    let n_buttons = machine.buttons.len();

    if n_counters == 0 {
        return Some(0);
    }
    if n_buttons == 0 {
        return if machine.joltage.iter().all(|&j| j == 0) {
            Some(0)
        } else {
            None
        };
    }

    // Set up the problem for the ILP solver
    // The problem is to minimize sum(x) subject to Ax = t and x >= 0.
    // This is converted to a set of inequalities for the simplex solver:
    // 1. Ax <= t
    // 2. -Ax <= -t
    // 3. -x_i <= 0 for all i (non-negativity)

    let num_constraints = 2 * n_counters + n_buttons;
    let num_vars = n_buttons;
    let mut matrix = vec![vec![0.0; num_vars + 1]; num_constraints];

    // Build the non-negativity constraints first: -x_i <= 0
    for (i, row) in matrix.iter_mut().take(n_buttons).enumerate() {
        row[i] = -1.0;
        // The right-hand side is 0, which is default
    }

    // Build the equality constraints (as two inequalities)
    let mut a = vec![vec![0i64; n_buttons]; n_counters];
    for (button_idx, button) in machine.buttons.iter().enumerate() {
        for &counter_idx in button {
            if counter_idx < n_counters {
                a[counter_idx][button_idx] = 1;
            }
        }
    }

    // Ax <= t
    for i in 0..n_counters {
        for j in 0..n_buttons {
            matrix[n_buttons + i][j] = a[i][j] as f64;
        }
        matrix[n_buttons + i][num_vars] = machine.joltage[i] as f64;
    }

    // -Ax <= -t
    for i in 0..n_counters {
        for j in 0..n_buttons {
            matrix[n_buttons + n_counters + i][j] = -a[i][j] as f64;
        }
        matrix[n_buttons + n_counters + i][num_vars] = -machine.joltage[i] as f64;
    }

    // Objective function: minimize sum(x), so all coefficients are 1.0
    let obj_coeffs = vec![1.0; n_buttons];

    Some(solve_ilp_bnb(matrix, &obj_coeffs))
}

fn part_two(data: &[String]) -> i64 {
    data.iter()
        .map(|line| {
            let machine = parse_machine(line);
            solve_machine_joltage(&machine).unwrap_or(0)
        })
        .sum()
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
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}";

    #[test]
    fn test_parse_diagram_from_line() {
        assert_eq!(
            parse_diagram("[.##.]").unwrap(),
            vec![false, true, true, false]
        );
        assert_eq!(
            parse_diagram("[...#.]").unwrap(),
            vec![false, false, false, true, false]
        );
        assert_eq!(
            parse_diagram("[.###.#]").unwrap(),
            vec![false, true, true, true, false, true]
        );
    }

    #[test]
    fn test_parse_buttons_from_strs_fn() {
        let buttons = parse_buttons_from_strs(&["3", "1,3", "0,2"]);
        assert_eq!(buttons, vec![vec![3], vec![1, 3], vec![0, 2]]);
    }

    #[test]
    fn test_parse_machine() {
        let line = "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}";
        let machine = parse_machine(line);
        assert_eq!(machine.num_lights, 4);
        assert_eq!(machine.target, vec![false, true, true, false]);
        assert_eq!(machine.buttons.len(), 6);
        assert_eq!(machine.buttons[0], vec![3]);
        assert_eq!(machine.buttons[1], vec![1, 3]);
        assert_eq!(machine.joltage, vec![3, 5, 4, 7]);
    }

    #[test]
    fn test_parse_input() {
        let data = parse_input(EXAMPLE);
        assert_eq!(data.len(), 3);
    }

    #[test]
    fn test_xor_rows() {
        let mut matrix = GF2Matrix {
            rows: 2,
            cols: 2,
            data: vec![vec![true, false, true], vec![true, true, false]],
        };
        matrix.xor_rows(0, 1);
        assert_eq!(matrix.data[0], vec![false, true, true]);
        assert_eq!(matrix.data[1], vec![true, true, false]);
    }

    #[test]
    fn test_hamming_weight() {
        assert_eq!(hamming_weight(&[true, false, true, false]), 2);
        assert_eq!(hamming_weight(&[false, false, false]), 0);
        assert_eq!(hamming_weight(&[true, true, true]), 3);
    }

    #[test]
    fn test_solve_machine_example1() {
        let line = "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}";
        let machine = parse_machine(line);
        let result = solve_machine(&machine);
        assert_eq!(result, Some(2));
    }

    #[test]
    fn test_solve_machine_example2() {
        let line = "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}";
        let machine = parse_machine(line);
        let result = solve_machine(&machine);
        assert_eq!(result, Some(3));
    }

    #[test]
    fn test_solve_machine_example3() {
        let line = "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}";
        let machine = parse_machine(line);
        let result = solve_machine(&machine);
        assert_eq!(result, Some(2));
    }

    #[test]
    fn test_part_one() {
        let data = parse_input(EXAMPLE);
        // Machine 1: 2 presses, Machine 2: 3 presses, Machine 3: 2 presses
        // Total: 2 + 3 + 2 = 7
        assert_eq!(part_one(&data), 7);
    }

    #[test]
    fn test_part_two() {
        let data = parse_input(EXAMPLE);
        // Machine 1: 10 presses, Machine 2: 12 presses, Machine 3: 11 presses
        // Total: 10 + 12 + 11 = 33
        assert_eq!(part_two(&data), 33);
    }

    #[test]
    fn test_solve_machine_joltage_example1() {
        let line = "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}";
        let machine = parse_machine(line);
        // Joltage: {3,5,4,7} requires minimum 10 presses
        // One way: (3)×1, (1,3)×3, (2,3)×3, (0,2)×1, (0,1)×2
        assert_eq!(solve_machine_joltage(&machine), Some(10));
    }

    #[test]
    fn test_solve_machine_joltage_example2() {
        let line = "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}";
        let machine = parse_machine(line);
        // Joltage: {7,5,12,7,2} requires minimum 12 presses
        // One way: (0,2,3,4)×2, (2,3)×5, (0,1,2)×5
        assert_eq!(solve_machine_joltage(&machine), Some(12));
    }

    #[test]
    fn test_solve_machine_joltage_example3() {
        let line = "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}";
        let machine = parse_machine(line);
        // Joltage: {10,11,11,5,10,5} requires minimum 11 presses
        // One way: (0,1,2,3,4)×5, (0,1,2,4,5)×5, (1,2)×1
        assert_eq!(solve_machine_joltage(&machine), Some(11));
    }
}

