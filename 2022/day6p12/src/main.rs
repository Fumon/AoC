use std::{collections::HashSet, fs::read_to_string, io};

fn main() -> io::Result<()> {
    let first_start_packet_complete_index = find_first_run_of_unique_chars(read_to_string("./input")?,14) + 14;


    println!(
        "The first start of packet marker finishes on the {}th character",
        first_start_packet_complete_index
    );
    Ok(())
}

fn find_first_run_of_unique_chars<T: AsRef<str>>(haystack: T, n: usize) -> usize {
    haystack.as_ref().as_bytes()
    .windows(n)
    .enumerate()
    .map(|(ind, wind)| (ind, wind)) // Offset for puzzle
    .find_map(|(index, window)| {
        let mut m: HashSet<&u8> = HashSet::default();
        if window.iter().map(|i| m.insert(i)).any(|x| !x) {
            None
        } else {
            Some(index)
        }
    })
    .expect("There will be an answer")
}