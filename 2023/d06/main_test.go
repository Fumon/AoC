package main

import (
	"fuaoc2023/day06/u"
	"testing"
)

func TestPart1(t *testing.T) {
	expected := 288
	received := part1(form_input(u.Linewisefile_chan("testinput")))
	if received != expected {
		t.Errorf("Part1 Failed: Expected %v and received %v", expected, received)
	}
}