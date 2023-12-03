#![feature(iter_array_chunks)]
use std::{
    fs::File,
    io::{BufRead, BufReader},
    str::FromStr,
};

mod cpu;
mod part1;
mod part2;
mod crt;

use cpu::Op;

fn main() {
    let instructions: Vec<Op> = BufReader::new(File::open("./input").unwrap())
        .lines()
        .flatten()
        .map(|s| Op::from_str(s.as_ref()))
        .flatten()
        .collect();

    part1::part_1(instructions.clone());

    part2::part_2(instructions);
}

#[cfg(test)]
mod test {
    use std::str::FromStr;

    use crate::{
        cpu::{Op, Register, CPU},
        part1::signal_strength, part2,
    };

    const SMALLPROGRAM: &'static str = "noop\naddx 3\naddx -5";
    const LARGERPROGRAM: &'static str = "addx 15\naddx -11\naddx 6\naddx -3\naddx 5\naddx -1\naddx -8\naddx 13\naddx 4\nnoop\naddx -1\naddx 5\naddx -1\naddx 5\naddx -1\naddx 5\naddx -1\naddx 5\naddx -1\naddx -35\naddx 1\naddx 24\naddx -19\naddx 1\naddx 16\naddx -11\nnoop\nnoop\naddx 21\naddx -15\nnoop\nnoop\naddx -3\naddx 9\naddx 1\naddx -3\naddx 8\naddx 1\naddx 5\nnoop\nnoop\nnoop\nnoop\nnoop\naddx -36\nnoop\naddx 1\naddx 7\nnoop\nnoop\nnoop\naddx 2\naddx 6\nnoop\nnoop\nnoop\nnoop\nnoop\naddx 1\nnoop\nnoop\naddx 7\naddx 1\nnoop\naddx -13\naddx 13\naddx 7\nnoop\naddx 1\naddx -33\nnoop\nnoop\nnoop\naddx 2\nnoop\nnoop\nnoop\naddx 8\nnoop\naddx -1\naddx 2\naddx 1\nnoop\naddx 17\naddx -9\naddx 1\naddx 1\naddx -3\naddx 11\nnoop\nnoop\naddx 1\nnoop\naddx 1\nnoop\nnoop\naddx -13\naddx -19\naddx 1\naddx 3\naddx 26\naddx -30\naddx 12\naddx -1\naddx 3\naddx 1\nnoop\nnoop\nnoop\naddx -9\naddx 18\naddx 1\naddx 2\nnoop\nnoop\naddx 9\nnoop\nnoop\nnoop\naddx -1\naddx 2\naddx -37\naddx 1\naddx 3\nnoop\naddx 15\naddx -21\naddx 22\naddx -6\naddx 1\nnoop\naddx 2\naddx 1\nnoop\naddx -10\nnoop\nnoop\naddx 20\naddx 1\naddx 2\naddx 2\naddx -6\naddx -11\nnoop\nnoop\nnoop";

    fn nprogram(program: &'static str) -> Vec<Op> {
        program.lines().map(Op::from_str).flatten().collect()
    }

    #[test]
    fn test_small_p1() {
        let prog = nprogram(SMALLPROGRAM);
        let mut cpu = CPU::new(prog);

        let expected = [
            (Some(()), Some(&1)),
            (Some(()), Some(&1)),
            (Some(()), Some(&1)),
            (Some(()), Some(&4)),
            (Some(()), Some(&4)),
            (None, Some(&-1)),
            (None, Some(&-1)),
        ];

        for (_cycle, exp) in expected.iter().copied().enumerate() {
            let ret = cpu.cycle();
            let x = cpu.program_memory.get(&Register::X);
            assert_eq!((ret, x), exp);
        }
    }

    #[test]
    fn test_large_p1() {
        let mut cpu = CPU::new(nprogram(LARGERPROGRAM));

        let expected = [
            (20, &21, 420),
            (60, &19, 1140),
            (100, &18, 1800),
            (140, &21, 2940),
            (180, &16, 2880),
            (220, &18, 3960),
        ];

        let mut cycle = 0;
        let mut signal_sum = 0;
        for exp in expected {
            while cycle < exp.0 {
                if let None = cpu.cycle() {
                    panic!("CPU died while cycling");
                }
                cycle += 1;
            }
            let x = cpu.program_memory.get(&Register::X).unwrap();
            let signal_strength = signal_strength(&cycle, x);

            assert_eq!((cycle, x, signal_strength), exp);
            signal_sum += signal_strength;
        }

        assert_eq!(signal_sum, 13140);
    }

    #[test]
    fn test_large_p2() {
        let instructions = nprogram(LARGERPROGRAM);

        part2::part_2(instructions);
    }
}


