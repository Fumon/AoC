use std::fs::read_to_string;

use heightmap::{parse_heightmap, ParsedHeightMap};

use crate::pathing::one_step_up_reach;

fn main() {
    // Parse input
    let input = read_to_string("./input").unwrap();
    let phm = parse_heightmap(input.as_str());

    part_1(phm);
}

fn part_1(phm: ParsedHeightMap) {

    let sp = pathing::find_shortest(&phm, one_step_up_reach);

    println!("Shortest path is: {}", sp);
}

mod heightmap;

mod pathing;
