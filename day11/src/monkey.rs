use std::{
    ops::{Add, Mul},
};

use nom::{
    branch::alt,
    bytes::complete::tag,
    character::complete::{digit1, line_ending, not_line_ending, one_of, space1},
    combinator::{all_consuming, map_res, opt},
    multi::separated_list0,
    sequence::{preceded, terminated, tuple},
    IResult, Parser,
};

pub(crate) type Worry = u64;

pub(crate) struct Monkey {
    pub(crate) index: Worry,
    pub(crate) worry_op: WorryOperation,
    pub(crate) test: MonkeyTest,
}

impl Monkey {
    fn parse_items(input: &str) -> IResult<&str, Vec<Worry>> {
        let p = parse_line
            .and_then(preceded(
                tuple((opt(tag("  ")), tag("Starting items: "))),
                all_consuming(separated_list0(tag(", "), parse_digits)),
            ))
            .parse(input);
        p
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
            .and(alt((
                tag("old").map(|_| WorryVal::Old),
                parse_digits.map(|v| WorryVal::Val(v)),
            )))
            .map(|(op, val)| WorryOperation { op, val });

        parse_line
            .and_then(preceded(
                tuple((opt(tag("  ")), tag("Operation: new = old "))),
                parse_operation,
            ))
            .parse(input)
    }

    pub(crate) fn parse_monkey(input: &str) -> IResult<&str, (Vec<Worry>, Monkey)> {
        let (input, monkeyindex) = parse_line
            .and_then(preceded(tag("Monkey "), terminated(parse_digits, tag(":"))))
            .parse(input)?;

        // let (input, items) = Self::parse_items(input)?;

        // let (input, worry_op) = Self::parse_worry_operation(input)?;

        // let (input, test) = MonkeyTest::parse_monkey_test(input)?;

        let (input, (items, worry_op, test)) = tuple((
            Self::parse_items,
            Self::parse_worry_operation,
            MonkeyTest::parse_monkey_test,
        ))(input)?;

        Ok((
            input,
            (items, Monkey {
                index: monkeyindex,
                worry_op,
                test,
            }),
        ))
    }
}

enum WorryVal {
    Val(Worry),
    Old,
}

pub(crate) struct WorryOperation {
    op: fn(Worry, Worry) -> Worry,
    val: WorryVal,
}

impl WorryOperation {
    pub(crate) fn exec(&self, w: Worry) -> Worry {
        (self.op)(
            w,
            match self.val {
                WorryVal::Val(v) => v,
                WorryVal::Old => w,
            },
        )
    }
    pub(crate) fn exec_tup(t: (&Self, Worry)) -> Worry {
        t.0.exec(t.1)
    }
}

#[derive(Debug, PartialEq)]
pub(crate) struct MonkeyTest {
    pub(crate) modulo: Worry,
    pub(crate) true_dst: Worry,
    pub(crate) false_dst: Worry,
}

impl MonkeyTest {
    fn parse_modulo(input: &str) -> IResult<&str, Worry> {
        parse_line
            .and_then(preceded(
                tag("  Test: divisible by "),
                parse_digits_whole_str,
            ))
            .parse(input)
    }

    fn parse_condition_line(cond: bool) -> impl Fn(&str) -> IResult<&str, Worry> {
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
            parse_line.and_then(Self::parse_condition_line(true)),
            parse_line.and_then(Self::parse_condition_line(false)),
        ))(input)
    }

    fn parse_monkey_test(input: &str) -> IResult<&str, MonkeyTest> {
        let (input, (modulo, (true_dst, false_dst))) =
            tuple((Self::parse_modulo, Self::parse_conditions))(input)?;

        Ok((
            input,
            MonkeyTest {
                modulo,
                true_dst,
                false_dst,
            },
        ))
    }
}

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
    use nom::error::{Error, ErrorKind};

    use crate::monkey::{Monkey, MonkeyTest};

    use super::{parse_line, Worry};

    const REGULAR_MONKEY: &'static str = r#"Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0"#;

    const OLD_OP_MONKEY: &'static str = r#"Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3"#;

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
            Monkey::parse_items("  Starting items: 54, 65, 75, 74\nblah"),
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
        let test_str = REGULAR_MONKEY
            .split_inclusive("\n")
            .skip(2)
            .collect::<String>();

        let (input, worryop) = Monkey::parse_worry_operation(test_str.as_str()).unwrap();

        assert_eq!(
            input,
            "  Test: divisible by 19\n    If true: throw to monkey 2\n    If false: throw to monkey 0"
        );

        assert_eq!(worryop.exec(0), 6);
    }

    #[test]
    fn test_parse_monkey_test_condition() {
        const TEST_STR: &'static str =
            "    If true: throw to monkey 2\n    If false: throw to monkey 0";
        assert_eq!(MonkeyTest::parse_conditions(TEST_STR), Ok(("", (2, 0))))
    }

    #[test]
    fn test_parse_monkey_test() {
        const TEST_STR: &'static str ="  Test: divisible by 13\n    If true: throw to monkey 1\n    If false: throw to monkey 3\nJK\n";

        assert_eq!(
            MonkeyTest::parse_monkey_test(TEST_STR),
            Ok((
                "JK\n",
                MonkeyTest {
                    modulo: 13,
                    true_dst: 1,
                    false_dst: 3
                }
            ))
        )
    }

    fn check_monkey_full(
        monkey: &Monkey,
        index: Worry,
        opresult: Worry,
        modulo: Worry,
        true_dst: Worry,
        false_dst: Worry,
    ) {
        assert_eq!(monkey.index, index);

        let wop = &monkey.worry_op;

        assert_eq!(wop.exec(100), opresult);

        assert_eq!(
            monkey.test,
            MonkeyTest {
                modulo: modulo,
                true_dst: true_dst,
                false_dst: false_dst,
            }
        );
    }

    #[test]
    fn test_parse_monkey() {
        let (remaining, (items, monkey)) = Monkey::parse_monkey(REGULAR_MONKEY).unwrap();

        assert_eq!(remaining.len(), 0);

        assert_eq!(items, vec![54, 65, 75, 74]);

        check_monkey_full(&monkey, 1, 0, 106, 19, 2);

        let (remaining, (items, monkey)) = Monkey::parse_monkey(OLD_OP_MONKEY).unwrap();

        assert_eq!(remaining.len(), 0);

        assert_eq!(items, vec![79, 60, 97]);

        check_monkey_full(&monkey, 2, 10000, 13, 1, 3);
    }
}
