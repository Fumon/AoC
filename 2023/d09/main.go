package main

import (
	"fmt"
	"fuaoc2023/day09/u"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
}

func part1(lines <-chan string) int {
	var sum_of_next int
	for line := range lines {

		nums := u.ParseNums(strings.Split(line, " "))
		m := NewDynamic_calc_map(nums)
		n :=  m.Get(dcmKey{len(nums), 0})
		sum_of_next += n
	}
	return sum_of_next
}

type dcmKey [2]int
type Dynamic_calc_map map[dcmKey]int

func NewDynamic_calc_map(a []int) Dynamic_calc_map {
	out := make(Dynamic_calc_map)
	for i, v := range a {
		out[[2]int{i, 0}] = v
	}
	return out
}

func (d Dynamic_calc_map) calc(key dcmKey) (out int) {
	if key[1] == 0 {
		out = d[key]
		return
	}

	kupright := dcmKey{key[0] + 1, key[1] - 1}
	kup := dcmKey{key[0], key[1] - 1}
	
	
	vupright := d[kupright]
	vup, ok := d[kup]
	if !ok {
		vup = d.calc(kup)
	}

	out = vupright - vup
	d[key] = out
	return
}

func (d Dynamic_calc_map) Get(key dcmKey) int {
	v, ok := d[key]
	if !ok {
		kleft := dcmKey{key[0] - 1, key[1]}
		kleftdown := dcmKey{key[0] - 1, key[1] + 1}
		if key[0] == 0 {
			panic("x is zero")
		}
		
		vleft := d.calc(kleft)
		if vleft == 0 && d.calc(dcmKey{kleft[0] - 1, kleft[1]}) == 0 {
			return 0
		}
		vleftdown := d.Get(kleftdown)
		v = vleftdown + vleft
		d[key] = v 
	}
	return v
}