

pub(crate) struct Monkey {
    items: Vec<u32>,
    worry_op: WorryOperation,
    test: MonkeyTest,
    activity: u32,
}

type WorryOperation = fn(u32) -> u32;

struct MonkeyTest {
    modulo: u32;
    true_dst: u32;
    false_dst: u32;
}