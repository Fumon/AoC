use std::{fs::File, io::{BufRead, BufReader}};



fn main() {
    let bf = BufReader::new(File::open("input/input").unwrap());
    
    println!("{}", part1(bf.lines().into_iter().filter_map(Result::ok)));
}

fn part1(lines: impl IntoIterator<Item = impl AsRef<str>>) -> i64 {
    let mut digitsum: i64 = 0;
    for line in lines {
        let numbers: Vec<_> = line.as_ref().chars().filter_map(|ch| {
            match ch {
                '0'..='9' => Some(ch as i64 - '0' as i64),
                _ => None
            }
        }).collect();
        digitsum += numbers.first().unwrap() * 10 + numbers.last().unwrap()
    }

    return digitsum;
}