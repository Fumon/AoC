package main

import (
	"fmt"
	"fuaoc2023/day08/u"
	"strings"
)

func main() {
	fmt.Println("Part1: ", part1(u.Linewisefile_chan("input")))
	fmt.Println("Part2: ", part2(u.Linewisefile_chan("input")))
}

func part1(lines <-chan string) int {

	instructions := []rune((<-lines))
	<-lines

	var mapping = make(map[uint16]*Node)
	for line := range lines {
		g := strings.Split(line, " = ")
		name := g[0]

		leftname, rightname := g[1][1:4], g[1][6:9]
		mapping[ConvertNameToInt(name)] = &Node{
			name: name,
			L:    ConvertNameToInt(leftname),
			R:    ConvertNameToInt(rightname),
		}
	}

	instructionIndex := 0
	cur := mapping[ConvertNameToInt("AAA")]
	stepcount := 0
	fmt.Println(cur.name)
	for cur.name != "ZZZ" {
		switch instructions[instructionIndex] {
		case 'L':
			cur = mapping[cur.L]
		case 'R':
			cur = mapping[cur.R]
		default:
			panic("WHAT")
		}
		instructionIndex = (instructionIndex + 1) % len(instructions)
		stepcount++
	}

	return stepcount
}

func part2(lines <-chan string) int {

	instructions := []rune((<-lines))
	<-lines

	var starts []*Node
	var mapping = make(map[uint16]*Node)
	for line := range lines {
		g := strings.Split(line, " = ")
		name := g[0]

		leftname, rightname := g[1][1:4], g[1][6:9]
		code := ConvertNameToInt(name)

		nnode := &Node{
			name: name,
			L:    ConvertNameToInt(leftname),
			R:    ConvertNameToInt(rightname),
		}

		mapping[code] = nnode

		if name[2] == 'A' {
			starts = append(starts, nnode)
		}
	}

	// Link
	for _, v := range mapping {
		v.LL = mapping[v.L]
		v.RR = mapping[v.R]
	}

	curs := make([]*Node, len(starts))
	copy(curs, starts)
	ztimeslimit := 1
	var steplimit int = 2_000_000_000
	stepsarray := make([]int, len(starts))
	for i, st := range curs {
		instructionIndex := 0
		var stepcount int = 0

		for stepsarray[i] < ztimeslimit && stepcount < steplimit {
			switch instructions[instructionIndex] {
			case 'L':
				st = st.LL
			case 'R':
				st = st.RR
			default:
				panic("WHAT")
			}

			instructionIndex = (instructionIndex + 1) % len(instructions)
			stepcount++

			if st.name[2] == 'Z' {
				stepsarray[i] = stepcount
				stepcount = 0
			}
		}

		fmt.Println(starts[i].name, ": ", stepsarray[i])
	}

	gcd := func(a, b int) int {
		for b > 0 {
			a, b = b, a % b
		}
		return a
	}

	lcm := func(list []int) int {
		a := list[0]
		for _, b := range list[1:] {
			a = (a * b) / gcd(a, b)
		}
		return a
	}



	// instructionIndex := 0
	// stepcount := 0
	// for {
	// 	var zcount int
	// 	for i := 0; i < len(starts); i++ {
	// 		switch instructions[instructionIndex] {
	// 		case 'L':
	// 			starts[i] = starts[i].LL
	// 		case 'R':
	// 			starts[i] = starts[i].RR
	// 		default:
	// 			panic("WHAT")
	// 		}
	// 		if starts[i].name[2] == 'Z' {
	// 			zcount++
	// 		}
	// 	}

	// 	instructionIndex = (instructionIndex + 1) % len(instructions)
	// 	stepcount++
	// 	if zcount == len(starts) {
	// 		break
	// 	}
	// }

	return lcm(stepsarray)
}

func ConvertNameToInt(s string) (out uint16) {
	for _, b := range []byte(s) {
		out = (out << 5) | uint16(b&0x1F)
	}
	return
}

type Node struct {
	name string
	L    uint16
	R    uint16
	LL   *Node
	RR   *Node
}
