package main

import (
	"fuaoc2023/day05/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 35
	received := part1(u.Linewisefile_chan("testinput"))
	if received != expected {
		t.Errorf("Part1 Failed: Expected %v and received %v", expected, received)
	}
}

func TestPart2(t *testing.T) {
	expected := 46
	received := part2(u.Linewisefile_chan("testinput"))
	if received != expected {
		t.Errorf("Part2 Failed: Expected %v and received %v", expected, received)
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1(u.Linewisefile_chan("input"))
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2(u.Linewisefile_chan("input"))
	}
}