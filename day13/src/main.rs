use std::{fs::read_to_string, iter::zip};

use packets::{Value, parse_full_packet_input};

fn main() {
    let ppairs = parse_full_packet_input(read_to_string("./input").unwrap().as_str());

    part_1(&ppairs)
}

fn part_1(ppairs: &Vec<[Value;2]>) {
    let s: usize = zip(1.., ppairs.iter()).filter_map(|(i, [l, r])| if l < r { Some(i) } else { None }).sum();
    println!("The sum of indicies of in-order pairs of packets is:\n\t{}", s);
}

#[cfg(test)]
mod test {
    use std::iter::zip;

    use crate::packets::parse_full_packet_input;

    const EXAMPLE_1: &'static str = r#"[1,1,3,1,1]
[1,1,5,1,1]

[[1],[2,3,4]]
[[1],4]

[9]
[[8,7,6]]

[[4,4],4,4]
[[4,4],4,4,4]

[7,7,7,7]
[7,7,7]

[]
[3]

[[[]]]
[[]]

[1,[2,[3,[4,[5,6,7]]]],8,9]
[1,[2,[3,[4,[5,6,0]]]],8,9]"#;

    #[test]
    fn p1_example() {
        let ppairs = parse_full_packet_input(EXAMPLE_1);

        let s: usize = zip(1.., ppairs.into_iter())
            .filter_map(|(i, [l, r])| if l < r { Some(i) } else { None })
            .sum();
        assert_eq!(13, s);
    }
}

mod packets {
    use std::cmp::Ordering;

    use nom::{
        branch::alt,
        character::complete::{char, digit1, line_ending},
        combinator::{all_consuming, map_res, opt},
        multi::{many_m_n, separated_list0, separated_list1},
        sequence::{delimited, terminated, tuple},
        IResult, Parser,
    };

    #[derive(Debug, Clone, PartialEq, Eq)]
    pub(crate) enum Value {
        List(Vec<Value>),
        Int(u32),
    }

    impl Value {
        fn iterate(&self) -> ValueIter {
            ValueIter { v: self, index: 0 }
        }
    }

    impl<'a> Ord for Value {
        fn cmp(&self, other: &Self) -> Ordering {
            match (self, other) {
                (Value::Int(l), Value::Int(r)) => l.cmp(r),
                (Value::List(_), Value::Int(_)) | (Value::Int(_), Value::List(_)) => {
                    self.iterate().cmp(other.iterate())
                }
                (Value::List(l), Value::List(r)) => l.cmp(r),
            }
        }
    }

    impl<'a> PartialOrd for Value {
        fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
            Some(self.cmp(other))
        }
    }

    fn parse_packet(input: &str) -> IResult<&str, Value> {
        delimited(
            char('['),
            separated_list0(
                char(','),
                alt((
                    map_res(digit1, |d: &str| d.parse::<u32>()).map(|num| Value::Int(num)),
                    parse_packet,
                )),
            ),
            char(']'),
        )
        .map(move |vals| Value::List(vals))
        .parse(input)
    }

    fn parse_one_packet(line: &str) -> Value {
        let ("", val) = all_consuming(parse_packet)(line).unwrap() else {
            panic!()
        };
        val
    }

    fn parse_packet_pair(input: &str) -> IResult<&str, [Value; 2]> {
        map_res(
            many_m_n(2, 2, terminated(parse_packet, opt(line_ending))),
            |v| v.try_into(),
        )(input)
    }

    pub(crate) fn parse_full_packet_input(input: &str) -> Vec<[Value; 2]> {
        let ("", ppairs) = all_consuming(separated_list1(line_ending, parse_packet_pair))(input).unwrap() else {
            panic!()
        };
        ppairs
    }

    struct ValueIter<'a> {
        v: &'a Value,
        index: usize,
    }

    impl<'a> Iterator for ValueIter<'a> {
        type Item = &'a Value;

        fn next(&mut self) -> Option<Self::Item> {
            match self.v {
                Value::List(sli) => {
                    if self.index >= sli.len() {
                        None
                    } else {
                        let out = Some(&sli[self.index]);
                        self.index += 1;
                        out
                    }
                }
                Value::Int(_) => {
                    if self.index > 0 {
                        None
                    } else {
                        self.index += 1;
                        Some(self.v)
                    }
                }
            }
        }
    }

    #[cfg(test)]
    mod test {
        use super::{parse_one_packet, Value};

        #[test]
        fn creation() {}

        #[test]
        fn empty_lists() {
            // Construct
            let l = Value::List(Vec::from([Value::List(Vec::from([Value::List(
                Vec::from([]),
            )]))]));
            let r = Value::List(Vec::from([Value::List(Vec::from([]))]));

            assert!(l > r);
        }

        #[test]
        fn non_empty_lists() {
            // Construct
            let l = Value::List(Vec::from([
                Value::Int(7),
                Value::Int(7),
                Value::Int(7),
                Value::Int(7),
            ]));
            let r = Value::List(Vec::from([Value::Int(7), Value::Int(7), Value::Int(7)]));

            assert!(l > r);
        }

        #[test]
        fn list_mismatch() {
            // Construct
            let l = Value::List(Vec::from([
                Value::List(Vec::from([Value::Int(1)])),
                Value::List(Vec::from([Value::Int(2), Value::Int(3), Value::Int(4)])),
            ]));
            let r = Value::List(Vec::from([
                Value::List(Vec::from([Value::Int(1)])),
                Value::Int(4),
            ]));

            assert!(l < r);
        }

        #[test]
        fn parsing_1() {
            assert_eq!(
                Value::List(Vec::from([
                    Value::List(Vec::from([Value::Int(1)])),
                    Value::List(Vec::from([Value::Int(2), Value::Int(3), Value::Int(4)])),
                ])),
                parse_one_packet("[[1],[2,3,4]]")
            );
        }
    }
}
