package main

import (
	"fuaoc2023/day17/u"
	"testing"
)


func TestPart1(t *testing.T) {
	expected := 102
	result := Part1(u.Linewisefile_slice("testinput"))
	if result != expected {
		t.Error("Part1 failed. Expected ", expected, " got ", result)
	}
}

func TestPart2(t *testing.T) {
	expected := 94
	result := Part2(u.Linewisefile_slice("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}

	expected = 71
	result = Part2(u.Linewisefile_slice("testinput2"))
	if result != expected {
		t.Error("Part1 failed on testinput2. Expected ", expected, " got ", result)
	}
}