package main

import (
	"fuaoc2023/day07/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 6440
	result := part1(u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Errorf("Part1 failed: expected %v and got %v", expected, result)
	}
}

func TestPart1_mine(t *testing.T) {
	expected := 8810
	result := part1(u.Linewisefile_chan("testinput_mine"))
	if result != expected {
		t.Errorf("Part1 failed: expected %v and got %v", expected, result)
	}
}
