#![allow(unused)]
#![feature(iter_array_chunks)]

use std::{
    borrow::BorrowMut,
    collections::{hash_set::Intersection, HashSet},
    fs::File,
    io::{self, BufRead, BufReader},
    iter::from_fn,
};

fn to_priority(c: char) -> u32 {
    let b = u32::from(c);

    (b & 0x1F) + if (b & 0x20 > 0) { 0 } else { 26 }
}

#[derive(Default)]
struct Ruck {
    left: HashSet<char>,
    right: HashSet<char>,
}

impl Ruck {
    fn insert(&mut self, (left, right): (char, char)) {
        self.left.insert(left);
        self.right.insert(right);
    }

    fn check_wrong_pocket(&self, (left, right): (char, char)) -> Option<char> {
        match (self.left.contains(&right), self.right.contains(&left)) {
            (true, false) => Some(right),
            (false, true) => Some(left),
            (false, false) => None,
            (true, true) => {
                if left == right {
                    Some(left)
                } else {
                    panic!("Uhhhhh")
                }
            }
        }
    }

    fn union(&self) -> HashSet<char> {
        self.left.union(&self.right).copied().collect()
    }

    fn intersection(self, other: &Self) -> HashSet<char> {
        self.union().intersection(&other.union()).copied().collect()
    }
}

fn main() -> io::Result<()> {
    let newline = String::from("\n");
    let mut line_iter = BufReader::new(File::open("./input")?).lines();

    let mut prioritysum: u32 = 0;

    let mut shared_in_group_iter = line_iter
        .map(|x| x.unwrap())
        .map(|mut line| {
            let mut h: HashSet<char> = HashSet::new();
            let mut line_drain = line.drain(..);
            let mut ruck = Ruck::default();
            from_fn(|| match (line_drain.next(), line_drain.next_back()) {
                (Some(f), Some(b)) => Some((f, b)),
                (None, None) => None,
                _ => panic!("non even rucksack size"),
            })
            .fold(ruck, |mut r, pair| {
                r.insert(pair);
                r
            })
        })
        .array_chunks::<3>()
        .map(|ruckset| {
            let Some(g) = ruckset
                .iter()
                .map(|x| x.union())
                .reduce(|acc, x| acc.intersection(&x).copied().collect()) else {
                    panic!("Did not produce a hashset")
                };
            g
        });

    for shared_set in shared_in_group_iter {
        let Some(v) = shared_set.into_iter().next() else {
            panic!("Wrong")
        };

        prioritysum += to_priority(v);
    }

    println!("Priority sum is {}", prioritysum);

    Ok(())
}

#[cfg(test)]
mod tests {
    use crate::to_priority;

    #[test]
    fn priority() {
        assert_eq!(to_priority('a'), 1);
        assert_eq!(to_priority('z'), 26);
        assert_eq!(to_priority('A'), 27);
        assert_eq!(to_priority('Z'), 52);
    }
}
