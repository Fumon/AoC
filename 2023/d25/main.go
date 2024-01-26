package main

import (
	"fmt"
	"fuaoc2023/day25/u"
	"slices"
	"strings"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var partmap = make(map[string]map[string]struct{})
	var insertfunc = func(s, t string) {
		if m, ok := partmap[s]; ok {
			m[t] = struct{}{}
		} else {
			nm := make(map[string]struct{})
			nm[t] = struct{}{}
			partmap[s] = nm
		}
	}
	var edges = make(map[[2]string]struct{})
	var insertedge = func(s, t string) {
		keyslice := []string{s,t}
		slices.Sort(keyslice)
		edges[[2]string{keyslice[0], keyslice[1]}] = struct{}{}
	}

	for line := range lines {
		sp := strings.Split(line, ": ")
		r := sp[0]
		for _, oth := range strings.Split(sp[1], " ") {
			insertfunc(r, oth)
			insertfunc(oth, r)
			insertedge(r,oth)
		}
	}

	// Use Floyd Warshall algorithm while tracking paths
	var dist = make(map[[2]string]int)
	for k := range partmap {
		for j := range partmap {
			key := make_string_key(k,j)
			if k == j {
				dist[key] = 0
			} else if _, ok := edges[key]; ok {
				dist[key] = 1
			} else {
				dist[key] = 1000000
			}
		}
	}
	for k := range partmap {
		for j := range partmap {
			for i := range partmap {
				keyji := make_string_key(j,i)
				keyjk := make_string_key(j,k)
				keyki := make_string_key(k,i)
				inter_dist := dist[keyjk] + dist[keyki]
				if dist[keyji] > inter_dist {
					dist[keyji] = inter_dist
				}
			}
		}
	}

	// Find the 6 parts with the highest centrality
	var centrality []cent_record
	for i := range partmap {
		var running_sum float64
		for j := range partmap {
			if i != j {
				key := make_string_key(i,j)
				if val, ok := dist[key]; ok {
					running_sum += float64(val)
				}
			}
		}
		centrality = append(centrality, cent_record{i, running_sum})
	}
	slices.SortFunc(centrality, func(a, b cent_record) int {
		if a.cent < b.cent {
			return 1
		} else if a.cent > b.cent {
			return -1
		} else {
			return 0
		}
	})

	// Print the 6 parts with the highest centrality
	for i := 0; i < len(centrality); i++ {
		fmt.Println(centrality[i])
	}


	return 0
}

func make_string_key (s, t string) [2]string{
	keyslice := []string{s,t}
	slices.Sort(keyslice)
	return [2]string{keyslice[0], keyslice[1]}
}

type cent_record struct {
	part string
	cent float64
}

func Part2(lines <-chan string) int {
	return 0
}
