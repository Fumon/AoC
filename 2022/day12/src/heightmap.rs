use std::{collections::HashMap, ops::Add};

pub(crate) type Elevation = u8;

// Start is elevation a
// End is elevation z
// Elevation ranges from a (lowest) to z (highest)

pub(crate) struct ParsedHeightMap {
    pub hm: Heightmap,
    pub start: Point,
    pub end: Point,
}

pub(crate) struct Heightmap(pub(crate) HashMap<Point, Elevation>);

#[derive(Debug, Clone, Copy, Hash, PartialEq, PartialOrd, Eq, Ord)]
pub(crate) struct Point(pub(crate) i32, pub(crate) i32);

impl Add for &Point {
    type Output = Point;

    fn add(self, rhs: Self) -> Self::Output {
        Point(self.0 + rhs.0, self.1 + rhs.1)
    }
}

#[derive(PartialEq, Debug)]
enum EncodedElevation {
    Start,
    End,
    Height(Elevation),
}

impl EncodedElevation {
    fn from_char(c: char) -> Self {
        match c {
            'S' => EncodedElevation::Start,
            'E' => EncodedElevation::End,
            'a'..='z' => {
                let v = (u32::from(c) & 0x1F) - 1;
                let q: u8 = v.try_into().unwrap();
                EncodedElevation::Height(q)
            }
            _ => panic!("Invalid height!"),
        }
    }

    fn to_elevation(&self) -> Elevation {
        match self {
            EncodedElevation::Start => 0,
            EncodedElevation::End => 25,
            EncodedElevation::Height(h) => *h,
        }
    }
}

pub(crate) fn parse_heightmap(i: &str) -> ParsedHeightMap {
    // Create an iterator: (Point, EncodedElevation)
    let g = i.lines().enumerate().flat_map(|(y, line): (usize, &str)| {
        line.chars()
            .enumerate()
            .map(move |(x, c)| (Point(x as i32, y as i32), EncodedElevation::from_char(c)))
    });

    let (hm, Some(start), Some(end)) = ({
        let mut hm: HashMap<Point, Elevation> = Default::default();
        let mut start: Option<Point> = None;
        let mut end: Option<Point> = None;

        for (p, e) in g.into_iter() {
            let elevation = match e {
                EncodedElevation::Start => {start = Some(p.clone()); e.to_elevation()},
                EncodedElevation::End => {end = Some(p.clone()); e.to_elevation()},
                EncodedElevation::Height(ele) => ele,
            };
            hm.insert(p, elevation);
        }
        (hm, start, end)
    }) else {
        panic!("Couldn't find the end points");
    };

    ParsedHeightMap {
        hm: Heightmap(hm),
        start,
        end,
    }
}

#[cfg(test)]
mod test {
    use crate::heightmap::Point;

    use super::{parse_heightmap, EncodedElevation};

    const EXAMPLE_1: &'static str = r#"Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi"#;
    #[test]
    fn elevation_decoding() {
        [
            (EncodedElevation::Start, 'S'),
            (EncodedElevation::End, 'E'),
            (EncodedElevation::Height(3), 'd'),
        ]
        .map(|(t, c)| (t, EncodedElevation::from_char(c)))
        .into_iter()
        .for_each(|(t, e)| assert_eq!(t, e));
        [(0, 'a'), (25, 'z'), (0, 'S'), (25, 'E')]
            .map(|(t, c)| (t, EncodedElevation::from_char(c)))
            .into_iter()
            .for_each(|(t, e)| assert_eq!(t, e.to_elevation()));
    }

    #[test]
    fn heightmap_decoding() {
        let ph = parse_heightmap(EXAMPLE_1);

        assert_eq!(Point(0, 0), ph.start);
        assert_eq!(Point(5, 2), ph.end);
        assert_eq!(5 * 8, ph.hm.0.len());
    }
}
