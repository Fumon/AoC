#![allow(unused)]
use std::{fs::File, io::{BufReader, self, BufRead}};

fn main() -> io::Result<()> {
    let newline = String::from("\n");
    let mut line_iter = BufReader::new(File::open("./input")?).lines();

    let mut max = 0;
    let mut sum = 0;
    let mut update = |max: &mut u32, sum: &mut u32| {*max = if *sum > *max {*sum} else {*max}; *sum = 0;};
    loop {
        let line = line_iter.next();
        match line {
            Some(Err(a)) => return Err(a),
            None => {update(&mut max, &mut sum); break;},
            Some(Ok(l)) if l.len() == 0 => update(&mut max, &mut sum),
            Some(Ok(l)) => sum += l.parse::<u32>().unwrap(),
        }
    }

    println!("{}", max);

    Ok(())
}
