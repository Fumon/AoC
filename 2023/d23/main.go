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
	start_v, _, _ := BuildGraph(lines)

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
	return longest_path(start_v)
}

func BuildGraph(lines <-chan string) (*Vertex, *Vertex, map[[2]int]*Vertex) {
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
		in:     []*Edge{},
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
		var vertex_exists bool
		if v, vertex_exists = verts[cur]; vertex_exists {
		} else {
			v = &Vertex{
				coords: cur,
				out:    []*Edge{},
				in:     []*Edge{},
			}
			verts[cur] = v
		}

		var edge = Edge{From: start_vertex, To: v, Length: cur_length}
		start_vertex.out = append(start_vertex.out, &edge)
		v.in = append(v.in, &edge)
		if reached_end {
			return
		}

		if !vertex_exists {
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
	}

	explore(&start_v, [2]int{1, 1})
	return &start_v, &end_v, verts
}


func Part2(lines <-chan string) int {
	start_v, end_v, _ := BuildGraph(lines)

	// Find the longest path from start_v to end_v which treats all edges as if they were undirected (i.e. each vertex can be left from the edges in vertex.out via those edges' .to as well as vertex.in via those edges' .from) and avoids cycles of any length (i.e. for a given path, no vertex can appear more than once in that path)
	// The following recursive function signature is not acceptable as it only prevents cycles of length 1: func longest_path(*Vertex, *Vertex) int
	// Instead, we need to keep track of the path so far and check for cycles of any length
	// var longests = map[string]int{}
	var longest_path func(v *Vertex, path_so_far map[*Vertex]struct{}) (int, bool)
	longest_path = func(v *Vertex, path_so_far map[*Vertex]struct{}) (length int, reached_end bool) {
		if v == end_v {
			return 0, true
		}
		var longest int
		path_after_here := make(map[*Vertex]struct{}, len(path_so_far)+1)
		for vertex := range path_so_far {
			path_after_here[vertex] = struct{}{}
		}
		path_after_here[v] = struct{}{}

		
		Edgeloop1:
		for _, edge := range v.out {
			next := edge.To
			if _, ok := path_so_far[next]; ok {
				continue Edgeloop1
			}
			path_length, t_reached_end := longest_path(next, path_after_here)
			if t_reached_end {
				reached_end = true
				var length = edge.Length + path_length 
				if length > longest {
					longest = length
				}
			}
		}
		Edgeloop2:
		for _, edge := range v.in {
			next := edge.From
			if _, ok := path_so_far[next]; ok {
				continue Edgeloop2
			}
			path_length, t_reached_end := longest_path(next, path_after_here)
			if t_reached_end {
				reached_end = true
				var length = edge.Length + path_length 
				if length > longest {
					longest = length
				}
			}
		}
		return longest, reached_end
	}
	len, _ := longest_path(start_v, map[*Vertex]struct{}{})
	return len
}
