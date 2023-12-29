package main

import (
	"fmt"
	"fuaoc2023/day18/u"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var cur = [2]int{0,0}
	var prev = cur
	var area_running int
	var bounds = Square_Bounds{}
	var dug = make(map[[2]int]Color)
	dug[cur] = Color{255, 255, 255}

	for line := range lines {
		dir, run, color := parse_command(line)

		offset := dir_to_offset[dir]
		for i := 0; i < run; i++ {
			prev = cur
			cur[0] += offset[0]
			cur[1] += offset[1]
			dug[cur] = color
			bounds.update(cur)
			
			area_running += (prev[0] * cur[1] - cur[0] * prev[1])
		}
	}
	// Can skip the closing point since it's zero

	print_dug(dug, bounds)
	area := area_running / 2
	if area < 0 {
		area = -area
	}
	return area + len(dug) / 2 + 1
}

var byte_to_offset = map[byte][2]int64{
	'3': {0, -1},
	'1': {0, 1},
	'2': {-1, 0},
	'0': {1, 0},
}

func Part2(lines <-chan string) int64 {
	var area_running int64
	var cur = [2]int64{0,0}
	var dug int64
	var prev = cur

	for line := range lines {
		dir, run := parse_command2(line)
		offset := byte_to_offset[dir]
		prev = cur
		cur[0] += offset[0]*run
		cur[1] += offset[1]*run
		area_running += (prev[0] * cur[1] - cur[0] * prev[1])
		dug += run
	}

	area := area_running / 2
	if area < 0 {
		area = -area
	}

	return area + dug / 2 + 1
}

var dir_to_offset = map[byte][2]int{
	'U': {0, -1},
	'D': {0, 1},
	'L': {-1, 0},
	'R': {1, 0},
}

type Color struct {
	R, G, B byte
}

func parse_command(line string) (dir byte, run int, color Color) {
	comp := strings.Split(line, " ")
	dir = comp[0][0]
	run, _ = strconv.Atoi(comp[1])
	R, _ := strconv.ParseInt(comp[2][2:4], 16, 8)
	G, _ := strconv.ParseInt(comp[2][4:6], 16, 8)
	B, _ := strconv.ParseInt(comp[2][6:8], 16, 8)

	color = Color{
		R: byte(R),
		G: byte(G),
		B: byte(B),
	}
	return
}

func parse_command2(line string) (dir byte, run int64) {
	comp := strings.Split(line, " ")
	run, _ = strconv.ParseInt(comp[2][2:7], 16, 64)
	dir = byte(comp[2][7])

	return
}


type Square_Bounds struct {
	xmin, ymin int
	xmax, ymax int
}

func (b *Square_Bounds) update(new_point [2]int) {
	if new_point[0] > b.xmax {
		b.xmax = new_point[0]
	} else if new_point[0] < b.xmin {
		b.xmin = new_point[0]
	}

	if new_point[1] > b.ymax {
		b.ymax = new_point[1]
	} else if new_point[1] < b.ymin {
		b.ymin = new_point[1]
	}
}

func print_dug(dug map[[2]int]Color, bounds Square_Bounds) {
	file, err := os.Create("dug.ppm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	width := bounds.xmax - bounds.xmin + 1
	height := bounds.ymax - bounds.ymin + 1

	_, err = file.WriteString(fmt.Sprintf("P3\n%d %d\n255\n", width, height))
	if err != nil {
		panic(err)
	}
	// fmt.Print("Bounds: ", bounds)

	for y := bounds.ymin; y <= bounds.ymax; y++ {
		for x := bounds.xmin; x <= bounds.xmax; x++ {
			if color, found := dug[[2]int{x, y}]; found {
				// if x == 0 && y == 0 {
				// 	color = Color{255, 0, 255}
				// } else {
				// 	color = Color{255, 255, 255}
				// }
				file.WriteString(fmt.Sprintf("%d %d %d ", color.R, color.G, color.B))
			} else {
				color = Color{0,0,0}
				file.WriteString(fmt.Sprintf("%d %d %d ", color.R, color.G, color.B))
			}
		}
		file.WriteString("\n")
	}
	// fmt.Print("Wrote dug map to dug.ppm")
}
