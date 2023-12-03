package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	red = "\033[31m"
	green = "\033[32m"
	yellow = "\033[33m"
	blue = "\033[34m"
	reset = "\033[0m"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	ls := bufio.NewScanner(f)

	var partsum int

	var prevlinesym []bool
	var prevlineparts []*Part

	for ls.Scan() {

		line := ls.Text()

		var curlinesym = make([]bool, len(line))
		var curlineparts = make([]*Part, len(line))

		var prevsym int = -1000

		var curpart *Part = nil

		for ind, ch := range line {
			if ch == '.' {
				// Check for valid part
				if curpart != nil {
					// Diagonal previous line check
					if !curpart.valid {
						if prevlinesym != nil && prevlinesym[ind] {
							curpart.valid = true
							// log.Println(curpart.Val(), " by diagprevlinefwd")
						}
					}

					if curpart.valid {
						partsum += curpart.Val()
					}
					curpart = nil
				}
			} else if ch >= '0' && ch <= '9' {
				if curpart == nil {
					t := NewPart(ch, ind)
					curpart = &t
					// Check prev
					if pi := (ind - 1); pi >= 0 && ( pi == prevsym || (prevlinesym != nil && prevlinesym[pi])) {
						curpart.valid = true
						// log.Println(curpart.Val(), " by prevsym or diagprevlineback")
					}
				} else {
					curpart.Add(ch)
				}
				curlineparts[ind] = curpart
				// Check cur ind
				if !curpart.valid && (prevlinesym != nil && prevlinesym[ind]) {
					curpart.valid = true
					// log.Println(curpart.Val(), " by prevline")
				}
			} else {
				// Symbol
				curlinesym[ind] = true
				prevsym = ind

				if curpart != nil {
					curpart.valid = true
					// log.Println(curpart.Val(), " by hitsym")
					partsum += curpart.Val()
					curpart = nil
				}

				// Check previous line
				if prevlineparts != nil {
					indcheck := []int{ind}
					if ind > 0 {
						indcheck = append(indcheck, ind - 1)
					}
					if ind + 1 < len(line) {
						indcheck = append(indcheck, ind + 1)
					}
					// log.Println("Checking ", indcheck)
					for _, pind := range indcheck {
						c := prevlineparts[pind]
						if c != nil && !c.valid {
							c.valid = true
							// log.Println(c.Val(), " by symprevlinecheck {", string(ch), " ", ind, "->", indcheck, ": ", c.bounds)
							partsum += c.Val()
						}
					}
				}
			}
		}

		// New Line
		if curpart != nil && curpart.valid {
			partsum += curpart.Val()
			curpart = nil
		}

		// Print prevline
		for i, part := range prevlineparts {
			if part == nil {
				if prevlinesym[i] {
					fmt.Print(yellow + "o" + reset)
				} else {
					fmt.Print(".")
				}
			} else {
				if part.valid {
					fmt.Print(green + "#" + reset)
				} else {
					fmt.Print(red + "#" + reset)
				}
			}
		}

		prevlinesym = curlinesym
		prevlineparts = curlineparts


		fmt.Print("\n")
	}
	fmt.Println(partsum)
}

// Part type
type Part struct {
	s      []rune
	bounds []int
	valid  bool
}

func NewPart(r rune, ind int) Part {
	return Part{
		[]rune{r},
		[]int{ind, ind},
		false,
	}
}

func (p *Part) Add(r rune) {
	p.s = append(p.s, r)
	p.bounds[1] += 1
}

func (p *Part) Val() int {
	val, err := strconv.Atoi(string(p.s))
	if err != nil {
		log.Fatal(err)
	}
	return val
}
