use std::{collections::HashSet, fs::read_to_string, iter::from_fn};

use day15::{
    sensor::Sensor,
    taxicab::{Coord, Point},
};
use nom::{character::complete::line_ending, multi::separated_list1, Parser};

fn main() {
    let sensors = parse_all_sensors(read_to_string("./input").unwrap().as_str());

    const P1_Y_TARGET: Coord = 2_000_000;
    println!(
        "In the row y={}, {} positions cannot contain beacons",
        P1_Y_TARGET,
        part_1(sensors, P1_Y_TARGET)
    )
}

fn parse_all_sensors(input: &str) -> Vec<Sensor> {
    let (_, v) = separated_list1(line_ending, Sensor::parse_sensor)
        .parse(input)
        .unwrap();
    v
}

fn part_1(sensors: Vec<Sensor>, y_target: Coord) -> usize {
    let mut beacon_coords = HashSet::new();
    let mut sensor_iter = sensors.iter();

    from_fn(|| loop {
        let Some(s) = sensor_iter.next() else {
                return None;
            };

        if s.closest_beacon.1 == y_target {
            beacon_coords.insert(s.closest_beacon.0);
        }

        let radius = s.exclusion_radius as i32
            - s.position.taxicab_dist(&Point(s.position.0, y_target)) as i32;

        if radius < 0 {
            continue;
        }

        return Some((s.position.0 - radius)..=(s.position.0 + radius));
    })
    .flatten()
    .collect::<HashSet<Coord>>()
    .difference(&beacon_coords)
    .count()
}

#[cfg(test)]
mod test {
    use day15::sensor::Sensor;

    const EXAMPLE: &'static str = r#"Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3"#;

    fn parse_example() -> Vec<Sensor> {
        super::parse_all_sensors(EXAMPLE)
    }

    #[test]
    fn parse_count_matches() {
        assert_eq!(14, parse_example().len());
    }

    #[test]
    fn part_1() {
        assert_eq!(26, super::part_1(parse_example(), 10));
    }
}
