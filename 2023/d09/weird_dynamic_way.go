package main

import (
	"fmt"
	"fuaoc2023/day09/u"
	"strings"
)

func predict(future bool, lines <-chan string) int {
	var sum_of_next int
	for line := range lines {

		nums := u.ParseNums(strings.Split(line, " "))

		x := -1
		bound := len(nums)
		if future {
			x = len(nums)
			bound = 0
		}

		m := NewDynamic_calc_map(nums)
		n := m.Get(dcmKey{
			k:      [2]int{x, 0},
			future: future,
			boundry_x: bound,
		})
		fmt.Println(" -> ", n)
		sum_of_next += n
	}
	return sum_of_next
}

func part1(lines <-chan string) int {
	return predict(true, lines)
}

func part2(lines <-chan string) int {
	return predict(false, lines)
}

type dcmKey struct {
	k      [2]int
	future bool
	boundry_x int
}

func (dk dcmKey) lr() dcmKey {
	lrdelta := -1
	if !dk.future {
		lrdelta = 1
	}
	return dcmKey{
		k:      [2]int{dk.k[0] + lrdelta, dk.k[1]},
		future: dk.future,
		boundry_x: dk.boundry_x,
	}
}

func (dk dcmKey) u() dcmKey {
	return dcmKey{
		k:      [2]int{dk.k[0], dk.k[1] - 1},
		future: dk.future,
		boundry_x: dk.boundry_x,
	}
}

func (dk dcmKey) ulr() dcmKey {
	lrdelta := 1
	if !dk.future {
		lrdelta = -1
	}
	return dcmKey{
		k:      [2]int{dk.k[0] + lrdelta, dk.k[1] - 1},
		future: dk.future,
		boundry_x: dk.boundry_x,
	}
}

func (dk dcmKey) dlr() dcmKey {
	lrdelta := -1
	if !dk.future {
		lrdelta = 1
	}
	return dcmKey{
		k:      [2]int{dk.k[0] + lrdelta, dk.k[1] + 1},
		future: dk.future,
		boundry_x: dk.boundry_x,
	}
}

type Dynamic_calc_map map[[2]int]int

func NewDynamic_calc_map(a []int) Dynamic_calc_map {
	out := make(Dynamic_calc_map)
	for i, v := range a {
		out[[2]int{i, 0}] = v
	}
	return out
}

func (d Dynamic_calc_map) calc(key dcmKey) (out int) {
	if key.k[1] == 0 {
		out = d[key.k]
		return
	}

	vup, ok := d[key.u().k]
	if !ok {
		vup = d.calc(key.u())
	}

	vulr, ok := d[key.ulr().k]
	if !ok {
		vulr = d.calc(key.ulr())
	}

	out = vulr - vup
	d[key.k] = out
	return
}

func (d Dynamic_calc_map) Get(key dcmKey) int {
	v, ok := d[key.k]
	if !ok {
		klr := key.lr()
		vlr := d.calc(klr)
		if vlr == 0 {
			var nonzero bool
			for i := 1; i < 4 && klr.k[0] != key.boundry_x; i, klr = i+1, klr.lr() {
				if d.calc(klr) != 0 {
					nonzero = true
					break
				}
			}
			if !nonzero {
				return 0
			}
		}

		vdlr := d.Get(key.dlr())
		v = vdlr + vlr
		d[key.k] = v

		fmt.Print(" ", v)
	}
	return v
}
