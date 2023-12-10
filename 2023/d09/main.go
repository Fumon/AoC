package main

import (
	"fmt"
	"fuaoc2023/day09/u"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))
}

func predict_diff(future bool, lines <-chan string) int {

	var sum_of_next int
	for line := range lines {
		nums := u.ParseNums(strings.Split(line, " "))
		var difflines [][]int

		fmt.Print(nums, ": ")
		difflines = append(difflines, nums)
		for i := 0; ; i++ {
			var newdiffline []int
			var nonzeros int
			for j := 1; j < len(difflines[i]); j++ {
				val := difflines[i][j] - difflines[i][j-1]
				newdiffline = append(newdiffline, val)
				if val != 0 {
					nonzeros += 1
				}
			}
			difflines = append(difflines, newdiffline)
			if nonzeros == 0 {
				break
			} else {
				fmt.Print(" ", nonzeros)
			}
		}

		fmt.Print(" -> ")

		// Predict
		endlineindex := len(difflines) - 1
		var addindex int
		var indexdelta int
		if future {
			addindex = len(difflines[endlineindex]) - 1
			indexdelta = 1
		}

		fmt.Print("\n")
		for _, dl := range difflines {
			for _, v := range dl {
				fmt.Print("\t", v)
			}
			fmt.Print("\n")
		}

		var nextval int
		for i := endlineindex; i > -1; i, addindex = i-1, addindex+indexdelta {
			if future {
				nextval += difflines[i][addindex]
			} else {
				nextval = difflines[i][addindex] - nextval
			}
		}
		sum_of_next += nextval
		fmt.Print(nextval, "\n")
	}
	return sum_of_next
}
