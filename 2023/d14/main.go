package main

import (
	"fmt"
	"fuaoc2023/day14/u"
)

func main() {
	fmt.Println("Part1: ", Part1(u.Linewisefile_chan("input")))
	fmt.Println("Part2: ", Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	lines2, linewidth := u.PeekAtFirstLineWidth(lines)

	var columnstops = make([]int, linewidth)
	for i := range columnstops {
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

func Part2(lines <-chan string) int {

	lines2, linewidth := u.PeekAtFirstLineWidth(lines)

	var round_rocks [][2]int
	var cube_rocks []cube_rock

	var row_count int
	for line := range lines2 {
		for i, ch := range []byte(line) {
			switch ch {
			case '#':
				cube_rocks = append(cube_rocks, cube_rock{Column: i, Row: row_count})
			case 'O':
				round_rocks = append(round_rocks, [2]int{i, row_count})
			}
		}
		row_count++
	}

	var omap = obstacle_map{make([]obstacle_node, row_count*linewidth), linewidth, row_count}

	// Populate
	for i, cr := range cube_rocks {
		opointer := omap.get(cr.Column, cr.Row)
		opointer.rock = &cube_rocks[i]
	}

	// Densify
	// Columns
	for col := 0; col < linewidth; col++ {
		cube_rocks = append(cube_rocks, cube_rock{Column: col, Row: -1})
		north_obstacle := &cube_rocks[len(cube_rocks)-1]
		cube_rocks = append(cube_rocks, cube_rock{Column: col, Row: row_count})
		south_obstacle := &cube_rocks[len(cube_rocks)-1]

		var next int
		var recurse func(int) *cube_rock
		recurse = func(row int) (obstacle *cube_rock) {
			if row >= row_count {
				next = row
				return south_obstacle
			}
			obs := omap.get(col, row)
			if obs.rock != nil {
				next = row + 1
				north_obstacle = obs.rock
				return obs.rock
			}

			obs.north = north_obstacle
			m_south_obs := recurse(row + 1)
			obs.south = m_south_obs
			return m_south_obs
		}

		for next < row_count {
			recurse(next)
		}
	}

	// Rows
	for row := 0; row < row_count; row++ {
		cube_rocks = append(cube_rocks, cube_rock{Column: -1, Row: row})
		west_obstacle := &cube_rocks[len(cube_rocks)-1]
		cube_rocks = append(cube_rocks, cube_rock{Column: linewidth, Row: row})
		east_obstacle := &cube_rocks[len(cube_rocks)-1]

		var next int
		var recurse func(int) *cube_rock
		recurse = func(col int) (obstacle *cube_rock) {
			if col >= linewidth {
				next = col
				return east_obstacle
			}
			obs := omap.get(col, row)
			if obs.rock != nil {
				next = col + 1
				west_obstacle = obs.rock
				return obs.rock
			}

			obs.west = west_obstacle
			m_east_obs := recurse(col + 1)
			obs.east = m_east_obs
			return m_east_obs
		}

		for next < linewidth {
			recurse(next)
		}
	}

	var load_list [][4]int
	var generation int

	spin_cycle := func() (loads [4]int) {
		// North
		for j, rock := range round_rocks {
			obstacle := omap.get(rock[0], rock[1]).north
			newrow := obstacle.Row + obstacle.get(generation)
			round_rocks[j][1] = newrow
			loads[0] += row_count - newrow
		}
		generation++

		// West
		for j, rock := range round_rocks {
			obstacle := omap.get(rock[0], rock[1]).west
			newcolumn := obstacle.Column + obstacle.get(generation)
			round_rocks[j][0] = newcolumn
			loads[1] += linewidth - newcolumn
		}
		generation++

		// South
		for j, rock := range round_rocks {
			obstacle := omap.get(rock[0], rock[1]).south
			newrow := obstacle.Row - obstacle.get(generation)
			round_rocks[j][1] = newrow
			loads[2] += row_count - newrow
		}
		generation++

		// East
		for j, rock := range round_rocks {
			obstacle := omap.get(rock[0], rock[1]).east
			newcolumn := obstacle.Column - obstacle.get(generation)
			round_rocks[j][0] = newcolumn
			loads[3] += linewidth - newcolumn
		}
		generation++

		return
	}

	load_list = append(load_list, spin_cycle())
	load_list = append(load_list, spin_cycle())
	load_list = append(load_list, spin_cycle())
	tortoise := 1
	hare := 2

	for load_list[tortoise] != load_list[hare] {
		load_list = append(load_list, spin_cycle())
		load_list = append(load_list, spin_cycle())
		tortoise += 1
		hare += 2
	}

	mu := 0
	tortoise = 0
	for load_list[tortoise] != load_list[hare] {
		load_list = append(load_list, spin_cycle())
		tortoise += 1
		hare += 1
		mu += 1
	}

	lambda := 1
	hare = tortoise + 1
	for load_list[tortoise] != load_list[hare] {
		hare += 1
		lambda += 1
	}

	fmt.Println("Cycle @", mu, "width", lambda)

	one_billion_index := ((1_000_000_000 - (mu + 1)) % lambda) + mu
	return load_list[one_billion_index][2]
}

type cube_rock struct {
	Column      int
	Row         int
	generation  int
	stack_count int
}

func (cr *cube_rock) get(cycle_gen int) (count int) {
	if cycle_gen > cr.generation {
		cr.stack_count = 0
		cr.generation = cycle_gen
	}
	cr.stack_count++
	return cr.stack_count
}

type obstacle_node struct {
	rock  *cube_rock // nil if not a rock
	north *cube_rock
	west  *cube_rock
	south *cube_rock
	east  *cube_rock
}

type obstacle_map struct {
	omap   []obstacle_node
	width  int
	height int
}

func (om *obstacle_map) get(x, y int) *obstacle_node {
	return &om.omap[(y*om.width)+x]
}

func print_map(round_rocks [][2]int, cube_rocks []cube_rock, width, height int) {
	rmap := make([]byte, width*height)

	ind := func(x, y int) int {
		return y*width + x
	}

	for _, rrock := range round_rocks {
		i := ind(rrock[0], rrock[1])
		rmap[i] = 'O'
	}
	for _, crock := range cube_rocks {
		if crock.Column < 0 || crock.Column >= width || crock.Row < 0 || crock.Row >= height {
			continue
		}
		i := ind(crock.Column, crock.Row)
		rmap[i] = '#'
	}

	for i, ch := range rmap {
		if i%width == 0 {
			fmt.Print("\n")
		}
		switch ch {
		case 0:
			fmt.Print(".")
		default:
			fmt.Print(string(ch))
		}
	}
	fmt.Print("\n")
}
