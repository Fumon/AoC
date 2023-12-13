package main

import "fmt"

func main() {

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
	tiletype int
}

func (p PipeTile) String() string {
	switch p.tiletype {
	case 0:
		return "|"
	case 1:
		return "-"
	case 2:
		return "L"
	case 3:
		return "J"
	case 4:
		return "7"
	case 5:
		return "F"
	case 6:
		return "."
	case 7:
		return "S"
	default:
		panic("FDSAFDSAFDSAFDSAFASDFDSA")
	}
}

func NewPipeTile(b byte) (isStart bool, p PipeTile) {
	switch b {
	case '|':
		p.tiletype = 0
	case '-':
		p.tiletype = 1
	case 'L':
		p.tiletype = 2
	case 'J':
		p.tiletype = 3
	case '7':
		p.tiletype = 4
	case 'F':
		p.tiletype = 5
	case '.':
		p.tiletype = 6
	case 'S':
		p.tiletype = 7
		isStart = true
	}
	return
}

type Connection struct {
	PipeTile
	Connection byte
	Location [2]int
}

func part1(lines <-chan string) int {
	var field [][]PipeTile
	var start [2]int
	var startfound bool
	var lineindex int
	for line := range lines {
		var fieldline []PipeTile
		for i, v := range []byte(line) {
			isStart, newTile := NewPipeTile(v)
			if isStart {
				start = [2]int{lineindex, i}
				startfound = false
			}
			fieldline = append(fieldline, newTile)
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

	var connections []Connection
	{
		// Find start connections
		leftloc := [2]int{start[0], start[1] - 1}
		left := field[leftloc[0]][leftloc[1]]
		switch left.tiletype {
		case 1 | 2 | 5 :
			connections = append(connections, Connection{left, 'E', leftloc})
		}

		rightloc := [2]int{start[0], start[1] - 1}
		left := field[leftloc[0]][leftloc[1]]
		switch left.tiletype {
		case 1 | 2 | 5 :
			connections = append(connections, Connection{left, 'E', leftloc})
		}


	}
	return 4
}
