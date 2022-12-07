use support::{FileSystem, Directory};

fn main() -> std::io::Result<()> {
    let filesystem = FileSystem::load_from_dir_walk_file("./input")?;

    let combined_size = part_1(&filesystem);

    println!("Combined size of candidate directories: {}", combined_size);

    let directory_for_deletion = part_2(&filesystem, 30000000);

    println!("Directory for deltion {:?}", directory_for_deletion);

    Ok(())
}

fn part_1(filesystem: &FileSystem) -> u32 {
    // Sort and tally
    filesystem
        .get_dirs()
        .iter()
        .filter_map(|dir| {
            if dir.size <= 100000 {
                Some(dir.size)
            } else {
                None
            }
        })
        .sum()
}

fn part_2(filesystem: &FileSystem, needed: u32) -> Directory {
    
    let target_to_free =compute_target_to_free(needed,filesystem.find_available_space());

    println!("Target to free: {}", target_to_free);

    let dirs = filesystem.get_dirs();
    let mut candidates = dirs.iter().filter_map(|dir| {
        if dir.size >= target_to_free {
            Some(dir)
        } else {
            None
        }
    }).collect::<Vec<&Directory>>();

    candidates.sort_unstable_by_key(|dir| dir.size);

    for ele in candidates.iter().enumerate() {
        println!("Candidate {: >3}: {:?}", ele.0, ele.1);
    }
    
    candidates.iter().nth(0).copied().unwrap().clone()
}

fn compute_target_to_free(needed: u32, available: u32) -> u32 {
    needed - available
}

#[cfg(test)]
mod test {
    use std::{io::{Cursor, BufRead}, path::PathBuf};

    use crate::{support::FileSystem, part_1, part_2, compute_target_to_free};

    const EXAMPLE: &'static str = r"$ cd /
$ ls
dir a
14848514 b.txt
8504156 c.dat
dir d
$ cd a
$ ls
dir e
29116 f
2557 g
62596 h.lst
$ cd e
$ ls
584 i
$ cd ..
$ cd ..
$ cd d
$ ls
4060174 j
8033020 d.log
5626152 d.ext
7214296 k";
    
    #[test]
    fn part_1_test() {
        let filesystem = load_example_filesystem();

        let root = filesystem.get_dir(PathBuf::from("/")).unwrap();

        assert_eq!(root.size, 48381165);
        
        assert_eq!(part_1(&filesystem), 95437);
    }

    #[test]
    fn part_2_test() {
        let filesystem = load_example_filesystem();

        let needed = 30000000;
        let available = filesystem.find_available_space();

        assert_eq!(available, 21618835);

        let target_to_free = compute_target_to_free(needed, available);
        assert_eq!(target_to_free, 8381165);

        let candidate = part_2(&filesystem, needed);

        assert_eq!(candidate.name, "d");
        assert_eq!(candidate.size, 24933642)
    }

    fn load_example_filesystem() -> FileSystem{
        FileSystem::load_from_dir_walk(Cursor::new(EXAMPLE).lines().flatten()).unwrap()
    }
}

mod support {
    use std::{
        collections::HashMap,
        fs::File,
        io::{BufRead, BufReader},
        path::{Path, PathBuf},
        str::FromStr,
    };

    #[derive(Debug, Default)]
    pub(crate) struct FileSystem {
        file_table: HashMap<PathBuf, Directory>,
    }

    impl FileSystem {
        pub(crate) const TOTAL_DISK_SPACE: u32 = 70000000;

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

        pub(crate) fn find_available_space(&self) -> u32{
            Self::TOTAL_DISK_SPACE - self.get_dir(PathBuf::from("/")).unwrap().size
        }

        pub(crate) fn get_dirs(&self) -> Vec<Directory> {
            self.file_table
                .iter()
                .map(|(_, dir)| dir)
                .cloned()
                .collect()
        }

        pub(crate) fn get_dir(&self, p: PathBuf) -> Option<&Directory> {
            self.file_table.get(&p)
        }

        pub(crate) fn load_from_dir_walk(i: impl Iterator<Item = String>) -> std::io::Result<Self> {
            let mut p_line_iter = i.peekable();

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
                        current_path
                            .as_mut()
                            .map(|x| x.pop())
                            .expect("Only pop if not beyond root");

                        // Update sizes
                        current_dir_size = dirs.add_to_size_of_directory(
                            current_path.clone().unwrap(),
                            current_dir_size,
                        );
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
                        dirs.add_to_size_of_directory(
                            current_path.clone().unwrap(),
                            current_dir_size,
                        );
                    }
                };
            }
            Ok(dirs)
        }

        pub(crate) fn load_from_dir_walk_file<P: AsRef<Path>>(
            file_name: P,
        ) -> std::io::Result<Self> {
            let line_iter = BufReader::new(File::open(file_name)?)
                .lines()
                .flatten();
            Self::load_from_dir_walk(line_iter)
        }
    }

    #[derive(Debug, Clone)]
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
