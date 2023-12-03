#![allow(unused)]

use std::{io::{self,BufReader, BufRead}, fs::File, iter::from_fn, collections::HashSet};

fn to_priority(c: char) -> u32 {

    let b = u32::from(c);

    (b & 0x1F) + if (b & 0x20 > 0) {0} else {26}
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
            },
        }
    }
}

fn main() -> io::Result<()> {
    let newline = String::from("\n");
    let mut line_iter = BufReader::new(File::open("./input")?).lines();

    let mut prioritysum: u32 = 0;
    while let Some(Ok(mut line)) = line_iter.next() {

        let mut h: HashSet<char> = HashSet::new();
        let mut linedrain = line.drain(..);
        let mut doubleiter = from_fn(|| {
            match (linedrain.next(), linedrain.next_back()) {
                (Some(f), Some(b)) => Some((f,b)),
                (None, None) => None,
                _ => panic!("non even rucksack size")
            }
        });
        
        let mut ruck = Ruck::default();

        for pair in doubleiter {
            ruck.insert(pair);
            if let Some(odd) = ruck.check_wrong_pocket(pair) {
                prioritysum += to_priority(odd);
                break;
            };
        }
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