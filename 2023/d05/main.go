package main

import (
	"fmt"
	"fuaoc2023/day05/u"
	"slices"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
}

func part1(lines <-chan string) int {

	seeds_strings := strings.Split(strings.Split((<-lines), ": ")[1], " ")
	cur_vals := u.ParseNums(seeds_strings)

	// Blank
	<- lines

	// Map Loop
	for {
		_, ok := <- lines
		if !ok {
			break
		}

		var new_vals []int
		for mapline := range lines {
			// Blank line, end of map
			if len(mapline) == 0 {
				break
			}
			
			// Parse
			ma := u.ParseNums(strings.Split(mapline, " "))
			dest_start := ma[0]
			source_start := ma[1]
			width := ma[2]
			smax := source_start + width
			diff := dest_start - source_start

			// Scan
			var ncurs []int
			for _, v := range cur_vals {
				if v >= source_start && v < smax {
					new_vals = append(new_vals, v + diff)
				} else {
					ncurs = append(ncurs, v)
				}
			}
			if len(ncurs) > 0 {
				cur_vals = ncurs
			} else {
				cur_vals = ncurs
				for len(<-lines) > 0 {}
				break
			}
		}
		// Move remaining values with their existing value
		if len(cur_vals) > 0 {
			new_vals = append(new_vals, cur_vals...)
		}
		cur_vals = new_vals
	}
	// Sort and pick smallest
	slices.Sort(cur_vals)
	return cur_vals[0]
}