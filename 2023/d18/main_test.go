package main

import (
	"fuaoc2023/day18/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 62
	result := Part1(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}

func TestPart2(t *testing.T) {
	expected := int64(952408144115)
	result := Part2(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}
