

with open("input") as f:
    lines = f.readlines()


# Check if a line does not contain a hash
count_of_no_hash = 0
for line in lines:
    sp = line.split(" ")
    board = '?'.join([sp[0]]*5)
    groups = [int(x) for x in sp[1].split(",")]*5

    if all([x != '#' for x in board]):
        count_of_no_hash += 1
        print(board)
        print(groups)
    


# Print the results of the checks as percentages of the total number of lines
print("Percentage of lines with no hash: {:.2f}%".format(count_of_no_hash / len(lines) * 100))
