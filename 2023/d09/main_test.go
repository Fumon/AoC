package main

import (
	"fuaoc2023/day09/u"
	"testing"
)

func TestPart1_Testinput(t *testing.T) {
	u.Assert(t, 114, part1(u.Linewisefile_chan("testinput")))
}

func TestPart2_Testinput(t *testing.T) {
	u.Assert(t, 2, part2(u.Linewisefile_chan("testinput")))
}


func TestPart1_Input(t *testing.T) {
	u.Assert(t, 1987402313, part1(u.Linewisefile_chan("input")))
}

func TestPart2_Input(t *testing.T) {
	u.Assert(t, 900, part2(u.Linewisefile_chan("input")))
}


func TestPart1_diff_Testinput(t *testing.T) {
	u.Assert(t, 114, predict_diff(true, u.Linewisefile_chan("testinput")))
}

func TestPart2_diff_Testinput(t *testing.T) {
	u.Assert(t, 2, predict_diff(false, u.Linewisefile_chan("testinput")))
}


func TestPart1_diff_Input(t *testing.T) {
	u.Assert(t, 1987402313, predict_diff(true, u.Linewisefile_chan("input")))
}

func TestPart2_diff_Input(t *testing.T) {
	u.Assert(t, 900, predict_diff(false, u.Linewisefile_chan("input")))
}