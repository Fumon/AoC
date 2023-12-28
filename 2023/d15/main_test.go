package main

import (
	"fuaoc2023/day15/u"
	"testing"
)

func Assert[T comparable](t *testing.T, e, r T) {
	if e != r {
		t.Error("Expected ", e, " but got ", r)
	}
}

func TestPart1(t *testing.T) {
	Assert(t, 52, Part1(u.Linewisefile_chan("testinput1")))
	Assert(t, 1320, Part1(u.Linewisefile_chan("testinput2")))
}

func TestPart2(t *testing.T) {
	Assert(t, 145, Part2(u.Linewisefile_chan("testinput2")))
}