use std::ops::{Rem, Div};




pub(crate) struct CathodeRayTube {
    width: u8,
    height: u8,
}

impl CathodeRayTube {
    pub(crate) fn new(width: u8, height: u8) -> Self {
        if width.checked_mul(height) == None {
            panic!("Too many pixels!");
        }
        CathodeRayTube { width, height }
    }
    pub(crate) fn get_beam(&self, cycle: usize) -> (u8, u8) {
        let x = cycle.rem(self.width as usize);
        let y = cycle.div(self.width as usize).rem(self.width as usize * self.height as usize);

        (x.try_into().unwrap(), y.try_into().unwrap())
    }
}

pub(crate) struct Sprite {
    width: u8,
    range: u8,
}

impl Sprite {
    pub(crate) fn new(width: u8) -> Self {
        if width.rem(2) == 0 {
            panic!("Sprite cannot have even width");
        }

        Sprite {
            width,
            range: width/2
        }
    }
    pub(crate) fn hit_scan(&self, pos: u8, test: u8) -> bool{
        pos.abs_diff(test) <= self.range
    }
}