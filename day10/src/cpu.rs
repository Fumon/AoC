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
