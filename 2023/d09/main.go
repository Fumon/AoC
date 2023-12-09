package main

import (
	"fmt"
	"fuaoc2023/day09/u"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))
}

func predict(future bool, lines <-chan string) int {
	var sum_of_next int
	for line := range lines {

		
		nums := u.ParseNums(strings.Split(line, " "))

		x := -1
		if future {
			x = len(nums)
		}

		m := NewDynamic_calc_map(nums)
		n :=  m.Get(dcmKey{
			k:       [2]int{x, 0},
			future: future,
		})
		fmt.Println(n)
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
	k [2]int
	future bool
}

func (dk dcmKey) lr() dcmKey {
	lrdelta := -1
	if !dk.future {
		lrdelta = 1
	}
	return dcmKey{
		k:       [2]int{dk.k[0] + lrdelta, dk.k[1]},
		future: dk.future,
	}
}


func (dk dcmKey) u() dcmKey {
	return dcmKey {
		k:       [2]int{dk.k[0], dk.k[1] - 1},
		future: dk.future,
	}
}

func (dk dcmKey) ulr() dcmKey {
	lrdelta := 1
	if !dk.future {
		lrdelta = -1
	}
	return dcmKey{
		k:       [2]int{dk.k[0] + lrdelta, dk.k[1] - 1},
		future: dk.future,
	}
}

func (dk dcmKey) dlr() dcmKey {
	lrdelta := -1
	if !dk.future {
		lrdelta = 1
	}
	return dcmKey{
		k:       [2]int{dk.k[0] + lrdelta, dk.k[1] + 1},
		future: dk.future,
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

	vupright := d[key.ulr().k]

	out = vupright - vup
	d[key.k] = out
	return
}

func (d Dynamic_calc_map) Get(key dcmKey) int {
	v, ok := d[key.k]
	if !ok {
		vleft := d.calc(key.lr())
		if vleft == 0 && d.calc(key.lr().lr()) == 0 {
			return 0
		}

		vleftdown := d.Get(key.dlr())
		v = vleftdown + vleft
		d[key.k] = v 
	}
	return v
}