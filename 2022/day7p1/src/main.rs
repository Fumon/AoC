#![feature(iter_intersperse)]

use filesystem::{FileSystem, Path, Size};

use crate::filewalk::FileWalk;

fn main() -> std::io::Result<()> {
    let mut filesystem = FileSystem::from(FileWalk::from_file("./input"));
    filesystem.set_capacity(70_000_000);

    let combined_size = part_1(&mut filesystem);

    println!("Combined size of candidate directories: {}", combined_size);

    let directory_for_deletion = part_2(&mut filesystem, 30000000);

    println!("Directory for deltion {:?}", directory_for_deletion);

    Ok(())
}

fn part_1(filesystem: &mut FileSystem) -> u32 {
    // Sort and tally
    filesystem
        .get_dir_sizes()
        .unwrap()
        .iter()
        .filter_map(|(_, size)| if *size <= 100_000 { Some(size) } else { None })
        .sum()
}

fn part_2(filesystem: &mut FileSystem, needed: Size) -> (Path, Size) {
    let target_to_free = compute_target_to_free(needed, filesystem.get_available_space());

    println!("Target to free: {}", target_to_free);

    let dirs = filesystem.get_dir_sizes().unwrap();
    let mut candidates = dirs
        .iter()
        .filter_map(|(path, size)| {
            if *size >= target_to_free {
                Some((path.clone(), size.clone()))
            } else {
                None
            }
        })
        .collect::<Vec<(Path, Size)>>();

    candidates.sort_unstable_by_key(|dir| dir.1);

    for ele in candidates.iter().enumerate() {
        println!("Candidate {: >3}: {:?}", ele.0, ele.1);
    }

    candidates.iter().nth(0).cloned().unwrap()
}

fn compute_target_to_free(needed: u32, available: u32) -> u32 {
    needed - available
}

#[cfg(test)]
mod test {

    use crate::{
        compute_target_to_free,
        filesystem::{FileSystem, Path, Size},
        filewalk::FileWalk,
        part_1, part_2,
    };

    const FILE_SYSTEM_CAPACITY: Size = 70_000_000;
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
        let mut filesystem = load_example_filesystem();

        let rootsize = filesystem.get_dir_size(Path::default()).unwrap();

        assert_eq!(rootsize, 48381165);

        assert_eq!(part_1(&mut filesystem), 95437);
    }

    #[test]
    fn part_2_test() {
        let mut filesystem = load_example_filesystem();

        let needed = 30000000;
        let available = filesystem.get_available_space();

        assert_eq!(available, 21618835);

        let target_to_free = compute_target_to_free(needed, available);
        assert_eq!(target_to_free, 8381165);

        let candidate = part_2(&mut filesystem, needed);

        assert_eq!(candidate.0.to_string(), "d".to_owned());
        assert_eq!(candidate.1, 24933642)
    }

    fn load_example_filesystem() -> FileSystem {
        let fw = FileWalk::from_string(EXAMPLE.into());
        let mut fs = FileSystem::from(fw);
        fs.set_capacity(FILE_SYSTEM_CAPACITY);
        fs
    }
}

mod filesystem {
    use std::{
        collections::{HashMap, HashSet},
        iter::once,
        ops::Add,
    };

    use crate::filewalk::{ChangeDirArg, FileWalk, Listing, Step};

    #[derive(Debug)]
    pub(crate) struct FileSystemError {
        kind: ErrorType,
    }

    #[derive(Debug)]
    pub(crate) enum ErrorType {
        DoesNotExist,
        AlreadyExists,
        ParentDirectoryDoesNotExist,
    }

    pub(crate) struct FileSystem {
        capacity: Size,
        dir_table: HashMap<Path, HashSet<FileTableEntry>>,
        dir_sizes: HashMap<Path, Size>,
        size_scuffed: bool,
        used_space: Size,
    }

    impl Default for FileSystem {
        fn default() -> Self {
            // Root must exist
            let mut starting_table: HashMap<Path, HashSet<FileTableEntry>> = Default::default();
            starting_table.insert(Path::default(), Default::default());

            Self {
                capacity: Default::default(),
                dir_table: starting_table,
                dir_sizes: Default::default(),
                size_scuffed: true,
                used_space: Default::default(),
            }
        }
    }

    impl FileSystem {
        pub(crate) fn new(capacity: Size) -> Self {
            let mut me = Self::default();
            me.capacity = capacity;
            me
        }

        pub(crate) fn set_capacity(&mut self, cap: Size) {
            self.capacity = cap;
        }

        pub(crate) fn insert_at_path(
            &mut self,
            path: Path,
            entry: FileTableEntry,
        ) -> Result<(), FileSystemError> {
            // Find parent directory
            let Some(parent_dir) = self.dir_table.get_mut(&path) else {
                return Err(FileSystemError { kind: ErrorType::ParentDirectoryDoesNotExist });
            };

            // Insert entry into parent's child listing
            if !parent_dir.insert(entry.clone()) {
                return Err(FileSystemError {
                    kind: ErrorType::AlreadyExists,
                });
            }
            self.size_scuffed = true;

            if let FileTableEntry::Directory(name) = entry {
                // Insert a new directory entry into the table
                let new_path = &path + name;
                let None = self.dir_table.insert(new_path, Default::default()) else {
                    return Err(FileSystemError { kind: ErrorType::AlreadyExists })
                };
            }

            Ok(())
        }

        fn get_dir_contents(&self, path: Path) -> Option<&HashSet<FileTableEntry>> {
            self.dir_table.get(&path)
        }

        pub(crate) fn get_dir_size(&mut self, path: Path) -> Result<Size, FileSystemError> {
            if let Some(size) = self.dir_sizes.get(&path) {
                return Ok(size.clone());
            } else {
                // Otherwise get contents and recurse on directories
                let mut new_size = 0;

                // Get contents
                let contents: Vec<FileTableEntry> = {
                    if let Some(contents_cpy) = self.get_dir_contents(path.clone()) {
                        contents_cpy.iter().cloned().collect()
                    } else {
                        return Err(FileSystemError {
                            kind: ErrorType::DoesNotExist,
                        });
                    }
                };

                // Iterate
                for entry in contents.iter() {
                    new_size += match entry {
                        FileTableEntry::Directory(childname) => {
                            self.get_dir_size(&path + childname.clone())?
                        }
                        FileTableEntry::File(_, filesize) => filesize.clone(),
                    }
                }

                // Set size
                self.dir_sizes.insert(path.clone(), new_size.clone());

                return Ok(new_size);
            }
        }

        pub(crate) fn get_dir_sizes(&mut self) -> Result<HashMap<Path, Size>, FileSystemError> {
            if self.dir_sizes.len() == 0 || self.size_scuffed {
                self.get_dir_size(Default::default())?;
                self.size_scuffed = false;
            }

            Ok(self.dir_sizes.clone())
        }

        pub(crate) fn get_available_space(&mut self) -> Size {
            self.capacity - self.get_dir_size(Default::default()).unwrap()
        }
    }

    impl From<FileWalk> for FileSystem {
        fn from(fw: FileWalk) -> Self {
            let mut p = Path::default();
            let mut fs = FileSystem::default();

            for step in fw {
                match step {
                    Step::ChangeDir(ChangeDirArg::Root) => continue,
                    Step::ChangeDir(ChangeDirArg::Parent) => p = p.pop().unwrap(),
                    Step::ChangeDir(ChangeDirArg::Child(dirname)) => p = p.push(dirname),
                    Step::DirListing(dlist) => dlist
                        .into_iter()
                        .map(|l| FileTableEntry::from(l))
                        .for_each(|entry| fs.insert_at_path(p.clone(), entry).unwrap()),
                }
            }
            fs
        }
    }

    pub(crate) type Name = String;
    pub(crate) type Size = u32;

    #[derive(Debug, Hash, PartialEq, Eq, Clone)]
    pub(crate) enum FileTableEntry {
        Directory(Name),
        File(Name, Size),
    }

    impl From<Listing> for FileTableEntry {
        fn from(value: Listing) -> Self {
            match value {
                Listing::Directory(name) => Self::Directory(name),
                Listing::File(size, name) => Self::File(name, size),
            }
        }
    }

    // Directories will have an associated list of contents
    //     Store this internally? Or store a reference to the list?
    // Files have sizes
    // Directories have sizes comprised of the files within and all directories within

    /// A zero length components vec is equivalent to the root directory
    #[derive(Debug, Hash, PartialEq, Eq, Default, Clone)]
    pub(crate) struct Path {
        components: Vec<String>,
    }

    impl Path {
        fn push(&self, c: String) -> Path {
            Path {
                components: self.components.iter().cloned().chain(once(c)).collect(),
            }
        }
        fn pop(&self) -> Option<Path> {
            if self.components.len() == 0 {
                None
            } else {
                Some(Path {
                    components: self
                        .components
                        .iter()
                        .cloned()
                        .take(self.components.len() - 1)
                        .collect(),
                })
            }
        }
    }

    impl Add<String> for &Path {
        type Output = Path;

        fn add(self, rhs: String) -> Self::Output {
            self.push(rhs)
        }
    }

    impl ToString for Path {
        fn to_string(&self) -> String {
            self.components
                .iter()
                .cloned()
                .intersperse("/".to_owned())
                .collect()
        }
    }
}

mod filewalk {
    use std::{
        fs::File,
        io::{BufRead, BufReader, Cursor},
        iter::from_fn,
        path::Path,
        str::FromStr,
    };

    pub(crate) struct FileWalk(Vec<Step>);

    impl FileWalk {
        fn load_filewalk(input: &mut std::iter::Peekable<impl Iterator<Item = String>>) -> Self {
            FileWalk(
                from_fn(|| {
                    let Some(command_line) = input.next() else {
                    return None;
                };

                    let mut c = command_line.split_ascii_whitespace().skip(1);

                    match (c.next(), c.next()) {
                        (Some("cd"), Some(d)) => {
                            Some(Step::ChangeDir(ChangeDirArg::from_str(d).unwrap()))
                        }
                        (Some("ls"), None) => {
                            let listing_lines = from_fn(|| input.next_if(|v| !v.starts_with("$")));

                            Some(Step::DirListing(DList::from_iter(listing_lines)))
                        }
                        (Some(z), _) => panic!("Invalid command {}", z),
                        (None, _) => panic!("No command in command line"),
                    }
                })
                .collect(),
            )
        }

        pub(crate) fn from_file<P: AsRef<Path>>(file: P) -> Self {
            Self::load_filewalk(
                &mut BufReader::new(File::open(file).unwrap())
                    .lines()
                    .flatten()
                    .peekable(),
            )
        }

        pub(crate) fn from_string(input: String) -> Self {
            let mut string_iter = Cursor::new(input).lines().flatten().peekable();
            Self::load_filewalk(&mut string_iter)
        }
    }

    impl IntoIterator for FileWalk {
        type Item = Step;
        type IntoIter = std::vec::IntoIter<Self::Item>;

        fn into_iter(self) -> Self::IntoIter {
            self.0.into_iter()
        }
    }

    pub(crate) enum Step {
        ChangeDir(ChangeDirArg),
        DirListing(DList),
    }

    pub(crate) enum ChangeDirArg {
        Root,
        Parent,
        Child(String),
    }

    impl FromStr for ChangeDirArg {
        type Err = &'static str;

        fn from_str(s: &str) -> Result<Self, Self::Err> {
            Ok(match s {
                "/" => Self::Root,
                ".." => Self::Parent,
                "" => return Err("Empty Change Directory Argument"),
                _ => Self::Child(s.to_owned()),
            })
        }
    }

    pub(crate) enum Listing {
        Directory(String),
        File(u32, String),
    }

    pub(crate) struct DList(Vec<Listing>);

    impl FromIterator<String> for DList {
        fn from_iter<T: IntoIterator<Item = String>>(input: T) -> Self {
            DList(
                input
                    .into_iter()
                    .map(|x| {
                        let Some((a, b)) = x.split_once(" ") else {
                            panic!("There are only two elements in a listing")
                        };

                        match a {
                            "dir" => Listing::Directory(b.to_owned()),
                            _ => Listing::File(a.parse().unwrap(), b.to_owned()),
                        }
                    })
                    .collect(),
            )
        }
    }

    impl IntoIterator for DList {
        type Item = Listing;
        type IntoIter = std::vec::IntoIter<Self::Item>;

        fn into_iter(self) -> Self::IntoIter {
            self.0.into_iter()
        }
    }
}
