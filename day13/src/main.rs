#![feature(array_chunks)]
use std::{collections::BinaryHeap, fs::read_to_string, iter::zip};

use packets::{parse_full_packet_input, Value};

fn main() {
    let ppairs = parse_full_packet_input(read_to_string("./input").unwrap().as_str());

    println!(
        "The sum of indicies of in-order pairs of packets is:\n\t{}",
        part_1(&ppairs)
    );
    
    println!(
        "The product of the indicies for the two separator packets is:\n\t{}",
        part_2(&ppairs)
    );
}

fn part_1(ppairs: &Vec<Value>) -> usize {
    zip(1.., ppairs.array_chunks())
        .filter_map(|(i, [l, r])| if l < r { Some(i) } else { None })
        .sum()
}

fn part_2(ppairs: &Vec<Value>) -> usize {
    let dividers = parse_full_packet_input("[[2]]\n[[6]]");

    let mut heap = BinaryHeap::new();
    heap.extend(ppairs.iter());
    heap.extend(dividers.iter());

    zip(1.., heap.into_sorted_vec().into_iter())
        .filter_map(|(i, packet)| {
            if dividers.contains(packet) {
                Some(i)
            } else {
                None
            }
        })
        .product()
}

#[cfg(test)]
mod test {
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

        assert_eq!(13, super::part_1(&ppairs));
    }

    #[test]
    fn p2_example() {
        let ppairs = parse_full_packet_input(EXAMPLE_1);
        assert_eq!(140, super::part_2(&ppairs));
    }
}

mod packets {
    use std::cmp::Ordering;

    use nom::{
        branch::alt,
        character::complete::{char, digit1, line_ending},
        combinator::{all_consuming, map_res},
        multi::{many1, separated_list0, many0},
        sequence::{delimited, terminated},
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

    #[allow(unused)]
    fn parse_one_packet(line: &str) -> Value {
        let ("", val) = all_consuming(parse_packet)(line).unwrap() else {
            panic!()
        };
        val
    }

    pub(crate) fn parse_full_packet_input(input: &str) -> Vec<Value> {
        let ("", ppairs) = all_consuming(many1(terminated(parse_packet, many0(line_ending))))(input).unwrap() else {
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
