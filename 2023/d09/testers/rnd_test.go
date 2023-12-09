package testers

import (
	"fmt"
	"testing"
)

func TestDifferenceSequence(t *testing.T) {
	input := []int{-3, 10, 39, 84, 150, 259, 464, 863, 1611, 2928, 5101, 8478, 13452, 20433, 29806, 41873, 56777, 74406, 94275, 115384, 136050}

	var diffs [][]int
	diffs = append(diffs, input)
	print_diff_line(diffs[0])
	for i := 0; len(diffs[i]) > 1 ; i++ {
		var new_diffs []int
		for j := 1; j < len(diffs[i]); j += 1 {
			new_diffs = append(new_diffs, diffs[i][j] - diffs[i][j - 1])
		}
		print_diff_line(new_diffs)
		diffs = append(diffs, new_diffs)
	}
}

func print_diff_line(in []int) {
	for _, s := range in {
		fmt.Print("\t", s)
	}
	fmt.Print("\n")
}