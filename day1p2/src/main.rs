#![allow(unused)]
use std::{fs::File, io::{BufReader, self, BufRead}, iter::from_fn};
use std::collections::BinaryHeap;

fn main() -> io::Result<()> {
    let newline = String::from("\n");
    let mut line_iter = BufReader::new(File::open("./input")?).lines();

    let mut h = BinaryHeap::new();
    let mut sum = 0;
    loop {
        let line = line_iter.next();
        match line {
            Some(Err(a)) => return Err(a),
            None => {h.push(sum); sum = 0; break;},
            Some(Ok(l)) if l.len() == 0 => {h.push(sum); sum = 0;},
            Some(Ok(l)) => sum += l.parse::<u32>().unwrap(),
        }
    }

    println!("Sum of top three elves' Calories: {}", from_fn(|| {h.pop()}).take(3).sum::<u32>());

    Ok(())
}
