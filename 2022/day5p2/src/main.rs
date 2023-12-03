#![allow(unused)]
#![feature(iter_array_chunks)]
#![feature(iter_advance_by)]
#![feature(iter_intersperse)]

use core::panic;
use std::{
    fs::File,
    io::{self, BufRead, BufReader},
    iter::from_fn,
    str::Chars,
};

struct ShipStacksState(Vec<String>);

struct ShipStacks<T> {
    stacks: Vec<T>,
}

impl ShipStacks<String> {
    fn execute_move(&mut self, m: Move) {
        let split_ind = self.stacks[m.from].len() - m.count;
        let elf_crates: Vec<char> = self.stacks[m.from].drain(split_ind..).collect();
        for elf_crate in elf_crates {
            self.stacks[m.to].push(elf_crate);
        }
    }

    fn get_top_row(&self) -> String {
        self.stacks
            .iter()
            .map(|stack| {
                if let Some(elf_crate) = stack.chars().last() {
                    elf_crate
                } else {
                    ' '
                }
            })
            .collect()
    }
}

impl From<ShipStacksState> for ShipStacks<String> {
    fn from(value: ShipStacksState) -> Self {
        let mut iter = value.0.into_iter().rev();

        let stack_count = iter
            .next()
            .expect("there should be at least one line in the state")
            .split_ascii_whitespace()
            .next_back()
            .expect("There should be at least one stack")
            .parse::<usize>()
            .expect("The last stack should have a parsable item");

        // intialize
        let mut sstacks = ShipStacks {
            stacks: (0..stack_count).map(|_| String::new()).collect(),
        };

        let stacks = &mut sstacks.stacks;

        // Parse each line
        // Careful of empty stacks
        for line in iter {
            // Convert into a friendlier iterator
            let mut char_iter = line.chars().skip(1);
            let mut stack_index = 0;
            let fill_state_iter = from_fn(|| loop {
                let Some(c) = char_iter.next() else {
                        return None;
                    };

                let r = Some((stack_index, c));
                let advance_result = char_iter.advance_by(3);
                stack_index += 1;

                if c != ' ' {
                    return r;
                } else {
                    if let Ok(_) = advance_result {
                        continue;
                    } else {
                        return None;
                    }
                }
            });

            for (stack_index, char) in fill_state_iter {
                stacks[stack_index].push(char)
            }
        }

        sstacks
    }
}

impl std::fmt::Display for ShipStacks<String> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let baseline = (1..=self.stacks.len())
            .map(|x| format!("{: ^3}", x))
            .intersperse(" ".to_string())
            .collect();

        let mut acc: Vec<String> = vec![baseline];
        let mut stack_iters: Vec<Chars> = self.stacks.iter().map(|stack| stack.chars()).collect();
        loop {
            let elf_crates: Vec<_> = stack_iters.iter_mut().map(Chars::next).collect();
            if elf_crates.iter().any(Option::is_some) {
                acc.push(
                    elf_crates
                        .into_iter()
                        .map(|co| co.map_or("   ".to_string(), |c| format!("[{}]", c)))
                        .intersperse(" ".to_string())
                        .collect(),
                );
            } else {
                break;
            }
        }

        let displayline = acc.into_iter().rev().intersperse("\n".to_string()).collect::<String>();

        write!(f, "{}", displayline);

        Ok(())
    }
}

#[derive(Debug)]
pub struct Move {
    pub count: usize,
    pub from: usize,
    pub to: usize,
}

impl From<String> for Move {
    fn from(line: String) -> Self {
        let [count, from, to]: [usize; 3] = line
            .split_ascii_whitespace()
            .skip(1)
            .step_by(2)
            .flat_map(|x| x.parse::<usize>())
            .take(3)
            .collect::<Vec<usize>>()
            .try_into()
            .ok()
            .expect("Couldn't unpack move properly");

        Move {
            count,
            from: (from - 1),
            to: (to - 1),
        }
    }
}

fn main() -> io::Result<()> {
    let mut line_iter = BufReader::new(File::open("./input")?).lines().flatten();

    let stacks_start_state = ShipStacksState(
        line_iter
            .by_ref()
            .take_while(|l| l.len() > 0)
            .collect::<Vec<String>>(),
    );

    let mut stacks = ShipStacks::from(stacks_start_state);

    line_iter
        .map(Move::from)
        .for_each(|m| stacks.execute_move(m));

    println!("Top of the stacks: {}", stacks.get_top_row());

    Ok(())
}

#[cfg(test)]
mod test {
    use std::{iter::from_fn, ops::Rem};

    #[test]
    fn take_while_by_ref_does_not_consume_iter() {
        let g = ["a", "b", "c", " ", "1", "2", "3"];

        let mut giter = g.into_iter();

        let first_take: String = giter.by_ref().take_while(|c| *c != " ").collect();

        assert_eq!(first_take, String::from("abc"));

        let second_take: String = giter.by_ref().collect();

        assert_eq!(second_take, String::from("123"));
        // Although, note that the " " is never emitted as `take_while` took it out and threw it away
        // For a method which does emit the separator: see `next_if`
    }

    #[test]
    fn next_if_while_let() {
        let g: &[u32] = [2, 4, 6, 8, 3, 5, 7, 9].as_slice();

        let mut giter = g.into_iter().peekable();

        let firsttake: Vec<u32> = from_fn(|| giter.next_if(|x| (x.rem(2) == 0)))
            .copied()
            .collect();

        assert_eq!(firsttake, [2, 4, 6, 8].as_slice());

        let secondtake: Vec<u32> = giter.copied().collect();

        assert_eq!(secondtake, [3, 5, 7, 9].as_slice());
    }
}
