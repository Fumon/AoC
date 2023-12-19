package main

import (
	"fuaoc2023/day12/u"
	"testing"
)

func TestPart1(t *testing.T) {
	u.Assert(t, 21, part1(u.Linewisefile_chan("testinput")))
}