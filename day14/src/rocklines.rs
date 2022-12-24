use std::ops::{Add, AddAssign, Sub};

use nom::{
    bytes::complete::tag,
    character::complete::{char, digit1, line_ending},
    combinator::map_res,
    multi::separated_list1,
    sequence::separated_pair,
    IResult, Parser,
};

#[derive(Debug, PartialEq, Eq)]
pub(crate) struct PathSet(pub(crate) Vec<Path>);

impl PathSet {
    #[allow(unused)]
    pub(crate) fn find_bounds(&self) -> PathSetBounds {
        let Some((min_right, max_right, max_depth)) = self.0
            .iter()
            .flat_map(|path| path.0.iter())
            .map(|PathPoint(right, down)| (right, right, down))
            .reduce(
                |(minr, maxr, maxdepth), (_, r, d)| match (r < minr, r > maxr, d > maxdepth) {
                    (true, true, true) => (r, r, d),
                    (true, true, false) => (r, r, maxdepth),
                    (true, false, true) => (r, maxr, d),
                    (true, false, false) => (r, maxr, maxdepth),
                    (false, true, true) => (minr, r, d),
                    (false, true, false) => (minr, r, maxdepth),
                    (false, false, true) => (minr, maxr, d),
                    (false, false, false) => (minr, maxr, maxdepth),
                },
            ) else {
                panic!()
            };

        PathSetBounds {
            min_right: *min_right,
            max_right: *max_right,
            max_depth: *max_depth,
        }
    }
}

#[derive(Debug, PartialEq, Eq)]
pub(crate) struct PathSetBounds {
    pub(crate) min_right: i32,
    pub(crate) max_right: i32,
    pub(crate) max_depth: i32,
}

#[derive(Debug, PartialEq, Eq)]
pub(crate) struct Path(pub(crate) Vec<PathPoint>);

#[derive(Debug, PartialEq, Eq, Hash, Clone, Copy)]
pub(crate) struct PathPoint(pub(crate) i32, pub(crate) i32);

impl PathPoint {
    pub(crate) fn pointing_vec(self, other: Self) -> Self {
        let PathPoint(x, y) = other - self;

        match (x, y) {
            (xd, 0) if xd > 0 => PathPoint(1, 0),
            (xd, 0) if xd < 0 => PathPoint(-1, 0),
            (0, yd) if yd > 0 => PathPoint(0, 1),
            (0, yd) if yd < 0 => PathPoint(0, -1),
            _ => panic!("Unexpected or diagonal pointing vec"),
        }
    }
}

impl Sub for PathPoint {
    type Output = PathPoint;

    fn sub(self, rhs: Self) -> Self::Output {
        PathPoint(self.0 - rhs.0, self.1 - rhs.1)
    }
}

impl Add for PathPoint {
    type Output = PathPoint;

    fn add(self, rhs: Self) -> Self::Output {
        Self(self.0 + rhs.0, self.1 + rhs.1)
    }
}

impl AddAssign for PathPoint {
    fn add_assign(&mut self, rhs: Self) {
        self.0 += rhs.0;
        self.1 += rhs.1;
    }
}


pub(crate) fn parse_full_paths(input: &str) -> PathSet {
    let (_, v) = separated_list1(line_ending, parse_path).map(|v| PathSet(v)).parse(input).unwrap();
    v
}

fn parse_path(input: &str) -> IResult<&str, Path> {
    separated_list1(tag(" -> "), parse_pathpoint)
        .map(|v| Path(v))
        .parse(input)
}

fn parse_pathpoint(input: &str) -> IResult<&str, PathPoint> {
    separated_pair(
        map_res(digit1, |d: &str| d.parse()),
        char(','),
        map_res(digit1, |d: &str| d.parse()),
    )
    .map(|(l, r)| PathPoint(l, r))
    .parse(input)
}

#[cfg(test)]
mod test {
    use crate::rocklines::PathSetBounds;

    use super::{parse_full_paths, Path, PathPoint};

    const EXAMPLE: &'static str = r#"498,4 -> 498,6 -> 496,6
503,4 -> 502,4 -> 502,9 -> 494,9"#;

    #[test]
    fn path_parse() {
        let paths = parse_full_paths(EXAMPLE);

        let test = vec![
            Path(vec![
                PathPoint(498, 4),
                PathPoint(498, 6),
                PathPoint(496, 6),
            ]),
            Path(vec![
                PathPoint(503, 4),
                PathPoint(502, 4),
                PathPoint(502, 9),
                PathPoint(494, 9),
            ]),
        ];

        assert_eq!(test, paths.0);
    }

    #[test]
    fn bounds() {
        let paths = parse_full_paths(EXAMPLE);

        assert_eq!(PathSetBounds{min_right: 494, max_right: 503, max_depth: 9}, paths.find_bounds());
    }
}
