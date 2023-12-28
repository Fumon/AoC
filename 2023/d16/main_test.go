package main

import (
	"fuaoc2023/day16/u"
	"testing"
)

func TestPart1(t *testing.T) {
	energized := Part1(u.Linewisefile_chan("testinput"))
	if energized != 46 {
		t.Error("Expected:", 46, "Got:", energized)
	}
}

func TestPart2(t *testing.T) {
	energized := Part2(u.Linewisefile_chan("testinput"))
	if energized != 51 {
		t.Error("Expected:", 51, "Got:", energized)
	}
}