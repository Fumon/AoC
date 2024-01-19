package main

import (
	"fmt"
	"fuaoc2023/day22/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var bricks []Brick
	var numbricks int
	var max_x, max_y uint16

	for line := range lines {
		new_brick := ParseBrick(line)
		max_x = max(new_brick.Start[0], new_brick.End[0], max_x)
		max_y = max(new_brick.Start[1], new_brick.End[1], max_y)

		bricks = append(bricks, new_brick)
		numbricks++
	}

	var bricks_at_rest []RestingBrick
	var colliding_lookup = make(map[uint16][]*RestingBrick)



	return 0
}

func Part2(lines <-chan string) int {
	return 0
}
