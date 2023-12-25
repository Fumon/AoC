package main

import (
	"fuaoc2023/day13/u"
	"testing"
)


func TestPart1(t *testing.T) {
	u.Assert(t, 405, Part1(u.Linewisefile_chan("testinput")))
}

func TestPart2(t *testing.T) {
	u.Assert(t, 400, Part2(u.Linewisefile_chan("testinput")))
}