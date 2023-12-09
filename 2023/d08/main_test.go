package main

import (
	"fuaoc2023/day08/u"
	"testing"
)

func TestSliceIndex(t *testing.T) {
	expected := "BBB"
	result := "(BBB, CCC)"[1:4]
	if expected != result {
		t.Error("Expected ", expected, " got ", result)
	}
}

func TestPart1(t *testing.T) {
	expected := 2
	result := part1(u.Linewisefile_chan("testinput1"))
	if expected != result {
		t.Error("Part1 Input1 Fail: e", expected, " != ", result)
	}
	expected = 6
	result = part1(u.Linewisefile_chan("testinput2"))
	if expected != result {
		t.Error("Part1 Input2 Fail: e", expected, " != ", result)
	}
}

func TestPart2(t *testing.T) {
	expected := 6
	result := part2(u.Linewisefile_chan("testinput3"))
	if expected != result {
		t.Error("Part1 Input1 Fail: e", expected, " != ", result)
	}
}