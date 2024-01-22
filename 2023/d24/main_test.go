package main

import (
	"fuaoc2023/day24/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 2
	result := Part1(u.Linewisefile_chan("testinput"), [2]int{7, 27})
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}

func TestPart2(t *testing.T) {
	expected := 0
	result := Part2(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}
