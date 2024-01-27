package main

import (
	"fmt"
	"fuaoc2023/day25/ek"
	"fuaoc2023/day25/u"
	"slices"
	"strings"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var injest_partmap = make(map[string]map[string]struct{})
	var injest_edges = make(map[[2]string]struct{})
	{
		var insertfunc = func(s, t string) {
			if m, ok := injest_partmap[s]; ok {
				m[t] = struct{}{}
			} else {
				nm := make(map[string]struct{})
				nm[t] = struct{}{}
				injest_partmap[s] = nm
			}
		}
		var insertedge = func(s, t string) {
			keyslice := []string{s, t}
			slices.Sort(keyslice)
			injest_edges[[2]string{keyslice[0], keyslice[1]}] = struct{}{}
		}

		for line := range lines {
			sp := strings.Split(line, ": ")
			r := sp[0]
			for _, oth := range strings.Split(sp[1], " ") {
				insertfunc(r, oth)
				insertfunc(oth, r)
				insertedge(r, oth)
			}
		}
	}

	var partcount = len(injest_partmap)
	var parts_name_map = make(map[string]int)
	var parts []string
	var connections = make([][]bool, partcount)
	{
		var part_index_count int
		for partname := range injest_partmap {
			parts_name_map[partname] = part_index_count
			connections[part_index_count] = make([]bool, partcount)
			parts = append(parts, partname)
			part_index_count++
		}
		for pn, links := range injest_partmap {
			pn_index := parts_name_map[pn]
			connection_array := connections[pn_index]
			for link := range links {
				connection_array[parts_name_map[link]] = true
			}
			// connection_array[pn_index] = true
		}
	}

	var last_most_popular [2]int
	for j := 0; j < 3; j++ {
		popularity_rankings := Rank_Edges(connections)

		fmt.Println("===== Iteration", j , "=====")
		most_popular := popularity_rankings[0]
		indicies := most_popular.e.Recover_indicies()
		connections[indicies[0]][indicies[1]] = false
		connections[indicies[1]][indicies[0]] = false

		for i, pop := range popularity_rankings[:10] {
			indicies := pop.e.Recover_indicies()
			firstname := parts[indicies[0]]
			secondname := parts[indicies[1]]
			fmt.Println(i, "-", firstname, "/", secondname, ":", pop.visits)
		}
		last_most_popular = indicies
	}

	var neighborhood_sizes []int
	for i := 0; i < 2; i++ {
		startpoint := last_most_popular[i]
		var neighborhood = map[int]struct{}{startpoint: {}}
		var front = []int{startpoint}
		for len(front) > 0 {
			n := front[0]
			front = front[1:]
			for dest, connected := range connections[n] {
				if !connected {
					continue
				} else if _, ok := neighborhood[dest]; ok {
					continue
				}

				front = append(front, dest)
				neighborhood[dest] = struct{}{}
			}
		}
		neighborhood_sizes = append(neighborhood_sizes, len(neighborhood))
	}

	fmt.Println(neighborhood_sizes)

	return neighborhood_sizes[0] * neighborhood_sizes[1]
}

func Rank_Edges(connections [][]bool) []pop_entry {
	var popularity_map = make(map[ek.EdgeKey]int)
	{
		for part_id := range connections {
			var visited = make(map[int]struct{})
			visited[part_id] = struct{}{}
			var front_ids = []int{part_id}
			for len(front_ids) > 0 {
				n := front_ids[0]
				front_ids = front_ids[1:]
				for dest, connected := range connections[n] {
					if !connected {
						continue
					} else if _, ok := visited[dest]; ok {
						continue
					}

					edge_key := ek.New(n, dest)
					popularity_map[edge_key]++
					front_ids = append(front_ids, dest)
					visited[dest] = struct{}{}
				}
			}
		}
	}

	var popularity_rankings []pop_entry
	{
		for edge_key, visit_count := range popularity_map {
			popularity_rankings = append(popularity_rankings, pop_entry{edge_key, visit_count})
		}
		slices.SortFunc(popularity_rankings, func(a, b pop_entry) int { return b.visits - a.visits })
	}
	return popularity_rankings
}

func Part2(lines <-chan string) int {
	return 0
}
