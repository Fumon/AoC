package main

import (
	"fmt"
	"fuaoc2023/day04/u"
	"strings"
)

func main() {
	lines := u.Linewisefile_chan("input")

	fmt.Println("Total Points:\n", part1(lines))
}

func splitToNums(line string) [][]string {
	sections := strings.Split(strings.Split(line, ": ")[1], " | ")
	return [][]string{
		strings.Split(sections[0], " "),
		strings.Split(sections[1], " "),
	}
}

func part1(lines <-chan string) (pointtotal int) {
	for line := range lines {
		numsections := splitToNums(line)
		prizeset := make(map[string]struct{}, len(numsections[0]))
		for _, p := range numsections[0] {
			if len(p) > 0 {
				prizeset[strings.Trim(p, " ")]=struct{}{}
			}
		}

		var matchcount int
		for _, t := range numsections[1] {
			if _, match := prizeset[t]; match {
				matchcount++
			}
		}
		if matchcount > 0 {
			pointtotal += 1 << (matchcount - 1)
		}
	}
	return 
}