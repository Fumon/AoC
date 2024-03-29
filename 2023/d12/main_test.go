package main

import (
	"fuaoc2023/day12/part1"
	"fuaoc2023/day12/u"
	"testing"
)

func TestPart1(t *testing.T) {
	u.Assert(t, 21, part1.Part1(u.Linewisefile_chan("testinput")))
	u.Assert(t, 7718, part1.Part1(u.Linewisefile_chan("input")))
}

func TestPart2(t *testing.T) {
	u.Assert(t, 525152, part2(u.Linewisefile_chan("testinput")))
	u.Assert(t, 128741994134728, part2(u.Linewisefile_chan("input")))
}

func BenchmarkPart1TestInput(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1.Part1(u.Linewisefile_chan("testinput"))
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1.Part1(u.Linewisefile_chan("input"))
	}
}