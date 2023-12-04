package main

import (
	"fmt"
	"fuaoc2023/day04/u"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Total Points:\n", part1(u.Linewisefile_chan("input")))
	fmt.Println("Total Cards: \n", part2_but_a_slice_fewerallocs_fixedarray(u.Linewisefile_slice("input")))
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
	// var maxmatch int
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
			// maxmatch = max(maxmatch, matchcount)
		}
	}
	// fmt.Println("Maximum match ", maxmatch)
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

func part2_but_a_slice(lines []string) (cardtotal int) {
	var multipliers = &MultiplierStack{}
	for _, line := range lines {
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

func part2_but_a_slice_fewerallocs_string(lines []string) (cardtotal int) {
	// Peek to get sizes
	fstline := lines[0]
	colon := strings.Index(fstline, ":")
	pipe := strings.Index(fstline, " |")
	llen := len(fstline)

	pstart := colon + 1
	plen := (pipe - colon)
	nstart := pipe + 2
	ncount := (llen - pipe) / 3

	var multipliers = &MultiplierStack{}
	for _, line := range lines {
		prizeset := line[pstart:pstart+plen]

		var matchcount int
		for i := 0; i < ncount; i++ {
			if strings.Contains(prizeset, line[nstart + i*3:nstart+(i+1)*3]) {
				matchcount++
			}
		}

		multipliers.Win(matchcount)
		cardtotal += 1 * multipliers.Pop()
		// fmt.Println(multipliers)
	}
	return 
}

func part2_but_a_slice_fewerallocs_hash(lines []string) (cardtotal int) {
	// Peek to get sizes
	fstline := lines[0]
	colon := strings.Index(fstline, ":")
	pipe := strings.Index(fstline, " |")
	llen := len(fstline)

	pstart := colon + 1
	pcount := (pipe - colon) / 3
	nstart := pipe + 2
	ncount := (llen - pipe) / 3

	var multipliers = &MultiplierStack{}
	for _, line := range lines {
		prizeset := make(map[string]struct{}, pcount)
		for i := 0; i < pcount; i++ {
			prizeset[line[pstart + i*3:pstart + (i+1)*3]] = struct{}{}
		}

		var matchcount int
		for i := 0; i < ncount; i++ {
			if _, match := prizeset[line[nstart + i*3:nstart+(i+1)*3]]; match {
				matchcount++
			}
		}

		multipliers.Win(matchcount)
		cardtotal += 1 * multipliers.Pop()
		// fmt.Println(multipliers)
	}
	return 
}

func part2_but_a_slice_fewerallocs_fixedarray(lines []string) (cardtotal int) {
	// Peek to get sizes
	fstline := lines[0]
	colon := strings.Index(fstline, ":")
	pipe := strings.Index(fstline, " |")
	llen := len(fstline)

	pstart := colon + 1
	pcount := (pipe - colon) / 3
	nstart := pipe + 2
	ncount := (llen - pipe) / 3

	// Can contain all the ascii combinations
	var multipliers = &MultiplierStack{}
	for _, line := range lines {
		var hashset = &u.ASCIINumberHashSet{}
		for i := 0; i < pcount; i++ {
			s := []byte(line[pstart + i*3 + 1:pstart + (i+1)*3])
			hashset.Insert(s)
		}

		var matchcount int
		for i := 0; i < ncount; i++ {
			s := []byte(line[nstart + i*3 + 1:nstart+(i+1)*3])
			if hashset.Exists(s) {
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