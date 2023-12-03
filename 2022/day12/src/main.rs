use std::fs::read_to_string;

use heightmap::{parse_heightmap, ParsedHeightMap};

use crate::pathing::{one_step_up_reach, EndGoal, one_step_down_reach};

fn main() {
    // Parse input
    let input = read_to_string("./input").unwrap();
    let phm = parse_heightmap(input.as_str());

    part_1(&phm);
    part_2(&phm);
}

fn part_1(phm: &ParsedHeightMap) {

    let sp = pathing::find_shortest(&phm.hm, &phm.start, EndGoal::Point(phm.end), one_step_up_reach);

    println!("Shortest start to end path is: {}", sp);
}

fn part_2(phm: &ParsedHeightMap) {

    let sp = pathing::find_shortest(&phm.hm, &phm.end, EndGoal::Height(0), one_step_down_reach);

    println!("Shortest zero to endpoint path is: {}", sp);
}

mod heightmap;

mod pathing;
