package testers

import (
	"fmt"
	"testing"
)

func TestDifferenceSequence(t *testing.T) {
	input := []int{-1, -6, -15, -28, -42, -45, 11, 267, 1053, 3051, 7603, 17344, 37537, 78854, 162989, 333510, 675784, 1351492, 2654874, 5099428, 9547513}

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