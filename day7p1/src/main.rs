use core::panic;
use std::{
    fs::File,
    io::{BufRead, BufReader},
    path::PathBuf,
};

use support::{CDArg, Command, Directory, FileSystem, Listing};

fn main() -> std::io::Result<()> {
    let mut line_iter = BufReader::new(File::open("./input")?).lines().flatten();

    let combined_size = part_1(&mut line_iter);

    println!("Combined size of candidate directories: {}", combined_size);

    Ok(())
}

fn part_1(line_iter: &mut impl Iterator<Item = String>) -> u32 {
    let mut p_line_iter = line_iter.peekable();

    let mut dirs: FileSystem = Default::default();
    let mut current_path: Option<PathBuf> = None;
    let mut current_dir_size: u32 = 0;
    loop {
        // Parse Command
        let Some(cline) = p_line_iter.next() else {
            // Add to parents
            while current_path.as_mut().map(|x| x.pop()).unwrap() {
                dirs.add_to_size_of_directory(current_path.clone().unwrap(), current_dir_size);
            }

            break;
        };
        let command = match cline.parse::<Command>() {
            Ok(command) => command,
            Err(e) => panic!("Error parsing command: {}", e),
        };

        match command {
            Command::CD(CDArg::Parent) => {
                // Pop up
                // Add to the directory's total
                current_path.as_mut()
                    .map(|x| {
                        x.pop()
                    })
                    .expect("Only pop if not beyond root");

                // Update sizes
                current_dir_size =
                    dirs.add_to_size_of_directory(current_path.clone().unwrap(), current_dir_size);
            }
            Command::CD(CDArg::Name(name)) => {
                // Descend
                let p = current_path.clone();
                current_path
                    .get_or_insert(PathBuf::new())
                    .push(name.clone());

                dirs.insert(Directory {
                    parent: p,
                    name,
                    size: 0,
                });

                current_dir_size = 0;
            }
            Command::LS => {
                // Consume until next command
                while let Some(ls_line) =
                    p_line_iter.next_if(|line| line.chars().nth(0) != Some('$'))
                {
                    let listing = match ls_line.parse::<Listing>() {
                        Ok(listing) => listing,
                        Err(e) => panic!("Error parsing listing: {}", e),
                    };

                    if let Listing::File(size, _) = listing {
                        current_dir_size += size;
                    }
                }

                // Update this directory
                dirs.add_to_size_of_directory(current_path.clone().unwrap(), current_dir_size);
            }
        };
    }

    // Sort and tally
    dirs.get_dirs().filter_map(|dir| if dir.size <= 100000 { Some(dir.size)} else {None}).sum()
}

#[cfg(test)]
mod test {}

mod support {
    use std::{collections::HashMap, path::PathBuf, str::FromStr};

    #[derive(Debug, Default)]
    pub(crate) struct FileSystem {
        file_table: HashMap<PathBuf, Directory>,
    }

    impl FileSystem {
        /// Looks up and updates the size of a directory returning the directory's new size
        pub(crate) fn add_to_size_of_directory(
            &mut self,
            path: PathBuf,
            size_increase: u32,
        ) -> u32 {
            let mut newsize: u32 = 0;
            self.file_table.entry(path).and_modify(|dir| {
                newsize = size_increase + dir.size;
                dir.size = newsize;
            });

            newsize
        }
        pub(crate) fn insert(&mut self, d: Directory) {
            let mut binding = d.parent.clone();
            let newpath = binding.get_or_insert(PathBuf::new());
            newpath.push(d.name.clone());

            self.file_table.insert(newpath.to_owned(), d);
        }
        pub(crate) fn get_dirs(self) -> impl Iterator<Item = Directory> {
            self.file_table.into_iter().map(|(_, dir)| dir)
        }
    }

    #[derive(Debug)]
    pub(crate) struct Directory {
        pub(crate) parent: Option<PathBuf>,
        pub(crate) name: String,
        pub(crate) size: u32,
    }

    // pub(crate) struct File {
    //     pub(crate) parent: PathBuf,
    //     pub(crate) name: String,
    //     pub(crate) size: u32,
    // }

    pub(crate) enum Command {
        LS,
        CD(CDArg),
    }

    impl FromStr for Command {
        type Err = &'static str;

        fn from_str(s: &str) -> Result<Self, Self::Err> {
            let mut b = s.split_ascii_whitespace();
            let Some(s1) = b.next() else {
                return Err("No input to parse");
            };
            if s1 != "$" {
                return Err("Tried to parse a command which didn't start with $");
            }

            match (b.next(), b.next()) {
                (Some("ls"), None) => Ok(Self::LS),
                (Some("cd"), Some("..")) => Ok(Self::CD(CDArg::Parent)),
                (Some("cd"), Some(path)) => Ok(Self::CD(CDArg::Name(path.to_owned()))),
                (_, _) => Err("invalid command"),
            }
        }
    }

    pub(crate) enum CDArg {
        Name(String),
        Parent,
    }

    pub(crate) enum Listing {
        File(u32, String),
        Dir(String),
    }

    impl FromStr for Listing {
        type Err = String;

        fn from_str(s: &str) -> Result<Self, Self::Err> {
            let mut b = s.split_ascii_whitespace();

            match (b.next(), b.next()) {
                (Some("dir"), Some(name)) => Ok(Self::Dir(name.to_owned())),
                (Some(size_str), Some(name)) => {
                    let size: u32 = size_str.parse().map_err(|x| format!("{}", x))?;
                    Ok(Self::File(size, name.to_owned()))
                }
                (_, _) => Err("invalid listing".to_string()),
            }
        }
    }
}
