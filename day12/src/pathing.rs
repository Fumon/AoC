use std::{
    cmp::Ordering,
    collections::{BinaryHeap, HashMap},
    ops::Add,
};

use crate::heightmap::{Elevation, Point};

#[derive(Clone, Copy, PartialEq, Eq)]
enum D {
    Inf,
    Dist(u32),
}

impl Add for D {
    type Output = D;

    fn add(self, rhs: Self) -> Self::Output {
        match (self, rhs) {
            (D::Inf, D::Inf) => D::Inf,
            (D::Inf, D::Dist(_)) | (D::Dist(_), D::Inf) => D::Inf,
            (D::Dist(a), D::Dist(b)) => D::Dist(a + b),
        }
    }
}

impl Ord for D {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        match (self, other) {
            (Self::Inf, Self::Inf) => Ordering::Equal,
            (Self::Inf, _) => Ordering::Greater,
            (_, Self::Inf) => Ordering::Less,
            (Self::Dist(sel), Self::Dist(oth)) => sel.cmp(oth),
        }
    }
}

impl PartialOrd for D {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

impl Default for D {
    fn default() -> Self {
        Self::Inf
    }
}

#[derive(Clone, Copy, Eq, PartialEq)]
struct FrontNode {
    dist: D,
    position: Point,
}

impl PartialOrd for FrontNode {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for FrontNode {
    /// Minimal Ordering
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        other
            .dist
            .cmp(&self.dist)
            .then_with(|| self.position.cmp(&other.position))
    }
}

pub(crate) fn find_shortest(input: &crate::heightmap::ParsedHeightMap, rt: ReachTest) -> u32 {
    // Reachable iterator
    let reach_from_here = |h: &Point, e: Elevation| {
        return [Point(1, 0), Point(0, -1), Point(-1, 0), Point(0, 1)]
            .map(|d| h + &d)
            .into_iter()
            .filter(move |p| {
                let Some(other_ele) = input.hm.0.get(p) else { return false };

                rt(&e, other_ele)
            })
            .into_iter();
    };

    // Distances
    let mut distances: HashMap<Point, D> = HashMap::new();

    // Frontier consists of all unvisited points reachable from those already visite
    let mut front = BinaryHeap::new();

    distances.insert(input.start, D::Dist(0));
    front.push(FrontNode {
        dist: D::Dist(0),
        position: input.start,
    });

    let mut shortest_path = None;

    while let Some(FrontNode { dist, position }) = front.pop() {
        if position == input.end {
            shortest_path = Some(dist);
            break;
        }

        if dist > distances[&position] {
            continue;
        }

        let cur_elevation = input.hm.0.get(&position).unwrap().clone();

        for o_point in reach_from_here(&position, cur_elevation) {
            let o_node = FrontNode {
                dist: dist + D::Dist(1),
                position: o_point,
            };

            if &o_node.dist < distances.entry(o_node.position).or_default() {
                front.push(o_node);
                distances.insert(o_node.position, o_node.dist);
            }
        }
    }

    let Some(D::Dist(shortpath_dist)) = shortest_path else {
        panic!("Shortest path was for naught")
    };

    shortpath_dist
}

type ReachTest = fn(&Elevation, &Elevation) -> bool;

pub(crate) fn one_step_up_reach(cur: &Elevation, test: &Elevation) -> bool {
    let step_up = cur.saturating_add(1);
    if test <= &step_up {
        true
    } else {
        false
    }
}

#[cfg(test)]
mod test {
    use crate::heightmap::parse_heightmap;

    use super::{find_shortest, one_step_up_reach};

    const EXAMPLE_1: &'static str = r#"Sabqponm
abcryxxl
accszExk
acctuvwj
abdefghi"#;

    #[test]
    fn shortest() {
        let phm = parse_heightmap(EXAMPLE_1);

        assert_eq!(31, find_shortest(&phm, one_step_up_reach));
    }
}