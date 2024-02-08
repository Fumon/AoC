use std::{fs::File, io::{BufRead, BufReader}};

use regex::{Regex, RegexSet};

fn main() {
    let bf = BufReader::new(File::open("input/input").unwrap());
    println!("{}", part1(bf.lines().into_iter().filter_map(Result::ok)));
    let bf = BufReader::new(File::open("input/input").unwrap());
    println!("{}", part2(bf.lines().into_iter().filter_map(Result::ok)));
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

fn part2(lines: impl IntoIterator<Item = impl AsRef<str>>) -> i64 {
    let set = RegexSet::new(&[
        r"0",
        r"(1|one)",
        r"(2|two)",
        r"(3|three)",
        r"(4|four)",
        r"(5|five)",
        r"(6|six)",
        r"(7|seven)",
        r"(8|eight)",
        r"(9|nine)",
    ]).unwrap();

    let pattern = &Regex::new(r"[0-9]|one|two|three|four|five|six|seven|eight|nine").unwrap();
    let mut digitsum: i64 = 0;
    for line in lines {
        let r = line.as_ref();
        let mut start = 0;
        let linenumber = std::iter::from_fn(move || {
            if let Some(rmatch) = pattern.find(&r[start..]) {
                start += rmatch.start() + 1;
                return Some(rmatch.as_str());
            } else {
                return None
            }
        }).fold(None, |acc, val| {
            if let Some((first, _)) = acc {
                Some((first, val))
            } else {
                Some((val, val))
            }
        }).map_or(0, |(first, last)| {
            let firstmatch = set.matches(first);
            let lastmatch = set.matches(last);
            firstmatch.into_iter().next().unwrap() * 10 + lastmatch.into_iter().next().unwrap()
        });
        digitsum += linenumber as i64;
    }

    digitsum
}


#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn part1_testinput() {
        assert_eq!(part1(vec!["apple1pear", "11"]), 22);

        let bf = BufReader::new(File::open("input/testinput01").unwrap()).lines().into_iter().filter_map(Result::ok);
        assert_eq!(part1(bf), 142);
    }

    #[test]
    fn part2_testinput() {
        let bf = BufReader::new(File::open("input/testinput02").unwrap()).lines().into_iter().filter_map(Result::ok);
        assert_eq!(part2(bf), 281);
    }
}