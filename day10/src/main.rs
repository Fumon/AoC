use std::{
    fs::File,
    io::{BufRead, BufReader},
    str::FromStr,
};

use cpu::Op;

fn main() {
    let instructions: Vec<Op> = BufReader::new(File::open("./input").unwrap())
        .lines()
        .flatten()
        .map(|s| Op::from_str(s.as_ref()))
        .flatten()
        .collect();

    part1::part_1(instructions);
}

mod part1 {
    use crate::cpu::{Op, CPU};

    pub(crate) const SIGNAL_STRENGTH_SAMPLE_CYCLES: [usize; 6] = [20, 60, 100, 140, 180, 220];

    pub(crate) fn part_1(instructions: Vec<Op>) {
        let mut cpu = CPU::new(instructions);
        
        let signal_sum: i64 = SIGNAL_STRENGTH_SAMPLE_CYCLES.iter().scan(0 , |cycle, target_cycle| {
            while *cycle < *target_cycle {
                if cpu.cycle() == None {
                    panic!("CPU died while iterating");
                }
                *cycle += 1;
            }
            let x = cpu.program_memory.get(&crate::cpu::Register::X).expect("retrieval of X is always possible");
            let signal_strength = signal_strength(&(*cycle).try_into().unwrap(), x);

            Some(signal_strength)
        }).sum();

        println!("The signal sum is {}", signal_sum);
    }

    pub(crate) fn signal_strength(cycle: &u32, xval: &i32) -> i64 {
        (*cycle as i64) * (*xval as i64)
    }
}



#[cfg(test)]
mod test {
    use std::str::FromStr;

    use crate::{
        cpu::{Op, Register, CPU},
        part1::signal_strength,
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
}

mod cpu {
    use std::{collections::HashMap, str::FromStr};

    type ProgramMemory = HashMap<Register, i32>;

    #[derive(Debug, Hash, PartialEq, Eq)]
    pub(crate) enum Register {
        X,
    }

    pub(crate) struct CPU {
        pub(crate) program_memory: ProgramMemory,
        instruction_memory: Vec<Op>,
        instruction_pointer: usize,
        instruction_cache: Op,
        instruction_phase: u32,
    }

    impl CPU {
        pub(crate) fn new(instructions: Vec<Op>) -> Self {
            let mem = ProgramMemory::from([(Register::X, 1)]);

            CPU {
                program_memory: mem,
                instruction_memory: instructions,
                instruction_pointer: 0,
                instruction_cache: Op::Noop,
                instruction_phase: 1,
            }
        }

        pub(crate) fn cycle(&mut self) -> Option<()> {
            // End of previous cycle
            self.instruction_phase = match self.instruction_phase.checked_sub(1) {
                Some(0) => {
                    self.instruction_cache.exec(&mut self.program_memory);
                    0
                }
                Some(v) => v,
                None => 0,
            };

            // Start of new cycle
            if self.instruction_phase == 0 {
                if self.fetch_instruction().is_none() {
                    return None;
                }
                self.instruction_pointer += 1;
            }

            // During cycle
            Some(())
        }

        fn fetch_instruction(&mut self) -> Option<()> {
            // Fetch instruction
            let Some(op) = self.instruction_memory.get(self.instruction_pointer) else {
                return None
            };

            self.instruction_cache = op.to_owned();
            self.instruction_phase = op.cycles();
            Some(())
        }
    }

    #[derive(Clone, Copy)]
    pub(crate) enum Op {
        Noop,
        Addx(i32),
    }

    impl Op {
        /// Returns how many cycles this operation will take before finishing
        const fn cycles(&self) -> u32 {
            match self {
                Self::Noop => 1,
                Self::Addx(_) => 2,
            }
        }

        fn exec(&self, mem: &mut ProgramMemory) {
            if let Self::Addx(v) = self {
                *mem.entry(Register::X).or_default() += v;
            }
        }
    }

    impl FromStr for Op {
        type Err = &'static str;

        fn from_str(s: &str) -> Result<Self, Self::Err> {
            match s.split(" ").collect::<Vec<_>>()[..] {
                ["noop"] => Ok(Self::Noop),
                ["addx", v] => Ok(Self::Addx(v.parse::<i32>().unwrap())),
                _ => Err("Invalid instruction"),
            }
        }
    }
}
