package main

import (
	"fmt"
	"fuaoc2023/day23/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var grid [][]byte
	for line := range lines {
		grid = append(grid, []byte(line))
	}

	start_v := Vertex{
		coords: [2]int{1, 0},
		out:    []*Edge{},
	}
	var end_v = Vertex{
		coords: [2]int{len(grid[0]) - 2, len(grid) - 1},
	}

	var verts = map[[2]int]*Vertex{start_v.coords: &start_v, end_v.coords: &end_v}

	var explore func(*Vertex, [2]int)
	explore = func(start_vertex *Vertex, start_coords [2]int) {
		var cur_length int = 1
		var previous, cur [2]int
		var reached_end bool
		previous = start_vertex.coords
		cur = start_coords
		TRAVERSE_LOOP:
		for {
			var delta [2]int
			delta[0], delta[1] = (cur[1] - previous[1]), -(cur[0] - previous[0])
			for turn := 0; turn < 3; turn++ {
				var next [2]int
				next[0], next[1] = cur[0]+delta[0], cur[1]+delta[1]
				switch grid[next[1]][next[0]] {
				case '#':
					// turn right
					delta[0], delta[1] = -delta[1], delta[0]
					continue
				case '.':
					cur_length++
					previous = cur
					cur = next

					if next == end_v.coords {
						reached_end = true
						break TRAVERSE_LOOP
					}
					continue TRAVERSE_LOOP
				case '>', '<', '^', 'v':
					cur_length += 2
					previous = next
					cur[0], cur[1] = next[0]+delta[0], next[1]+delta[1]
					break TRAVERSE_LOOP
				}
			}
			panic("Dead End!")
		}

		var v *Vertex
		if vert, ok := verts[cur]; ok {
			v = vert
		} else {
			v = &Vertex{
				coords: cur,
				out:    []*Edge{},
			}
			verts[cur] = v
		}

		var edge = Edge{From: start_vertex, To: v, Length: cur_length}
		start_vertex.out = append(start_vertex.out, &edge)
		if reached_end {
			return
		}

		var directions = map[[2]int]byte{{0, -1}: '^', {1, 0}: '>', {0, 1}: 'v', {-1, 0}: '<'}
		for delta, direction := range directions {
			var next [2]int
			next[0], next[1] = cur[0]+delta[0], cur[1]+delta[1]
			if next == previous {
				continue
			} else if grid[next[1]][next[0]] == direction {
				explore(v, next)
			}
		}
	}
	
	explore(&start_v, [2]int{1, 1})

	// Find the longest path from start_v to end_v
	var longests = map[*Vertex]int{}
	var longest_path func(*Vertex) int
	longest_path = func(v *Vertex) int {
		if longest, ok := longests[v]; ok {
			return longest
		}
		var longest int
		for _, edge := range v.out {
			var length = edge.Length + longest_path(edge.To)
			if length > longest {
				longest = length
			}
		}
		longests[v] = longest
		return longest
	}
	return longest_path(&start_v)
}


func Part2(lines <-chan string) int {
	return 0
}
