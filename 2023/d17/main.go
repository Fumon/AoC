package main

import (
	"container/heap"
	"fmt"
	"fuaoc2023/day17/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_slice("input")))
	fmt.Println(Part2(u.Linewisefile_slice("input")))
}

type node_map struct {
	lookup  map[node_name]*node
	storage []node
}

func NewNodeMap() node_map {
	return node_map{
		lookup: make(map[node_name]*node),
	}
}

func (nm *node_map) Get(name *node_name) *node {

	if val, ok := nm.lookup[*name]; ok {
		return val
	}

	nm.storage = append(nm.storage, node{name: name, distance: int(^uint(0) >> 1)})
	ptr := &nm.storage[len(nm.storage)-1]
	nm.lookup[*name] = ptr
	return ptr

}

var TurnMap = map[byte][2]byte{
	'S': {'E', 'W'},
	'N': {'W', 'E'},
	'E': {'N', 'S'},
	'W': {'S', 'N'},
}
var OffsetMap = map[byte][2]int{
	'S': {0, -1},
	'N': {0, 1},
	'E': {-1, 0},
	'W': {1, 0},
}

// P1 Stencil:
//
//

func Part1(lines []string) int {

	// Build graph
	nodes, _, _, end := populate_graph(lines, 1, 3)

	return Shortest_Path(nodes, end)

	// print_path(end, lines)
}

func Part2(lines []string) int {
		// Build graph
		nodes, _, _, end := populate_graph(lines, 4, 10)

		return Shortest_Path(nodes, end)
}

func Shortest_Path(nodes node_map, end *node) int {
	unvisited := NewNodeHeap()
	heap.Init(unvisited)
	for _, nodep := range nodes.lookup {
		heap.Push(unvisited, nodep)
	}

	for !end.visited {
		cur := heap.Pop(unvisited).(*node)

		for _, edge := range cur.edges_out {
			neighbour := edge.end
			if neighbour.visited {
				continue
			}

			m_distance := cur.distance + edge.cost
			if m_distance < neighbour.distance {
				neighbour.prev = cur
				unvisited.UpdateNode(neighbour, m_distance)
			}
		}
		cur.visited = true
	}
	return end.distance
}

// func print_path(end *node, lines []string) {
// 	var shortpath = make(map[[2]int]byte)

// 	{
// 		cur := end
// 		for cur != start {
// 			cur = cur.prev
// 			shortpath[[2]int{cur.name.x, cur.name.y}] = '#'
// 		}

// 		for y := 0; y < height; y++ {
// 			for x := 0; x < width; x++ {
// 				if _, ok := shortpath[[2]int{x, y}]; ok {
// 					fmt.Print("#")
// 				} else {
// 					fmt.Print(string(lines[y][x]))
// 				}
// 			}
// 			fmt.Print("\n")
// 		}
// 	}
// }

func populate_graph(lines []string, min_turn uint8, max_straight uint8) (node_map, []edge, *node, *node) {
	var nodes = NewNodeMap()
	var edges []edge
	var width = len(lines[0])
	var height = len(lines)
	for y, line := range lines {
		for x, ch := range []byte(line) {
			cost := int(ch & 0xF)

			for dir, offset := range OffsetMap {
				new_x, new_y := x+offset[0], y+offset[1]
				if new_x >= 0 && new_x < width && new_y >= 0 && new_y < height {
					{
						to_1 := nodes.Get(&node_name{
							x:     x,
							y:     y,
							dir:   dir,
							level: 1,
						})

						turns := TurnMap[dir]
						for sidelevel := uint8(min_turn); sidelevel <= max_straight; sidelevel++ {

							e_node := nodes.Get(&node_name{
								x:     new_x,
								y:     new_y,
								dir:   turns[0],
								level: sidelevel,
							})
							w_node := nodes.Get(&node_name{
								x:     new_x,
								y:     new_y,
								dir:   turns[1],
								level: sidelevel,
							})

							edges = append(edges, edge{
								start: e_node,
								end:   to_1,
								cost:  cost,
							})
							e_node.edges_out = append(e_node.edges_out, &edges[len(edges)-1])

							edges = append(edges, edge{
								start: w_node,
								end:   to_1,
								cost:  cost,
							})
							w_node.edges_out = append(w_node.edges_out, &edges[len(edges)-1])
						}
					}

					{
						for from_level := uint8(1); from_level < max_straight; from_level++ {
							from := nodes.Get(&node_name{
								x:     new_x,
								y:     new_y,
								dir:   dir,
								level: from_level,
							})

							to := nodes.Get(&node_name{
								x:     x,
								y:     y,
								dir:   dir,
								level: from_level+1,
							})

						

							edges = append(edges, edge{
								start: from,
								end:   to,
								cost:  cost,
							})
							from.edges_out = append(from.edges_out, &edges[len(edges)-1])
						}
					}
				}
			}
		}
	}

	// Add start and end
	start := nodes.Get(&node_name{
		x:     -1,
		y:     -1,
		dir:   0,
		level: 0,
	})
	start.distance = 0
	{
		edges = append(edges, edge{
			start: start,
			end: nodes.Get(&node_name{
				x:     1,
				y:     0,
				dir:   'E',
				level: 1,
			}),
			cost: int(lines[0][1] & 0xF),
		})
		start.edges_out = append(start.edges_out, &edges[len(edges)-1])

		edges = append(edges, edge{
			start: start,
			end: nodes.Get(&node_name{
				x:     0,
				y:     1,
				dir:   'S',
				level: 1,
			}),
			cost: int(lines[1][0] & 0xF),
		})
		start.edges_out = append(start.edges_out, &edges[len(edges)-1])
	}

	end := nodes.Get(&node_name{
		x:     -2,
		y:     -2,
		dir:   0,
		level: 0,
	})
	{
		for start_level := uint8(min_turn); start_level <= max_straight; start_level++ {
			lrc_e := nodes.Get(&node_name{
				x:     width - 1,
				y:     height - 1,
				dir:   'E',
				level: start_level,
			})

			lrc_s := nodes.Get(&node_name{
				x:     width - 1,
				y:     height - 1,
				dir:   'S',
				level: start_level,
			})

			edges = append(edges, edge{
				start: lrc_e,
				end:   end,
				cost:  0,
			})
			lrc_e.edges_out = append(lrc_e.edges_out, &edges[len(edges)-1])

			edges = append(edges, edge{
				start: lrc_s,
				end:   end,
				cost:  0,
			})
			lrc_s.edges_out = append(lrc_s.edges_out, &edges[len(edges)-1])
		}
	}

	return nodes, edges, start, end
}

type NodeHeap struct {
	nodes    []*node
	indexMap map[*node]int
}

func (h NodeHeap) Len() int { return len(h.nodes) }

func (h NodeHeap) Less(i, j int) bool {
	return h.nodes[i].distance < h.nodes[j].distance
}

func (h NodeHeap) Swap(i, j int) {
	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
	h.indexMap[h.nodes[i]] = i
	h.indexMap[h.nodes[j]] = j
}

func (h *NodeHeap) Push(x interface{}) {
	n := x.(*node)
	h.indexMap[n] = len(h.nodes)
	h.nodes = append(h.nodes, n)
}

func (h *NodeHeap) Pop() interface{} {
	old := h.nodes
	n := len(old)
	node := old[n-1]
	h.nodes = old[0 : n-1]
	delete(h.indexMap, node)
	return node
}

func NewNodeHeap() *NodeHeap {
	return &NodeHeap{
		nodes:    []*node{},
		indexMap: make(map[*node]int),
	}
}

// UpdateNode modifies a node's distance and reestablishes the heap invariants.
func (h *NodeHeap) UpdateNode(n *node, distance int) {
	index, ok := h.indexMap[n]
	if !ok {
		// Node not in heap
		return
	}
	n.distance = distance
	heap.Fix(h, index)
}

type node_name struct {
	x, y  int
	dir   byte
	level uint8
}

func (nn *node_name) String() string {
	return fmt.Sprint("(", nn.x, ",", nn.y, ") ", string(nn.dir), "@", nn.level)
}

type edge struct {
	start *node
	end   *node
	cost  int
}

func (e *edge) String() string {
	return fmt.Sprint("{", e.cost, "|", e.start.name, " -> ", e.end.name, "}")
}

type node struct {
	name      *node_name
	visited   bool
	distance  int
	prev      *node
	edges_out []*edge
}
