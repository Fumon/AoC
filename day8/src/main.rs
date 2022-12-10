#![feature(is_some_and)]
#![warn(clippy::shadow_same)]

use std::{
    collections::{HashMap, HashSet, VecDeque},
    fs::File,
    hash::Hash,
    io::{BufRead, BufReader},
    iter::{once, repeat},
};

fn main() {
    let liter = BufReader::new(File::open("./input").unwrap())
        .lines()
        .flatten()
        .peekable();

    // Allocate the horizontal queues
    // 10 levels
    // [VecDeque<&Tree>>; 10]

    let mut trees: HashMap<Point, Tree> = Default::default();

    const SIDE_LENGTH: usize = 99;
    let [Ok(mut row_queues), Ok(mut column_queues)]: [Result<[[VecDeque<Point>; 10]; SIDE_LENGTH], _>; 2] = [(), ()].map(|_| {
        (0..SIDE_LENGTH)
            .map(|_| { let r: [VecDeque<Point>; 10] = Default::default(); r} )
            .collect::<Vec<[VecDeque<Point>; 10]>>()
            .try_into()
    }) else {
        panic!("Couldn't construct queues");
    };

    liter.enumerate().for_each(|(y, line)| {
        line.chars()
            .enumerate()
            .map(|(h, byte)| (h, byte.to_string().parse::<u8>().unwrap()))
            .for_each(|(x, height)| match height {
                0..=9 => {
                    let pos = Point(x, y);
                    trees.insert(
                        pos,
                        Tree {
                            height,
                            position: Point(x, y),
                        },
                    );

                    let hu = height as usize;

                    column_queues[x][hu].push_back(pos);
                    row_queues[y][hu].push_back(pos);
                }
                _ => panic!("Invalid height"),
            })
    });

    // Now "ray cast" by pulling off each queue from each direction and deduplicating with a hashset
    let mut seen: HashSet<Point> = Default::default();

    let xaccessor = Point::x as Pointaccessor;
    let yaccessor = Point::y as Pointaccessor;

    let lt: fn(&usize, &usize) -> bool = usize::lt;
    let gt: fn(&usize, &usize) -> bool = usize::gt;

    let hqueues_iter = row_queues.iter().zip(repeat(xaccessor));
    let vqueues_iter = column_queues.iter().zip(repeat(yaccessor));

    for (view, coord) in hqueues_iter.chain(vqueues_iter) {
        // Iterate backwards from highest to lowest trees
        let mut last_visible_left_coord: Option<usize> = None;
        let mut last_visible_right_coord: Option<usize> = None;

        view.iter().rev().for_each(|q| {
            once((q.front(), &mut last_visible_left_coord, lt))
                .chain(once((q.back(), &mut last_visible_right_coord, gt)))
                .for_each(|(point, tracker, compareop)| {
                    if let Some(new_seen) = {
                        if let Some(p) = point {
                            let t = coord(p);
                            if tracker.is_none() || tracker.is_some_and(|old| compareop(&t, &old)) {
                                *tracker = Some(t);
                                point
                            } else {
                                None
                            }
                        } else {
                            None
                        }
                    } {
                        seen.insert(new_seen.to_owned());
                    }
                });
        });
    }

    println!("# of visible trees: {}", seen.len());
}

#[derive(Debug, Hash, PartialEq, Eq)]
struct Tree {
    height: u8,
    position: Point,
}

/// Tuple is (x, y): (horizontal, vertical)
#[derive(Debug, Hash, PartialEq, Eq, Clone, Copy)]
struct Point(usize, usize);

type Pointaccessor = fn(&Point) -> usize;
impl Point {
    fn x(&self) -> usize {
        self.0
    }
    fn y(&self) -> usize {
        self.1
    }
}
