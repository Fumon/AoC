#![allow(unused)]

use std::{
    num::ParseIntError,
    ops::{Add, Mul},
    str::FromStr,
};

use nom::{
    bytes::complete::{is_not, tag, take_until1, take_while, take_while1},
    character::complete::{line_ending, newline, not_line_ending, one_of, space0, digit1},
    combinator::{all_consuming, map_res, opt, recognize},
    error::{self, ErrorKind},
    multi::{many1, separated_list0},
    sequence::{preceded, terminated, tuple},
    IResult, Parser,
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
        let (input, line) = parse_line(input)?;

        let (line, _) = preceded(opt(tag("  ")), tag("Starting items: "))(line)?;

        let (_, nums) = all_consuming(|f| {
            separated_list0(
                tag(", "),
                map_res(recognize(many1(one_of("0123456789"))), |x: &str| x.parse()),
            )(f)
        })(line)?;

        Ok((input, nums))
    }

    fn parse_worry_operation(input: &str) -> IResult<&str, WorryOperation> {
        let (input, _) = tag("Operation: new = old ")(input)?;

        let (input, opc) = one_of("*+")(input)?;

        let op = match opc {
            '*' => Worry::mul,
            '+' => Worry::add,
            _ => {
                return Err(nom::Err::Error(error::Error {
                    code: ErrorKind::Fail,
                    input,
                }))
            }
        };

        let (input, val): (&str, Worry) =
            map_res(recognize(many1(one_of("0123456789"))), |x: &str| x.parse())(input)?;

        Ok((input, WorryOperation { op, val }))
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
    fn parse_digits(input: &str) -> IResult<&str, Worry> {
        all_consuming(map_res(recognize(digit1), |n: &str| n.parse()))(input)
    }

    fn parse_modulo(input: &str) -> IResult<&str, Worry> {
        let (input, line) = parse_line(input)?;

        let (_, modulo) = preceded(
            tag("  Test: divisible by "),
            Self::parse_digits,
        )(line)?;

        Ok((input, modulo))
    }

    fn parse_conditions(input: &str) -> IResult<&str, (Worry, Worry)> {
        
        (vec![(line_ending, "true"), (line_ending, "false")]).into_iter()
        .map(|(linecomb, cond)| {
            linecomb.map(|f| preceded(tuple((tag("    If"),tag(cond), tag("throw to monkey "))), Self::parse_digits))
        })

        todo!()
    }
}

impl FromStr for MonkeyTest {
    type Err = nom::error::Error<String>;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        // let (input, (modulo, true_dst, false_dst)) = tuple((

        // ))?;

        todo!()
    }
}

fn parse_line(input: &str) -> IResult<&str, &str> {
    terminated(not_line_ending, opt(line_ending))(input)
}

#[cfg(test)]
mod test {
    use std::panic;

    use nom::error::{Error, ErrorKind};

    use crate::monkey::Monkey;

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
}
