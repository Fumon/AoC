package main

import (
	"fmt"
	"fuaoc2023/day16/u"
	"sync"
)

func main() {
	fmt.Println("Part1:", Part1(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {

	var contraption [][]byte
	for line := range lines {
		contraption = append(contraption, []byte(line))
	}
	var width int = len(contraption[0])
	var height int = len(contraption)

	var bounds = twodrect{
		ulc: twodvec{},
		lrc: twodvec{
			x: width - 1,
			y: height - 1,
		},
	}

	var lightmap = make(map[[2]twodvec]uint8)
	var lightmap_lck sync.RWMutex

	var beamwait sync.WaitGroup

	var new_beam func(twodvec, twodvec)
	new_beam = func(coords, dir twodvec) {
		for bounds.inside(coords) {
			lightmap_lck.RLock()
			if lightmap[[2]twodvec{coords, dir}] > 0 {
				lightmap_lck.RUnlock()
				break
			}
			lightmap_lck.RUnlock()
			lightmap_lck.Lock()
			lightmap[[2]twodvec{coords, dir}]++
			lightmap_lck.Unlock()

			// Determine where to go next
			ch := contraption[coords.y][coords.x]
			switch ch {
			case '.':
			case '\\':
				dir = dir.reflect_y_equals_x()
			case '/':
				dir = dir.reflect_y_equals_neg_x()
			case '|':
				if dir.x != 0 {
					splitdir := dir.rot_clockwise()
					newcoord := coords
					newcoord.add(splitdir)
					beamwait.Add(1)
					go new_beam(newcoord, splitdir)

					dir = dir.rot_counterclockwise()
				}
			case '-':
				if dir.y != 0 {
					splitdir := dir.rot_clockwise()
					newcoord := coords
					newcoord.add(splitdir)
					beamwait.Add(1)
					go new_beam(newcoord, splitdir)

					dir = dir.rot_counterclockwise()
				}
			}
			coords.add(dir)
		}
		beamwait.Done()
	}

	beamwait.Add(1)
	go new_beam(twodvec{}, twodvec{1, 0})

	beamwait.Wait()

	var energized_counts = make(map[twodvec]int)
	for k := range lightmap {
		energized_counts[k[0]]++
	}

	return len(energized_counts)
}

type twodvec struct {
	x, y int
}

func (td *twodvec) add(oth twodvec) {
	td.x += oth.x
	td.y += oth.y
}

func (td twodvec) rot_clockwise() twodvec {
	return twodvec{
		x: td.y,
		y: -td.x,
	}
}

func (td twodvec) rot_counterclockwise() twodvec {
	return twodvec{
		x: -td.y,
		y: td.x,
	}
}

func (td twodvec) reflect_y_equals_x() twodvec {
	return twodvec{
		x: td.y,
		y: td.x,
	}
}

func (td twodvec) reflect_y_equals_neg_x() twodvec {
	return twodvec{
		x: -td.y,
		y: -td.x,
	}
}

type twodrect struct {
	ulc twodvec
	lrc twodvec
}

func (rect *twodrect) inside(pt twodvec) bool {
	return rect.ulc.x <= pt.x && pt.x <= rect.lrc.x && rect.ulc.y <= pt.y && pt.y <= rect.lrc.y
}
