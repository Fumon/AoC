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