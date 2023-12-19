package main

import (
	"fmt"
	"fuaoc2023/day12/u"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
}

func part1(lines <-chan string) int {
	var combinations_sum int
	for line := range lines {
		sp := strings.Split(line, " ")
		
		run_counts := u.ParseNums(strings.Split(sp[1], ","))
		gs := NewGamespace([]byte(sp[0]))

		// // Explore solution space
		// var permut Permut_Instance
		// copy(permut, gs.Input)

		var linesum int
		recursive_permutation_explore(gs.Input, gs.Unknown_indicies, run_counts, &linesum)
		combinations_sum += linesum
	}

	return combinations_sum
}

func recursive_permutation_explore(permut Permut_Instance, indicies []int, run_counts []int, combo_sum *int) {
	if len(indicies) == 0 {
		return
	}
	mindex := indicies[0]

	permut[mindex] = '#'
	res := validate_permutation(permut, run_counts)
	switch res {
	case 1:
		*combo_sum++
	case 0:
		recursive_permutation_explore(permut, indicies[1:], run_counts, combo_sum)
	}
	permut[mindex] = '.'
	res = validate_permutation(permut, run_counts)
	switch res {
	case 1:
		*combo_sum++
	case 0:
		recursive_permutation_explore(permut, indicies[1:], run_counts, combo_sum)
	}
	permut[mindex] = '?'
	return
}

type Permut_Instance []byte

type Gamespace struct {
	Input []byte
	Unknown_indicies []int
}

func NewGamespace(input []byte) Gamespace {
	var unknown_indicies []int
	for i, b := range input {
		if b == '?' {
			unknown_indicies = append(unknown_indicies, i)
		}
	}

	return Gamespace{
		Input:            input,
		Unknown_indicies: unknown_indicies,
	}
}


// 1 is pass
func validate_permutation(permut Permut_Instance, run_counts []int) (int) {
	_, _, rc := count_runs(permut)

	if len(rc) != len(run_counts) {
		return 0
	}

	for j := 0; j < len(run_counts); j++ {
		if rc[j] != run_counts[j] {
			return 0
		}
	}
	return 1
}

func count_runs(permut Permut_Instance) (pre_runcounts []int, post_runcounts []int, run_counts []int) {
	var curcount int
	var i int

	spit := func(sl *[]int) {
		if curcount > 0 {
			*sl = append(*sl, curcount)
			run_counts = append(run_counts, curcount)
			curcount = 0
		}
	}

	LOOP:
	for i < len(permut) {
		switch permut[i] {
		case '.':
			spit(&pre_runcounts)
		case '?':
			break LOOP
		case '#':
			curcount++
		}
		i++
	}
	spit(&pre_runcounts)

	for i < len(permut) {
		switch permut[i] {
		case '.':
			spit(&post_runcounts)
		case '?':
			spit(&post_runcounts)
		case '#':
			curcount++
		}
		i++
	}
	spit(&post_runcounts)

	return
}