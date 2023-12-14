package main

import (
	"fuaoc2023/day11/u"
	"testing"
)


func TestPart1(t *testing.T) {
	u.Assert(t, 374, part1(u.Linewisefile_chan("testinput")))
}