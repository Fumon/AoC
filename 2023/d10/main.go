package main

import (
	"fmt"
	"fuaoc2023/day10/u"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))

}

// 0: | is a vertical pipe connecting north and south.
// 1: - is a horizontal pipe connecting east and west.
// 2: L is a 90-degree bend connecting north and east.
// 3: J is a 90-degree bend connecting north and west.
// 4: 7 is a 90-degree bend connecting south and west.
// 5: F is a 90-degree bend connecting south and east.
// 6: . is ground; there is no pipe in this tile.
// 7: S is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.

type PipeTile struct {
	tiletype byte
}

func (p PipeTile) String() string {
	// switch p.tiletype {
	// case 0:
	// 	return "|"
	// case 1:
	// 	return "-"
	// case 2:
	// 	return "L"
	// case 3:
	// 	return "J"
	// case 4:
	// 	return "7"
	// case 5:
	// 	return "F"
	// case 6:
	// 	return "."
	// case 7:
	// 	return "S"
	// default:
	// 	panic("FDSAFDSAFDSAFDSAFASDFDSA")
	// }
	return string(p.tiletype)
}

func NewPipeTile(b byte) (isStart bool, p PipeTile) {
	p.tiletype = b
	if b == 'S' {
		isStart = true
	}
	return
}

var ConnectionNext = map[[2]byte]offsetFunc{
	{'|', 'S'}: Up,
	{'J', 'W'}: Up,
	{'L', 'E'}: Up,

	{'|', 'N'}: Down,
	{'7', 'W'}: Down,
	{'F', 'E'}: Down,

	{'-', 'E'}: Left,
	{'J', 'N'}: Left,
	{'7', 'S'}: Left,

	{'-', 'W'}: Right,
	{'L', 'N'}: Right,
	{'F', 'S'}: Right,
}

type Connection struct {
	G        byte
	Location [2]int
}

type TileConnection struct {
	*PipeTile
	Connection
}

func (c TileConnection) Next() Connection {
	return ConnectionNext[[2]byte{c.tiletype, c.G}](c.Connection)
}

func Up(cl Connection) Connection {
	return Connection{
		G:        'S',
		Location: [2]int{cl.Location[0] - 1, cl.Location[1]},
	}
}
func Down(cl Connection) Connection {
	return Connection{
		G:        'N',
		Location: [2]int{cl.Location[0] + 1, cl.Location[1]},
	}
}
func Left(cl Connection) Connection {
	return Connection{
		G:        'E',
		Location: [2]int{cl.Location[0], cl.Location[1] - 1},
	}
}
func Right(cl Connection) Connection {
	return Connection{
		G:        'W',
		Location: [2]int{cl.Location[0], cl.Location[1] + 1},
	}
}

type offsetFunc func(Connection) Connection

func part1(lines <-chan string) int {
	var field [][]*PipeTile
	var start [2]int
	start, field = parseField(lines, start, field)

	// Find start connections
	var connections []TileConnection
	connections = findStartConnections(field, start, connections)

	var distance_until_met int = 1
	for {
		for i, v := range connections {
			c := v.Next()
			connections[i] = TileConnection{
				PipeTile:   field[c.Location[0]][c.Location[1]],
				Connection: c,
			}
		}
		distance_until_met++
		if connections[0].PipeTile == connections[1].PipeTile {
			break
		}
	}
	return distance_until_met
}

type Edge struct {
	*Connection
	H byte
}

func part2(lines <-chan string) int {
	var field [][]*PipeTile
	var start [2]int
	start, field = parseField(lines, start, field)

	var bounds = NewFieldBounds(start)

	// Find start connections
	var connections []TileConnection
	connections = findStartConnections(field, start, connections)[0:1]

	
	var edges = make(map[[2]int]*Edge)
	con := connections[0]
	bounds.update(con.Location)
	edges[con.Location] = &Edge{&con.Connection, 0}

	for con.tiletype != 'S' {
		n := con.Next()
		bounds.update(n.Location)
		v := TileConnection{
			PipeTile:   field[n.Location[0]][n.Location[1]],
			Connection: n,
		}
		edges[con.Location].H = n.G
		edges[n.Location] = &Edge{&n, 0}
		connections = append(connections, v)
		con = v
	}
	// Close
	first := edges[connections[0].Location]
	edges[con.Location].H = first.G
	first.H = con.G

	var inside = make(map[[2]int]bool)
	var inside_count int
	// Generate all candidate points
	for i := bounds.ymin + 1; i < bounds.ymax; i++ {
		for j := bounds.xmin + 1; j < bounds.xmax; j++ {
			var p = [2]int{i, j}
			if _, ok := edges[p]; ok {
				continue
			}

			dir, conn := bounds.closestBoundDelta(p)


			var crossings int
			for !bounds.testOutside(p) {
				p[0] += dir[0]
				p[1] += dir[1]
				if pt, ok := edges[p]; ok {
					crossings += crossingMap[[2]byte{conn, pt.G}]
					crossings += crossingMap[[2]byte{conn, pt.H}]
				}
			}

			if crossings != 0 {
				inside_count++
				inside[[2]int{i, j}] = true
			} else {
				inside[[2]int{i, j}] = false
			}
		}
	}


	for i, fieldline := range field {
		for j, tile := range fieldline {
			isInside, ok := inside[[2]int{i, j}]
			if ok {
				if isInside {
					fmt.Print("\033[1;92mI\033[0m")
				} else {
					fmt.Print("\033[1;91mO\033[0m")
				}
			} else {
				fmt.Print(string(tile.tiletype))
			}
		}
		fmt.Println("")
	}

	return inside_count
}

type fieldBounds struct {
	xmin int
	xmax int
	ymin int
	ymax int
}

func NewFieldBounds(s [2]int) fieldBounds {
	return fieldBounds{
		s[1],
		s[1],
		s[0],
		s[0],
	}
}

func (fb *fieldBounds) update(p [2]int) {
	if p[1] < fb.xmin {
		fb.xmin = p[1]
	} else if p[1] > fb.xmax {
		fb.xmax = p[1]
	}

	if p[0] < fb.ymin {
		fb.ymin = p[0]
	} else if p[0] > fb.ymax {
		fb.ymax = p[0]
	}
}

// Right to left is 1
var crossingMap = map[[2]byte]int{
	{'N', 'N'}: 0,
	{'N', 'S'}: 0,
	{'N', 'E'}: 1,
	{'N', 'W'}: -1,

	{'S', 'S'}: 0,
	{'S', 'N'}: 0,
	{'S', 'E'}: -1,
	{'S', 'W'}: 1,

	{'E', 'E'}: 0,
	{'E', 'W'}: 0,
	{'E', 'N'}: 1,
	{'E', 'S'}: -1,


	{'W', 'W'}: 0,
	{'W', 'E'}: 0,
	{'W', 'N'}: -1,
	{'W', 'S'}: 1,
}

func (fb *fieldBounds) closestBoundDelta(p [2]int) (dir [2]int, connection byte) {
	g1 := p[0] - fb.ymin
	g2 := fb.ymax - p[0]
	
	var gneg bool
	var g = g2
	if g1 < g2 {
		g = g1
		gneg = true
	}

	f1 := p[1] - fb.xmin
	f2 := fb.xmax - p[1]

	var fneg bool
	var f = f2
	if f1 < f2 {
		f = f1
		fneg = true
	}

	if g < f {
		connection = 'N'
		dir[0] = 1
		if gneg {
			connection = 'S'
			dir[0] = -1
		}
	} else {
		connection = 'W'
		dir[1] = 1
		if fneg {
			connection = 'E'
			dir[1] = -1
		}
	}
	return
}

func (fb *fieldBounds) testOutside(p [2]int) bool {
	return p[0] < fb.ymin || p[0] > fb.ymax || p[1] < fb.xmin || p[1] > fb.xmax
}

func findStartConnections(field [][]*PipeTile, start [2]int, connections []TileConnection) []TileConnection {
	var xbound int = len(field[0])
	var ybound int = len(field)
	{
		if start[1] > 0 {
			leftloc := [2]int{start[0], start[1] - 1}
			left := field[leftloc[0]][leftloc[1]]
			switch left.tiletype {
			case '-', 'L', 'F':
				connections = append(connections, TileConnection{PipeTile: left, Connection: Connection{'E', leftloc}})
			}
		}

		if start[0] > 0 {
			uploc := [2]int{start[0] - 1, start[1]}
			up := field[uploc[0]][uploc[1]]
			switch up.tiletype {
			case '|', '7', 'F':
				connections = append(connections, TileConnection{PipeTile: up, Connection: Connection{'S', uploc}})
			}
		}

		if start[1] < xbound-1 {
			rightloc := [2]int{start[0], start[1] + 1}
			right := field[rightloc[0]][rightloc[1]]
			switch right.tiletype {
			case '-', 'J', '7':
				connections = append(connections, TileConnection{PipeTile: right, Connection: Connection{'W', rightloc}})
			}
		}

		if start[0] < ybound-1 {
			downloc := [2]int{start[0] + 1, start[1]}
			down := field[downloc[0]][downloc[1]]
			switch down.tiletype {
			case '|', 'L', 'J':
				connections = append(connections, TileConnection{PipeTile: down, Connection: Connection{'N', downloc}})
			}
		}
	}

	if len(connections) != 2 {
		panic("Wrong number of connections")
	}
	return connections
}

func parseField(lines <-chan string, start [2]int, field [][]*PipeTile) ([2]int, [][]*PipeTile) {
	var startfound bool
	var lineindex int
	for line := range lines {
		var fieldline []*PipeTile
		for i, v := range []byte(line) {
			isStart, newTile := NewPipeTile(v)
			if isStart {
				start = [2]int{lineindex, i}
				startfound = true
			}
			fieldline = append(fieldline, &newTile)
		}
		field = append(field, fieldline)
		lineindex++
	}
	if !startfound {
		panic("Start not found")
	}
	return start, field
}
