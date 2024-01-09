package main

import (
	"fuaoc2023/day21/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 16
	result := Part1(6, u.Linewisefile_chan("testinput"))
	if result != expected {
		t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
	}
}

func TestPart2(t *testing.T) {
	g := [][2]int {
		{6, 16},
		{10, 50},
		{50, 1594},
		{100, 6536},
		{500, 167004},
		{1000, 668697},
		{5000, 16733044},
	}
	for _, v := range g {
		expected := v[1]
		result := Part2(v[0], u.Linewisefile_chan("testinput"))
		if result != expected {
			t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
		}
	}
}

func BenchmarkTest2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Part2(1000, u.Linewisefile_chan("testinput"))
	}
}
