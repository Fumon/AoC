package main

import (
	"fuaoc2023/day04/u"
	"testing"
)

func TestPart1OnInput(t *testing.T) {
	total := part1(u.Linewisefile("testinput"))
	if total != 13 {
		t.Errorf("Part1 should sum to 13, but got %d", total)
	}
}
func TestPart2OnInput(t *testing.T) {
	total := part2(u.Linewisefile("testinput"))
	if total != 30 {
		t.Errorf("Part2 should sum to 30, but got %d", total)
	}
}