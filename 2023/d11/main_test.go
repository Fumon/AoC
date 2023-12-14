package main

import (
	"fuaoc2023/day11/u"
	"testing"
)


func TestPart1(t *testing.T) {
	u.Assert(t, 374, part1(u.Linewisefile_chan("testinput")))
}

func TestGrowthFactor(t *testing.T) {
	u.Assert(t, 1030, sum_shortest_paths_with_factor(u.Linewisefile_chan("testinput"), 10))
	u.Assert(t, 8410, sum_shortest_paths_with_factor(u.Linewisefile_chan("testinput"), 100))
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