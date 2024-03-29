package main

import (
	"fmt"
	"fuaoc2023/day11/u"
	"math"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))
}

type Galaxy struct {
	position [2]int
}

func part1(lines <-chan string) int {
	// Adjust for columns
	// Shortest Paths
	sum_shortest_path := sum_shortest_paths_with_factor(lines, 2)

	return sum_shortest_path
}

func part2(lines <-chan string) int {
	// Adjust for columns
	// Shortest Paths
	sum_shortest_path := sum_shortest_paths_with_factor(lines, 1_000_000)

	return sum_shortest_path
}

func sum_shortest_paths_with_factor(lines <-chan string, growthfactor int) int {
	addmul := (growthfactor - 1)
	lines_2, line_len := PeekAtFirstLineWidth(lines)

	var galaxies []*Galaxy
	var galaxy_columns = make(map[int][]*Galaxy)

	var blank_rows []int
	var column_counts = make([]int, line_len)
	{
		var linecount = 0
		for line := range lines_2 {
			var line_galaxy_count = 0
			for i, c := range []byte(line) {
				if c == '#' {
					line_galaxy_count++
					column_counts[i]++
					g := &Galaxy{[2]int{i, linecount + (addmul * len(blank_rows))}}
					galaxies = append(galaxies, g)
					galaxy_columns[i] = append(galaxy_columns[i], g)
				}
			}

			if line_galaxy_count == 0 {
				blank_rows = append(blank_rows, linecount)
			}
			linecount++
		}
	}

	{
		var blank_column_count int
		for i, count := range column_counts {
			if count == 0 {
				blank_column_count++
			} else {
				for _, gal := range galaxy_columns[i] {
					gal.position[0] += addmul * blank_column_count
				}
			}
		}
	}

	var sum_shortest_path int
	{
		var pathcount = (len(galaxies) * (len(galaxies) - 1)) / 2
		sumch := make(chan int, len(galaxies))
		for i := len(galaxies) - 1; i > 0; i-- {
			root := galaxies[i]
			leaves := galaxies[0:i]
			go func() {
				x1, y1 := root.position[0], root.position[1]
				for _, l := range leaves {
					sumch <- int(math.Abs(float64(l.position[0]-x1)) + math.Abs(float64(l.position[1]-y1)))
				}
			}()
		}

		for i := pathcount; i > 0; i-- {
			sum_shortest_path += <-sumch
		}
		close(sumch)
	}
	return sum_shortest_path
}

func PeekAtFirstLineWidth(lines <-chan string) (<-chan string, int) {
	line1 := <-lines
	nc := make(chan string, 3)
	nc <- line1
	go func() {
		for line := range lines {
			nc <- line
		}
		close(nc)
	}()

	return nc, len(line1)
}
