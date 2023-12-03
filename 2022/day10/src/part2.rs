use crate::{cpu::{Op, CPU}, crt::{Sprite, CathodeRayTube}};

const SPRITE_WIDTH: u8 = 3;
const TUBE_WIDTH: u8 = 40;
const TUBE_WIDTH_U: usize = 40;
const TUBE_HEIGHT: u8 = 6;

pub(crate) fn part_2(instructions: Vec<Op>) {
    let mut cpu = CPU::new(instructions);

    let sprite = Sprite::new(SPRITE_WIDTH);
    let crt = CathodeRayTube::new(TUBE_WIDTH, TUBE_HEIGHT);

    let tube = (0..(TUBE_WIDTH*TUBE_HEIGHT) as usize).map(|cycle| {
        if cpu.cycle() == None {
            panic!("CPU died while iterating");
        }

        let x = cpu.program_memory.get(&crate::cpu::Register::X).unwrap();
        let Ok(xpos) = (*x).try_into() else {
            return '.';
        };

        let (beamx, _beamy) = crt.get_beam(cycle);

        if sprite.hit_scan(xpos, beamx) {
            '#'
        } else {
            '.'
        }
    })
    .array_chunks::<TUBE_WIDTH_U>()
    .map(|line| format!("{}\n", String::from_iter(line.iter()))).collect::<String>();

    print!("{}", tube);
    ()
}