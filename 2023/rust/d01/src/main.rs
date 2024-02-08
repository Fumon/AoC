use std::{fs::File, io::{BufRead, BufReader}};



fn main() {
    let bf = BufReader::new(File::open("input/input").unwrap());

    let mut digitsum: u32 = 0;
    for line in bf.lines().into_iter().filter_map(Result::ok) {
        let numbers: Vec<u32> = line.chars().filter_map(|ch| {
            match ch {
                '0'..='9' => Some(ch as u32 - '0' as u32),
                _ => None
            }
        }).collect();
        digitsum += numbers.first().unwrap() * 10 + numbers.last().unwrap()
    }

    println!("{}", digitsum);
}