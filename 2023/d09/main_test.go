package main

import (
	"fuaoc2023/day09/u"
	"testing"
)

func TestPart1_Testinput(t *testing.T) {
	u.Assert(t, 114, part1(u.Linewisefile_chan("testinput")))
}