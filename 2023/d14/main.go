package main

import (
	"fmt"
	"fuaoc2023/day14/u"
)

func main() {
	fmt.Println("Part1: ", Part1(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	lines2, linewidth := u.PeekAtFirstLineWidth(lines)

	var columnstops = make([]int, linewidth)
	for i, _ := range columnstops {
		columnstops[i] = -1
	}
	var columnsum int
	var rock_count int
	var line_count int
	for line := range lines2 {
		for i, ch := range []byte(line) {
			switch ch {
			case '#':
				columnstops[i] = line_count
			case 'O':
				columnstops[i]++
				columnsum += columnstops[i]
				rock_count++
			}
		}
		line_count++
	}

	return (line_count * rock_count) - columnsum
}