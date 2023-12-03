use std::ops::{Add, Sub};

pub type Coord = i32;
pub type Distance = u32;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Point(pub Coord, pub Coord);

impl Sub for Point {
    type Output = Self;

    fn sub(self, rhs: Self) -> Self::Output {
        Self(self.0 - rhs.1, self.1 - rhs.1)
    }
}

impl Add for Point {
    type Output = Point;

    fn add(self, rhs: Self) -> Self::Output {
        Self(self.0 + rhs.0, self.1 + rhs.1)
    }
}

impl Point {
    pub fn taxicab_dist(&self, other: &Self) -> Distance {
        self.0.abs_diff(other.0) + self.1.abs_diff(other.1)
    }
}
