#![allow(unused)]
use std::{fs::File, io::{BufReader, self, BufRead}, iter::from_fn, fmt::Display, ops::{Add, AddAssign}};
use std::collections::BinaryHeap;

#[derive(Debug, Clone, Copy)]
struct Score(u32);

impl Add<Score> for Score {
    type Output = Score;

    fn add(self, rhs: Score) -> Self::Output {
        Score(self.0 + rhs.0)
    }
}

impl AddAssign<Score> for Score {
    fn add_assign(&mut self, rhs: Score) {
        *self = *self + rhs;
    }
}

enum Outcome {
    Win,
    Loss,
    Draw
}

impl Outcome {
    fn score(self) -> Score {
        match self {
            Outcome::Win => Score(6),
            Outcome::Loss => Score(0),
            Outcome::Draw => Score(3),
        }
    }
}

impl From<&str> for Outcome {
    fn from(value: &str) -> Self {
        match value {
            "X" => Outcome::Loss,
            "Y" => Outcome::Draw,
            "Z" => Outcome::Win,
            _ => panic!(),
        }
    }
}

#[derive(PartialEq)]
enum RPS {
    Rock,
    Paper,
    Scissors
}

impl RPS {
    /// x R P S
    /// R D W L
    /// P L D L
    /// S W L D
    fn vs(self, other: RPS) -> Outcome {
        match (self, other) {
            (RPS::Rock, RPS::Rock) => Outcome::Draw,
            (RPS::Rock, RPS::Paper) => Outcome::Loss,
            (RPS::Rock, RPS::Scissors) => Outcome::Win,
            (RPS::Paper, RPS::Rock) => Outcome::Win,
            (RPS::Paper, RPS::Paper) => Outcome::Draw,
            (RPS::Paper, RPS::Scissors) => Outcome::Loss,
            (RPS::Scissors, RPS::Rock) => Outcome::Loss,
            (RPS::Scissors, RPS::Paper) => Outcome::Win,
            (RPS::Scissors, RPS::Scissors) => Outcome::Draw,
        }
    }

    fn winning_move(self) -> RPS {
        match self {
            RPS::Rock => RPS::Paper,
            RPS::Paper => RPS::Scissors,
            RPS::Scissors => RPS::Rock,
        }
    }

    fn losing_move(self) -> RPS {
        match self {
            RPS::Rock => RPS::Scissors,
            RPS::Paper => RPS::Rock,
            RPS::Scissors => RPS::Paper,
        }
    }

    fn get_desired_outcome_shape(self, o: &Outcome) -> RPS {
        match o {
            Outcome::Win => self.winning_move(),
            Outcome::Loss => self.losing_move(),
            Outcome::Draw => self.drawing_move(),
        }
    }

    fn drawing_move(self) -> RPS {
        return self
    }

    fn shape_score(&self) -> Score {
        match self {
            RPS::Rock => Score(1),
            RPS::Paper => Score(2),
            RPS::Scissors => Score(3),
        }
    }
}

impl From<OpponentPlayStr> for RPS {
    fn from(value: OpponentPlayStr) -> Self {
        match value.0.as_str() {
            "A" => RPS::Rock,
            "B" => RPS::Paper,
            "C" => RPS::Scissors,
            _ => panic!()
        }
    }
}

impl From<MyPlayStr> for RPS {
    fn from(value: MyPlayStr) -> Self {
        match value.0.as_str() {
            "X" => RPS::Rock,
            "Y" => RPS::Paper,
            "Z" => RPS::Scissors,
            _ => panic!()
        }
    }
}

struct OpponentPlayStr(String);

struct MyPlayStr(String);


fn main() -> io::Result<()> {
    let newline = String::from("\n");
    let mut line_iter = BufReader::new(File::open("./input")?).lines();

    let mut my_score = Score(0);
    while let Some(Ok(line)) = line_iter.next() {
        let mut i =line.split(" ");
        let (Some(opponent_str), Some(my_str)) =(i.next(), i.next()) else {
            panic!("AHHHH");
        };
        let opponent = RPS::from(OpponentPlayStr(String::from(opponent_str)));
        let desired_outcome = Outcome::from(my_str);



        my_score += opponent.get_desired_outcome_shape(&desired_outcome).shape_score() + desired_outcome.score();
    }

    println!("My tournament score is: {:?}", my_score);

    Ok(())
}