#![allow(unused)]

use std::{
    num::ParseIntError,
    ops::{Add, Mul},
    str::FromStr,
};

use nom::{
    bytes::complete::{is_not, tag, take_until1, take_while, take_while1},
    character::complete::{digit1, line_ending, newline, not_line_ending, one_of, space0, space1},
    combinator::{all_consuming, flat_map, map, map_parser, map_res, opt, recognize},
    error::{self, Error, ErrorKind},
    multi::{many1, separated_list0},
    sequence::{preceded, terminated, tuple},
    Finish, IResult, Parser,
};

type Worry = u32;

pub(crate) struct Monkey {
    items: Vec<Worry>,
    worry_op: WorryOperation,
    test: MonkeyTest,
    activity: u32,
}

impl Monkey {
    fn parse_items(input: &str) -> IResult<&str, Vec<Worry>> {
        map_parser(
            parse_line,
            preceded(
                tuple((opt(tag("  ")), tag("Starting items: "))),
                all_consuming(separated_list0(tag(", "), parse_digits)),
            ),
        )(input)
    }

    fn parse_worry_operation(input: &str) -> IResult<&str, WorryOperation> {
        let parse_op_char = map_res(
            terminated(one_of("*+"), space1),
            |opc| -> Result<fn(Worry, Worry) -> Worry, &str> {
                match opc {
                    '*' => Ok(Worry::mul),
                    '+' => Ok(Worry::add),
                    _ => Err("Invalid Operator"),
                }
            },
        );

        let parse_operation = parse_op_char
            .and(parse_digits)
            .map(|(op, val)| WorryOperation { op, val });

        preceded(
            tuple((opt(tag("  ")), tag("Operation: new = old "))),
            parse_operation,
        )(input)
    }
}

struct WorryOperation {
    op: fn(Worry, Worry) -> Worry,
    val: Worry,
}

impl WorryOperation {
    fn exec(&self, w: Worry) -> Worry {
        (self.op)(w, self.val)
    }
}

#[derive(Debug)]
struct MonkeyTest {
    modulo: Worry,
    true_dst: u32,
    false_dst: u32,
}

impl MonkeyTest {
    fn parse_modulo(input: &str) -> IResult<&str, Worry> {
        map_parser(
            parse_line,
            preceded(tag("  Test: divisible by "), parse_digits_whole_str)
        )(input)
    }

    fn parse_condition_line(cond: bool) -> (impl Fn(&str) -> IResult<&str, Worry>) {
        move |input: &str| {
            preceded(
                tuple((
                    tag("    If "),
                    tag(cond.to_string().as_str()),
                    tag(": throw to monkey "),
                )),
                parse_digits_whole_str,
            )(input)
        }
    }

    fn parse_conditions(input: &str) -> IResult<&str, (Worry, Worry)> {
        tuple((
            map_parser(parse_line, Self::parse_condition_line(true)),
            map_parser(parse_line, Self::parse_condition_line(false)),
        ))(input)
    }

    fn parse_monkey_test(input: &str) -> IResult<&str, (Worry, (Worry, Worry))> {
        // let (input, (modulo, (true_dst, false_dst))) =
        tuple((Self::parse_modulo, Self::parse_conditions))(input)
    }
}

// impl FromStr for MonkeyTest {
//     type Err = nom::error::Error<String>;

//     fn from_str(s: &str) -> Result<Self, Self::Err> {
        
//         match Self::parse_monkey_test(s).finish() {
//             Ok((_, (modulo, (true_dst, false_dst)))) => Ok(MonkeyTest{modulo, true_dst, false_dst}),
//             Err(Error{input, code}) => Err(Error {
//                 input: input.to_string(),
//                 code,
//             })
//         };

//         todo!()
//     }
// }

fn parse_line(input: &str) -> IResult<&str, &str> {
    terminated(not_line_ending, opt(line_ending))(input)
}

fn parse_digits(input: &str) -> IResult<&str, Worry> {
    map_res(digit1, |n: &str| n.parse())(input)
}

fn parse_digits_whole_str(input: &str) -> IResult<&str, Worry> {
    all_consuming(parse_digits)(input)
}

#[cfg(test)]
mod test {
    use std::panic;

    use nom::error::{Error, ErrorKind};

    use crate::monkey::{Monkey, MonkeyTest};

    use super::parse_line;

    #[test]
    fn test_parse_line() {
        assert_eq!(
            parse_line("hatoful boyfriend\nfdsafdsa\nanother one"),
            Ok(("fdsafdsa\nanother one", "hatoful boyfriend"))
        );
        assert_eq!(parse_line("nothing after"), Ok(("", "nothing after")));
    }

    #[test]
    fn test_parse_monkey_items() {
        assert_eq!(
            Monkey::parse_items("Starting items: 54, 65, 75, 74"),
            Ok(("", vec![54, 65, 75, 74]))
        );

        assert_eq!(
            Monkey::parse_items("Starting items: 54, 65, 75, 74\nblah"),
            Ok(("blah", vec![54, 65, 75, 74]))
        );

        assert_eq!(
            Monkey::parse_items("Starting items: 54, 65, AB, 74"),
            Err(nom::Err::Error(Error {
                input: ", AB, 74",
                code: ErrorKind::Eof
            }))
        );

        assert_eq!(
            Monkey::parse_items("Starting items: 1, 2, 3, 534523465435643254325342, 4, 9"),
            Err(nom::Err::Error(Error {
                input: ", 534523465435643254325342, 4, 9",
                code: ErrorKind::Eof
            }))
        );
    }

    #[test]
    fn test_parse_worry_operation() {
        const TEST_STR: &'static str =
            "    If true: throw to monkey 2\n    If false: throw to monkey 0";
        assert_eq!(MonkeyTest::parse_conditions(TEST_STR), Ok(("", (2, 0))))
    }

    #[test]
    fn test_parse_monkey_test_condition() {
        const TEST_STR: &'static str =
            "    If true: throw to monkey 2\n    If false: throw to monkey 0";
        assert_eq!(MonkeyTest::parse_conditions(TEST_STR), Ok(("", (2, 0))))
    }

    #[test]
    fn test_parse_monkey_test() {
        const TEST_STR: &'static str ="  Test: divisible by 13\n    If true: throw to monkey 1\n    If false: throw to monkey 3";

        assert_eq!(
            MonkeyTest::parse_monkey_test(TEST_STR),
            Ok(("", (13, (1, 3))))
        )
    }
}
