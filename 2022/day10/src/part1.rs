use crate::cpu::{Op, CPU};

pub(crate) const SIGNAL_STRENGTH_SAMPLE_CYCLES: [usize; 6] = [20, 60, 100, 140, 180, 220];

pub(crate) fn part_1(instructions: Vec<Op>) {
    let mut cpu = CPU::new(instructions);
    
    let signal_sum: i64 = SIGNAL_STRENGTH_SAMPLE_CYCLES.iter().scan(0 , |cycle, target_cycle| {
        while *cycle < *target_cycle {
            if cpu.cycle() == None {
                panic!("CPU died while iterating");
            }
            *cycle += 1;
        }
        let x = cpu.program_memory.get(&crate::cpu::Register::X).expect("retrieval of X is always possible");
        let signal_strength = signal_strength(&(*cycle).try_into().unwrap(), x);

        Some(signal_strength)
    }).sum();

    println!("The signal sum is {}", signal_sum);
}

pub(crate) fn signal_strength(cycle: &u32, xval: &i32) -> i64 {
    (*cycle as i64) * (*xval as i64)
}
