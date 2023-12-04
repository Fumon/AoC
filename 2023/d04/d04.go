package main

import (
	"fmt"
	"fuaoc2023/day04/u"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Total Points:\n", part1(u.Linewisefile("input")))
	fmt.Println("Total Cards: \n", part2(u.Linewisefile("input")))
}

func splitToNums(line string) (cardnumber int, output [][]string) {
	cardnum_and_numbers := strings.Split(line, ": ")
	cardnum_split := strings.Split(cardnum_and_numbers[0], " ")
	var err error
	if cardnumber, err = strconv.Atoi(cardnum_split[len(cardnum_split) - 1]); err != nil {
		panic(err)
	}
	sections := strings.Split(cardnum_and_numbers[1], " | ")
	output = [][]string{
		strings.Split(sections[0], " "),
		strings.Split(sections[1], " "),
	}

	return
}

func part1(lines <-chan string) (pointtotal int) {
	var maxmatch int
	for line := range lines {
		_, numsections := splitToNums(line)
		prizeset := make(map[string]struct{}, len(numsections[0]))
		for _, p := range numsections[0] {
			if len(p) > 0 {
				prizeset[strings.Trim(p, " ")]=struct{}{}
			}
		}

		var matchcount int
		for _, t := range numsections[1] {
			if _, match := prizeset[t]; match {
				matchcount++
			}
		}
		if matchcount > 0 {
			pointtotal += 1 << (matchcount - 1)
			maxmatch = max(maxmatch, matchcount)
		}
	}
	fmt.Println("Maximum match ", maxmatch)
	return 
}

func part2(lines <-chan string) (cardtotal int) {
	var multipliers = &MultiplierStack{}
	for line := range lines {
		_, numsections := splitToNums(line)
		prizeset := make(map[string]struct{}, len(numsections[0]))
		for _, p := range numsections[0] {
			if len(p) > 0 {
				prizeset[strings.Trim(p, " ")]=struct{}{}
			}
		}

		var matchcount int
		for _, t := range numsections[1] {
			if _, match := prizeset[t]; match {
				matchcount++
			}
		}

		multipliers.Win(matchcount)
		cardtotal += 1 * multipliers.Pop()
		// fmt.Println(multipliers)
	}
	return 
}


type MultiplierStack struct {
	ring [11]int
	cur int
}

func (m *MultiplierStack) wrap(index int) int {
	return index % 11
}

func (m *MultiplierStack) Pop() (out int) {
	// Retrieve
	out = m.get()
	
	// Advance
	m.ring[m.cur] = 0
	m.cur = m.wrap(m.cur + 1)

	return
}

func (m *MultiplierStack) get() int {
	return m.ring[m.cur] + 1
}

func (m *MultiplierStack) Win(count int) {
	multiplier := m.get()
	for i := 0; i < count; i++ {
		m.ring[m.wrap(m.cur + 1 + i)] += 1 * multiplier
	}
}