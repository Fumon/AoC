package main

import (
	"fmt"
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
	t.SkipNow()
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
		} else {
			fmt.Println("Pass", v[0], "steps")
		}
	}
}

func TestPart2Assumptions(t *testing.T) {
	g := [][2]int {
		{196, 34700},
		{588, 312300},
		{982, 867500},
		{1375, 1700300},
		{1768, 2810700},
	}
	for _, v := range g {
		expected := v[1]
		result := Part2(v[0], u.Linewisefile_chan("input"))
		if result != expected {
			t.Error("Part1 failed on testinput. Expected ", expected, " got ", result)
		} else {
			fmt.Println("Pass", v[0], "steps")
		}
	}
}


func BenchmarkTest2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Part2(5000, u.Linewisefile_chan("input"))
	}
}
