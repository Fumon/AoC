
from itertools import combinations

def find_combos(groups, size):
    glen = len(groups)
    blanks = size - sum(groups) - (glen - 1)
    return combinations(range(glen + blanks), glen)

with open("testinput") as f:
    lines = f.readlines()

for line in lines:
    sp = line.split(" ")
    board = '?'.join([sp[0]]*5)
    groups = [int(x) for x in sp[1].split(",")]*5
    
    combos = find_combos(groups, len(board))
    count_of_combos = 0
    for x in combos:
        count_of_combos += 1

    print(f"{line.rstrip()} - {count_of_combos}")
    print(board)
    print(groups)
    # print(f"\t{combos}")
    print()