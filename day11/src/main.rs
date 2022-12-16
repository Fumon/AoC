use std::{
    fs::read_to_string,
    iter::{from_fn, repeat},
    ops::Rem,
};

use monkey::{Monkey, Worry, WorryOperation};
use nom::{character::complete::line_ending, combinator::all_consuming, multi::separated_list1};

mod monkey;

fn main() {
    let monkeys_and_items = get_monkeys_from_input(read_to_string("./input").unwrap());
    //part_1(monkeys_and_items);
    part_2(monkeys_and_items);
}

fn part_1(monkeys_and_items: Vec<(Vec<Worry>, Monkey)>) {
    let mut mb = MonkeyBusiness::new(monkeys_and_items, false);

    let total_monkey_business: usize = mb.business_after_n_rounds(20);

    println!(
        "Total Monkey Business after 20 rounds: {}",
        total_monkey_business
    )
}

fn part_2(monkeys_and_items: Vec<(Vec<Worry>, Monkey)>) {
    let mut mb = MonkeyBusiness::new(monkeys_and_items, true);

    let total_monkey_business: usize = mb.business_after_n_rounds(10_000);

    println!(
        "Total Monkey Business after 10_000 rounds: {}",
        total_monkey_business
    );
}

struct MonkeyBusiness {
    monkeys: Vec<Monkey>,
    item_lists: Vec<Vec<Worry>>,
    mod_product: Worry,
    part2: bool,
}

impl MonkeyBusiness {
    fn new(input: Vec<(Vec<Worry>, Monkey)>, part2: bool) -> Self {
        let (item_lists, monkeys): (Vec<Vec<Worry>>, Vec<Monkey>) = input.into_iter().unzip();
        let mod_product = monkeys
            .iter()
            .map(|m| m.test.modulo.clone())
            .product::<u64>();
        MonkeyBusiness {
            monkeys,
            item_lists,
            mod_product,
            part2,
        }
    }

    fn run_round(&mut self) -> Vec<usize> {
        let (monkeys, item_lists) = (&mut self.monkeys, &mut self.item_lists);
        monkeys
            .iter()
            .map(|monkey| {
                let items: Vec<Worry> = item_lists
                    .get_mut(monkey.index as usize)
                    .unwrap()
                    .drain(..)
                    .collect();
                let inspect_len = items.len();

                let op_iter = repeat(&monkey.worry_op);
                let item_iter = items.into_iter().map(|w| {
                    if !self.part2 {
                        w
                    } else {
                        w.rem(self.mod_product)
                    }
                });

                let before_reduce_iter = op_iter
                    .zip(item_iter)
                    .map(WorryOperation::exec_tup)
                    .map(|v| if self.part2 {v} else {v.saturating_div(3)});

                let (true_vec, false_vec): (Vec<Worry>, Vec<Worry>) = before_reduce_iter
                    .partition(|w| w.rem(monkey.test.modulo) == 0);

                item_lists
                    .get_mut(monkey.test.true_dst as usize)
                    .unwrap()
                    .extend(true_vec);
                item_lists
                    .get_mut(monkey.test.false_dst as usize)
                    .unwrap()
                    .extend(false_vec);

                inspect_len
            })
            .collect()
    }

    fn print_items(&self) {
        for (i, list) in self.item_lists.iter().enumerate() {
            println!("Monkey {}: {:?}", i, list);
        }
    }

    fn business_after_n_rounds(&mut self, rounds: usize) -> usize {
        let mut business_totals = from_fn(|| {
            let run = self.run_round();
            Some(run)
        })
        .take(rounds)
        .reduce(|monkey_business, round_business| {
            monkey_business
                .into_iter()
                .zip(round_business.into_iter())
                .map(|(a, b)| a + b)
                .collect()
        })
        .unwrap();
        business_totals.sort();
        business_totals.iter().rev().take(2).product()
    }
}

fn get_monkeys_from_input<T: AsRef<str>>(input: T) -> Vec<(Vec<Worry>, Monkey)> {
    let (_, monkeys) =
        all_consuming(separated_list1(line_ending, Monkey::parse_monkey))(input.as_ref()).unwrap();

    return monkeys;
}

#[cfg(test)]
mod test {
    use crate::{get_monkeys_from_input, monkey::Worry, MonkeyBusiness};

    const SAMPLE_MONKEYS: &'static str = r#"Monkey 0:
  Starting items: 79, 98
  Operation: new = old * 19
  Test: divisible by 23
    If true: throw to monkey 2
    If false: throw to monkey 3

Monkey 1:
  Starting items: 54, 65, 75, 74
  Operation: new = old + 6
  Test: divisible by 19
    If true: throw to monkey 2
    If false: throw to monkey 0

Monkey 2:
  Starting items: 79, 60, 97
  Operation: new = old * old
  Test: divisible by 13
    If true: throw to monkey 1
    If false: throw to monkey 3

Monkey 3:
  Starting items: 74
  Operation: new = old + 3
  Test: divisible by 17
    If true: throw to monkey 0
    If false: throw to monkey 1"#;

    #[test]
    fn one_round() {
        let round1_results: Vec<Vec<Worry>> = Vec::from([
            Vec::from([20, 23, 27, 26]),
            Vec::from([2080, 25, 167, 207, 401, 1046]),
            Vec::from([]),
            Vec::from([]),
        ]);

        let mut mb = MonkeyBusiness::new(get_monkeys_from_input(SAMPLE_MONKEYS), false);

        let monkey_business_round = mb.run_round();

        assert_eq!(monkey_business_round, vec![2, 4, 3, 5]);

        assert_eq!(mb.item_lists, round1_results);
    }

    #[test]
    fn total_after_20() {
        let mut mb = MonkeyBusiness::new(get_monkeys_from_input(SAMPLE_MONKEYS), false);

        assert_eq!(mb.business_after_n_rounds(20), 10605);
    }

    #[test]
    fn total_after_10000() {
        let mut mb = MonkeyBusiness::new(get_monkeys_from_input(SAMPLE_MONKEYS), true);

        assert_eq!(mb.business_after_n_rounds(10_000), 2713310158);
    }
}
