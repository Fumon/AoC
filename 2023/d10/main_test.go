package main

import (
	"fuaoc2023/day10/u"
	"testing"
)

func TestPart1_Testinputs(t *testing.T) {
	u.Assert(t, 4, part1(u.Linewisefile_chan("testinput1")))
	u.Assert(t, 8, part1(u.Linewisefile_chan("testinput2")))
}

func TestPart2_Testinputs(t *testing.T) {
	u.Assert(t, 4, part2(u.Linewisefile_chan("part2input1")))
	u.Assert(t, 8, part2(u.Linewisefile_chan("part2input2")))
	u.Assert(t, 10, part2(u.Linewisefile_chan("part2input3")))
}