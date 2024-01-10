package main

import (
	"fmt"
	"fuaoc2023/day21/u"
	"os"
)

func main() {
	fmt.Println(Part1(64, u.Linewisefile_chan("input")))
	fmt.Println(Part2(26501365, u.Linewisefile_chan("input")))
}

func Part1(steps int, lines <-chan string) int {
	lines2, width := u.PeekAtFirstLineWidth(lines)

	var elfstart point
	var playboard [][]bool
	var linecount int
	for line := range lines2 {
		var newline = make([]bool, width)
		for x, ch := range line {
			switch ch {
			case 'S':
				elfstart = point{linecount, x}
				fallthrough
			case '.':
				newline[x] = true
			}
		}
		playboard = append(playboard, newline)
		linecount++
	}

	var bounds = Bounds{
		width:  width,
		height: linecount,
	}
	var curmap = make(map[point]struct{})
	var nextmap = make(map[point]struct{})
	curmap[elfstart] = struct{}{}
	for step := 1; step <= steps; step++ {
		for p := range curmap {
			for _, pa := range p.around() {
				if bounds.check(pa) && playboard[pa[0]][pa[1]] {
					nextmap[pa] = struct{}{}
				}
			}
		}
		curmap, nextmap = nextmap, make(map[point]struct{})
	}

	return len(curmap)
}

func Part2(steps int, lines <-chan string) int {
	lines2, width := u.PeekAtFirstLineWidth(lines)

	var elfstart point
	var playboard [][]bool
	var linecount int
	for line := range lines2 {
		var newline = make([]bool, width)
		for x, ch := range line {
			switch ch {
			case 'S':
				elfstart = point{linecount, x}
				fallthrough
			case '.':
				newline[x] = true
			}
		}
		playboard = append(playboard, newline)
		linecount++
	}

	var curmap = make(map[point]struct{})
	var nextmap = make(map[point]struct{})
	var front []point
	curmap[elfstart] = struct{}{}
	front = append(front, elfstart)

	f, _ := os.Create("131steps.csv")
	// fmt.Fprintf(f, "%d,%d\n", 0, len(curmap))

	for step := 1; step <= steps; step++ {
		var nextfront []point
		
		for _, p := range front {
			for _, pa := range p.around() {
				if _, in := nextmap[pa]; !in {
					if playboard[((pa[0] % linecount) + linecount) % linecount][((pa[1] % width) + width) % width] {
						nextmap[pa] = struct{}{}
						nextfront = append(nextfront, pa)
					}
				}
			}
		}

		curmap, nextmap = nextmap, curmap
		front = nextfront
		tt := step - 65
		if tt >= 0 && tt % 131 == 0 {

			fmt.Fprintf(f, "%d,%d\n", tt/131, len(curmap))
		}
	}

	// one_ninety_six_count := len(curmap)
	// num_iterations := (steps - 196)/393
	// total := ((num_iterations)*2) - 1
	// total *= total
	// total *= one_ninety_six_count

	return len(curmap)
}
