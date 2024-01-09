package main

import (
	"fmt"
	"fuaoc2023/day21/u"
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
	curmap[elfstart] = struct{}{}
	for step := 1; step <= steps; step++ {
		for p := range curmap {
			for _, pa := range p.around() {
				if playboard[((pa[0] % linecount) + linecount) % linecount][((pa[1] % width) + width) % width] {
					nextmap[pa] = struct{}{}
				}
			}
		}
		curmap, nextmap = nextmap, make(map[point]struct{}, len(nextmap))
	}

	return len(curmap)
}
