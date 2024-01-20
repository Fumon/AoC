package main

import (
	"fmt"
	"fuaoc2023/day12/u"
	"strings"
)

func main() {
	// fmt.Println(part1.Part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))
}

func part2(lines <-chan string) int {
	var memostore memoized_spring_combo_compute = memoized_spring_combo_compute{make(map[string]bool), make(map[string]int)}
	var combinations_sum int

	for line := range lines {
		sp := strings.Split(line, " ")
		

		// Initialize
		folded_run_counts := u.ParseNums(strings.Split(sp[1], ","))
		var run_counts []int = folded_run_counts[:]

		folded_game_String := sp[0]
		var game_string string = folded_game_String[:]

		for i := 0; i < 4; i++ {
			run_counts = append(run_counts, folded_run_counts...)
			game_string = game_string + "?" + folded_game_String
		}

		combos := memostore.count_combos(game_string, run_counts)
		combinations_sum += combos
	}

	return combinations_sum
}

type memoized_spring_combo_compute struct {
	cont_memo map[string]bool
	spring_memo map[string]int
}

func (m memoized_spring_combo_compute) count_combos(game_string string, run_counts []int) int {
	spring_memo := fmt.Sprint(run_counts, ":", game_string)
	if v, ok := m.spring_memo[spring_memo]; ok {
		return v
	}

	if len(run_counts) == 0 {
		if len(game_string) == 0 {
			return 1
		} else if game_string[0] == '.' || game_string[0] == '?' {
			return m.count_combos(game_string[1:], run_counts)
		} else {
			return 0
		}
	} else if len(game_string) == 0 || len(game_string) < run_counts[0] {
		return 0
	}
	

	right := 0
	if game_string[0] == '.' || game_string[0] == '?' {
		right = m.count_combos(game_string[1:], run_counts)
	}

	if game_string[0] == '.' {
		return right
	}
	
	left := 0
	can_fit := m.can_fit_contiguous(game_string[:run_counts[0]])
	if can_fit {
		if len(game_string) == run_counts[0] && len(run_counts) == 1 {
			left = 1
		} else if len(game_string) > run_counts[0] {
			if game_string[run_counts[0]] == '.' || game_string[run_counts[0]] == '?' {
				left = m.count_combos(game_string[run_counts[0] + 1:], run_counts[1:])
			}
		}
	} else {
		left = 0
	}

	m.spring_memo[spring_memo] = left + right
	return left + right
}

func (m memoized_spring_combo_compute) can_fit_contiguous(s string) bool {
	if v, ok := m.cont_memo[s]; ok {
		return v
	}

	for _, c := range s {
		if c == '.' {
			m.cont_memo[s] = false
			return false
		}
	}
	m.cont_memo[s] = true
	return true
}