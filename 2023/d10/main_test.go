package main

import (
	"fuaoc2023/day10/u"
	"testing"
)

func TestPart1_Testinput1(t *testing.T) {
	u.Assert(t, 4, part1(u.Linewisefile_chan("testinput1")))
	u.Assert(t, 8, part1(u.Linewisefile_chan("testinput2")))
}