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

	// // Print partmap
	// for p := range partmap {
	// 	var othslice []string
	// 	for oth := range partmap[p] {
	// 		othslice = append(othslice, oth)
	// 	}
	// 	slices.Sort(othslice)

	// 	fmt.Print(p, " : {")
	// 	for _, o := range othslice {
	// 		fmt.Print(" ", o)
	// 	}
	// 	fmt.Print(" }\n")
	// }

	// Find all paths
	// Enumerate nodes
	var node_enum []string
	var pcount int
	for p := range partmap {
		node_enum = append(node_enum, p)
		pcount++
	}


	var distance = make([][]int, pcount)
	var next = make([][]int, pcount)
	for i, name := range node_enum {
		distance[i] = make([]int, pcount)
		next[i] = make([]int, pcount)
		for j, oname := range node_enum {
			keyslice := []string{name,oname}
			slices.Sort(keyslice)
			if _, ok := edges[[2]string{keyslice[0], keyslice[1]}]; ok {
				distance[i][j] = 1
				next[i][j] = j
			} else if i == j {
				distance[i][j] = 0
				next[i][j] = -1
			} else {
				distance[i][j] = 900000
				next[i][j] = -1
			}
		}
	}

	for k := range node_enum {
		for i := range node_enum {
			for j := range node_enum {
				inter_dist := distance[i][k] + distance[k][j]
				if inter_dist < distance[i][j] {
					distance[i][j] = inter_dist
					next[i][j] = next[i][k]
				}
			}
		}
	}

	fmt.Println("Distances")
	for i, name := range node_enum {
		for j, oname := range  node_enum {
			fmt.Printf("%v  %v: %d\n", name, oname, distance[i][j])
		}
	}

	// Most popular edges
	var pop_edge_map = make(map[[2]string]int, len(edges))
	for i := 0; i < len(node_enum) - 1; i++ {
		for j := i; j < len(node_enum); j++ {
			if i == j {
				continue
			}
			if next[i][j] == -1 {
				continue
			}
			var mi, mj = i, j
			for mi != mj {
				n := next[mi][mj]
				nname := node_enum[n]
				miname := node_enum[mi]

				edge_key_slice := []string{nname, miname}
				slices.Sort(edge_key_slice)
				edge_key := [2]string{edge_key_slice[0], edge_key_slice[1]}
				pop_edge_map[edge_key]++
				// if ep, ok := pop_edge_map[edge_key]; ok {
				// 	(*ep).count++
				// } else {
				// 	name := rev_enum[edge_key[0]] + "/" + rev_enum[edge_key[1]]
				// 	pop_edge_map[edge_key] = &pop_edge{
				// 		name:  name,
				// 		count: 1,
				// 	}
				// }
				mi = n
			}
		}
	}

	// Dump them all out and sort
	var pop_edge_list []*pop_edge
	for name, ep := range pop_edge_map {
		pop_edge_list = append(pop_edge_list, &pop_edge{count:ep, name: name[0] + "/" + name[1]})
	}
	slices.SortFunc(pop_edge_list, func(a, b *pop_edge) int {
		return -(a.count - b.count)
	})

	for i, ep := range pop_edge_list {
		fmt.Printf("%d: %v - %v\n", i, ep.name, ep.count)
	}
	

	// Save edges as a graphviz dot file
	// f, err := os.Create(".dot")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// fmt.Fprintln(f, "graph {")
	// for e := range edges {
	// 	fmt.Fprintf(f, "  %s -- %s;\n", e[0], e[1])
	// }
	// fmt.Fprintln(f, "}")

	return 0
}

type pop_edge struct {
	name string
	count int
}

func Part2(lines <-chan string) int {
	return 0
}
