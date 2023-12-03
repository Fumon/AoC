#![feature(array_windows)]
use std::{
    collections::HashSet,
    fs::read_to_string,
    iter::{from_fn, repeat, successors, zip},
};

use rocklines::{parse_full_paths, PathPoint, PathSet};

use crate::rocklines::PathSetBounds;

fn main() {
    let paths = parse_full_paths(read_to_string("./input").unwrap().as_str());

    let sand_units = part_1(&paths);
    println!("Total Sand Settled: {}", sand_units);

    let sand_units2 = part_2(&paths);
    println!("Total Sand Settled including floor: {}", sand_units2);
}

fn fill_cave(paths: &PathSet) -> HashSet<PathPoint> {
    let mut cave = HashSet::new();

    paths
        .0
        .iter()
        .flat_map(|path| {
            path.0.array_windows().flat_map(|[a, b]| {
                let pointing = a.pointing_vec(*b);

                successors(Some(a.clone()), move |p| {
                    if p == b {
                        None
                    } else {
                        Some(*p + pointing)
                    }
                })
            })
        })
        .for_each(|pt| {
            cave.insert(pt);
        });

    cave
}

fn part_1(paths: &PathSet) -> usize {
    let mut cave = fill_cave(paths);
    
    let PathSetBounds {
        min_right: _,
        max_right: _,
        max_depth,
    } = paths.find_bounds();

    let mut sandcount = 0;
    let mut sp = SandPath::new();

    let mut s = PathPoint(500, 0);

    loop {
        // One hop this time
        {
            let np = s + PathPoint(0, 1);
            if np.1 > max_depth {
                break;
            }
            if !cave.contains(&np) {
                sp.push(s);
                s = np;
                continue;
            }
        }

        // Slide to the left
        {
            let sl = s + PathPoint(-1, 1);
            if !cave.contains(&sl) {
                sp.push(s);
                s = sl;
                continue;
            }
        }

        // Slide to the right
        {
            let sr = s + PathPoint(1, 1);
            if !cave.contains(&sr) {
                sp.push(s);
                s = sr;
                continue;
            }
        }

        // CRISSCROSS (settle)
        cave.insert(s);
        sandcount += 1;
        match sp.pop() {
            Some(prev) => s = prev,
            None => panic!("Not supposed to run out of points"),
        }
    }

    sandcount
}

const DOWN: PathPoint = PathPoint(0, 1);
const LEFTDOWN: PathPoint = PathPoint(-1, 1);
const RIGHTDOWN: PathPoint = PathPoint(1, 1);

const DESENT_SEQUENCE: [PathPoint; 3] = [DOWN, LEFTDOWN, RIGHTDOWN];

fn part_2(paths: &PathSet) -> usize {
    let mut cave = fill_cave(paths);
    
    let PathSetBounds {
        min_right: _,
        max_right: _,
        max_depth,
    } = paths.find_bounds();

    let floor_depth = max_depth + 2;

    let mut sandcount = 0;
    let mut sp = SandPath::new();

    let mut s = PathPoint(500, 0);

    'outer: loop {
        for a in DESENT_SEQUENCE {
            let sd = s + a;
            if sd.1 < floor_depth && !cave.contains(&sd) {
                sp.push(s);
                s = sd;
                continue 'outer;
            }
        }

        // CRISSCROSS (settle)
        cave.insert(s);
        sandcount += 1;
        match sp.pop() {
            Some(prev) => s = prev,
            None => break,
        }
    }

    sandcount
}

type SandPath = Vec<PathPoint>;

#[cfg(test)]
mod test {
    use crate::rocklines::{parse_full_paths, PathSet};

    const EXAMPLE: &'static str = r#"498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9"#;

    fn parse_example() -> PathSet {
        parse_full_paths(EXAMPLE)
    }

    #[test]
    fn part_1() {
        assert_eq!(24, super::part_1(&parse_example()));
    }

    #[test]
    fn part_2() {
        assert_eq!(93, super::part_2(&parse_example()));
    }
}

mod rocklines;
