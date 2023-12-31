package main

import (
	"fuaoc2023/day19/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 19114
	result := Part1(u.Linewisefile_chan("testinput"))
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
