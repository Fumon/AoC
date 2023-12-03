pub mod taxicab;

pub mod sensor;

pub mod lines {
    use crate::taxicab::Point;

    /// Line describes the slope sign and y intercept of a line with slope 1 or -1
    #[derive(Debug, PartialEq, Hash, Eq, PartialOrd, Ord)]
    pub struct Line(pub Sign, pub i32);

    impl Line {
        pub fn intersection_point(&self, oth: &Self) -> Result<Point, &'static str> {
            let (p, n) = match (&self.0, &oth.0) {
                (Sign::Neg, Sign::Neg) | (Sign::Pos, Sign::Pos) => {
                    return Err("Parallel lines aren't guaranteed to intersect")
                }
                (Sign::Neg, Sign::Pos) => (oth, self),
                (Sign::Pos, Sign::Neg) => (self, oth),
            };

            let xint = match n.1 - p.1 {
                x if x % 2 != 0 => return Err("Non-integer x intersect"),
                x => x / 2,
            };

            let yint = xint + p.1;

            Ok(Point(xint, yint))
        }
    }

    #[derive(Debug, PartialEq, Hash, Eq, PartialOrd, Ord)]
    pub enum Sign {
        Neg,
        Pos,
    }
}
