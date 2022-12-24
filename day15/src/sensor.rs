use nom::{
    bytes::complete::tag,
    character::complete::{digit1, line_ending},
    combinator::{map_res, opt, recognize},
    multi::many1,
    sequence::{preceded, separated_pair, terminated},
    IResult, Parser,
};

use crate::taxicab::{Coord, Distance, Point};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Sensor {
    pub position: Point,
    pub closest_beacon: Point,
    pub exclusion_radius: Distance,
}

impl Sensor {
    pub fn new(position: Point, closest_beacon: Point) -> Self {
        let exclusion_radius = position.taxicab_dist(&closest_beacon);

        Self {
            position,
            closest_beacon,
            exclusion_radius,
        }
    }

    fn parse_num(input: &str) -> IResult<&str, Coord> {
        map_res(recognize(preceded(opt(tag("-")), digit1)), |s: &str| {
            s.parse()
        })
        .parse(input)
    }
    fn parse_coords(input: &str) -> IResult<&str, Point> {
        separated_pair(
            preceded(tag("x="), Self::parse_num),
            tag(", "),
            preceded(tag("y="), Self::parse_num),
        )
        .map(|(x, y)| Point(x, y))
        .parse(input)
    }
    pub fn parse_sensor(input: &str) -> IResult<&str, Sensor> {
        preceded(
            tag("Sensor at "),
            separated_pair(
                Self::parse_coords,
                tag(": closest beacon is at "),
                Self::parse_coords,
            ),
        )
        .map(|(position, closest_beacon)| Sensor::new(position, closest_beacon))
        .parse(input)
    }

    pub fn parse_from_line(input: &str) -> Self {
        let (_, sensor) = terminated(Self::parse_sensor, opt(many1(line_ending)))
            .parse(input)
            .unwrap();

        sensor
    }
}

#[cfg(test)]
mod test {
    use crate::taxicab::Point;

    use super::Sensor;

    #[test]
    fn single_parse() {
        const LINE: &'static str = "Sensor at x=12, y=-14: closest beacon is at x=10, y=16";

        let s = Sensor::parse_from_line(LINE);
        let test = Sensor {
            position: Point(12, -14),
            closest_beacon: Point(10, 16),
            exclusion_radius: 32,
        };

        assert_eq!(test, s);
    }
}
