#![feature(iter_array_chunks)]
use std::{
    fs::File,
    io::{self, BufRead, BufReader},
};

trait Interval<T>
where
    T: ?Sized + PartialOrd<T>,
{
    fn start(&self) -> &T;
    fn end(&self) -> &T;

    fn contains(&self, other: &T) -> bool {
        self.start() <= other && self.end() >= other
    }

    fn covers(&self, other: impl Interval<T>) -> bool {
        self.contains(other.start()) && self.contains(other.end())
    }
}

impl Interval<u32> for (u32, u32) {
    fn start(&self) -> &u32 {
        &self.0
    }

    fn end(&self) -> &u32 {
        &self.1
    }
}

fn main() -> io::Result<()> {
    // Open file
    // Shove into buffered reader
    // Split by lines
    // Separate by commas
    // Separate by hyphen
    // Parse into int pairs
    // Check if either is contained within the other
    // Sum the result

    let num_contains: u32 = BufReader::new(File::open("./input")?)
        .lines()
        .flatten()
        .map(|x| {
            let elf_pair = x
                .split_once(",")
                .map(|pair| [pair.0, pair.1])
                .unwrap()
                .map(|x| x.split_once("-").unwrap());

            elf_pair.map(|(a, b)| {
                let (Ok(i), Ok(j)) =  (a.parse::<u32>(), b.parse::<u32>()) else {
                    panic!("Numeric conversion failed")
                };
                (i, j)
            })
        })
        .map(|[x, y]| if x.covers(y) || y.covers(x) { 1 } else { 0 })
        .sum();

    println!("Number of fully contained subsets {}", num_contains);
    Ok(())
}
