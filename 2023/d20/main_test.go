package main

import (
	"fuaoc2023/day20/u"
	"testing"
)

func TestPart1(t *testing.T) {
	var g = map[string]int{
		"testinput1": 32000000,
		"testinput2": 11687500,
	}
	for f, expected := range g {
	result := Part1(u.Linewisefile_chan(f))
	if result != expected {
		t.Error("Part1 failed on ",f, ". Expected ", expected, " got ", result)
	}
	}
}

func TestPart2(t *testing.T) {
	expected := 237878264003759
	result := Part2(u.Linewisefile_chan("input"))
	if result != expected {
		t.Error("Part2 failed on input. Expected ", expected, " got ", result)
	}
}


func TestPulseQueue(t *testing.T) {
	var pulse_queue PulseQueue
	
	if _, ok := pulse_queue.Pop(); ok != false {
		t.Error("empty queue fail")
	}

	p1 := Pulse{}
	pulse_queue.Push(p1)

	if pop1, ok := pulse_queue.Pop(); pop1 != p1 || ok != true {
		t.Error("pop error")
	}
}

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Part1(u.Linewisefile_chan("input"))
	}
}