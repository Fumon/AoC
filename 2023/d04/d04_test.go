package main

import (
	"fuaoc2023/day04/u"
	"testing"
)

func TestPart1OnInput(t *testing.T) {
	total := part1(u.Linewisefile_chan("testinput"))
	if total != 13 {
		t.Errorf("Part1 should sum to 13, but got %d", total)
	}
}
func TestPart2OnInput(t *testing.T) {
	total := part2(u.Linewisefile_chan("testinput"))
	if total != 30 {
		t.Errorf("Part2 should sum to 30, but got %d", total)
	}

	tests := map[string]func([]string)int {
		"Part2slice": part2_but_a_slice,
		"Part2slice_fewerallocs_string": part2_but_a_slice_fewerallocs_string,
		"Part2slice_fewerallocs_hash": part2_but_a_slice_fewerallocs_hash,
		"Part2slice_fewerallocs_fixedarray": part2_but_a_slice_fewerallocs_fixedarray,
	}
	lineslice := u.Linewisefile_slice("testinput")
	for name, f := range tests {
		total := f(lineslice)
		if total != 30 {
			t.Errorf("On %s: should be 30 but got %d", name, total)
		}
	}
}

func BenchmarkChannelPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1(u.Linewisefile_chan("input"))
	}
}
func BenchmarkChannelPart2_original(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2(u.Linewisefile_chan("input"))
	}
}

func BenchmarkChannelPart2_slice(b *testing.B) {
	input := u.Linewisefile_slice("input")
	for i := 0; i < b.N; i++ {
		part2_but_a_slice(input)
	}
}

func BenchmarkChannelPart2_slice_fewerallocs_string(b *testing.B) {
	input := u.Linewisefile_slice("input")
	for i := 0; i < b.N; i++ {
		part2_but_a_slice_fewerallocs_string(input)
	}
}

func BenchmarkChannelPart2_slice_fewerallocs_hash(b *testing.B) {
	input := u.Linewisefile_slice("input")
	for i := 0; i < b.N; i++ {
		part2_but_a_slice_fewerallocs_hash(input)
	}
}

func BenchmarkChannelPart2_slice_fewerallocs_fixedarray(b *testing.B) {
	input := u.Linewisefile_slice("input")
	for i := 0; i < b.N; i++ {
		part2_but_a_slice_fewerallocs_fixedarray(input)
	}
}