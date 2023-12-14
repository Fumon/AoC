package main

import (
	"fmt"
	"fuaoc2023/day10/u"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
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
	// switch b {
	// case '|':
	// 	p.tiletype = 0
	// case '-':
	// 	p.tiletype = 1
	// case 'L':
	// 	p.tiletype = 2
	// case 'J':
	// 	p.tiletype = 3
	// case '7':
	// 	p.tiletype = 4
	// case 'F':
	// 	p.tiletype = 5
	// case '.':
	// 	p.tiletype = 6
	// case 'S':
	// 	p.tiletype = 7
	// 	isStart = true
	// }
	p.tiletype = b
	if b == 'S' {
		isStart = true
	}
	return
}

var ConnectionNext = map[[2]byte]offsetFunc {
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
	G byte
	Location [2]int
}

type TileConnection struct {
	*PipeTile
	Connection
}

func (c TileConnection) Next() Connection {
	return ConnectionNext[[2]byte{c.tiletype, c.G}](c.Connection)
}

// func (c TileConnection) Next() (out [2]int, connection byte) {
// 	out[0] = c.Location[0]
// 	out[1] = c.Location[1]

// 	switch c.Connection {
// 	case 'N':
// 		switch c.tiletype {
// 		case '|':
// 			connection = 'N'
// 		case 'L':
// 			connection = 'W'
// 		case 'J':
// 			connection = 'E'
// 		}
// 	case 'S':
// 		switch c.tiletype {
// 		case '|':
// 			connection = 'S'
// 		case 'F':
// 			connection = 'W'
// 		case '7':
// 			connection = 'E'
// 		}
// 	case 'W':
// 		switch c.tiletype {
// 		case '-':
// 			connection = 'W'
// 		case 'J':
// 			connection = 'S'
// 		case '7':
// 			connection = 'N'
// 		}
// 	case 'E':
// 		switch c.tiletype {
// 		case '-':
// 			connection = 'E'
// 		case 'L':
// 			connection = 'S'
// 		case 'F':
// 			connection = 'N'
// 		}
// 	default:
// 		panic("bad connection")
// 	}

// 	switch connection {
// 	case 'N':
// 		out[0] += 1
// 	case 'S':
// 		out[0] -= 1
// 	case 'E':
// 		out[1] -= 1
// 	case 'W':
// 		out[1] += 1
// 	default:
// 		panic("bad after connection")
// 	}

// 	return
// }

func Up(cl Connection) Connection {
	return Connection{
		G: 'S',
		Location:   [2]int{cl.Location[0] - 1, cl.Location[1]},
	}
}
func Down(cl Connection) Connection {
	return Connection{
		G: 'N',
		Location:   [2]int{cl.Location[0] + 1, cl.Location[1]},
	}
}
func Left(cl Connection) Connection {
	return Connection{
		G: 'E',
		Location:   [2]int{cl.Location[0], cl.Location[1] - 1},
	}
}
func Right(cl Connection) Connection {
	return Connection{
		G: 'W',
		Location:   [2]int{cl.Location[0], cl.Location[1] + 1},
	}
}

type offsetFunc func (Connection) Connection


func part1(lines <-chan string) int {
	var field [][]*PipeTile
	var start [2]int
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

	// Print
	for _, fieldline := range field {
		fmt.Println(fieldline)
	}

	// Find start connections
	var xbound int = len(field[0])
	var ybound int = len(field)
	var connections []TileConnection
	{
		if start[1] > 0 {
			leftloc := [2]int{start[0], start[1] - 1}
			left := field[leftloc[0]][leftloc[1]]
			switch left.tiletype {
			case '-', 'L', 'F':
				connections = append(connections, TileConnection{PipeTile: left, Connection: Connection{ 'E', leftloc}})
			}
		}

		if start[1] < xbound - 1 {
			rightloc := [2]int{start[0], start[1] + 1}
			right := field[rightloc[0]][rightloc[1]]
			switch right.tiletype {
			case '-', 'J', '7':
				connections = append(connections, TileConnection{PipeTile: right, Connection: Connection{ 'W', rightloc}})
			}
		}

		if start[0] > 0{
			uploc := [2]int{start[0] - 1, start[1]}
			up := field[uploc[0]][uploc[1]]
			switch up.tiletype {
			case '|', '7', 'F':
				connections = append(connections, TileConnection{PipeTile: up, Connection: Connection{ 'S', uploc}})
			}
		}

		if start[0] < ybound - 1 {
			downloc := [2]int{start[0] + 1, start[1]}
			down := field[downloc[0]][downloc[1]]
			switch down.tiletype {
			case '|', 'L', 'J':
				connections = append(connections, TileConnection{PipeTile: down, Connection: Connection{ 'N', downloc}})
			}
		}
	}

	if len(connections) != 2 {
		panic("Wrong number of connections")
	}

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
