package main

import (
	"fuaoc2023/day14/u"
	"testing"
)

func Assert[T comparable](t *testing.T, e, r T) {
	if e != r {
		t.Error("Expected ", e, " but got ", r)
	}
}
func TestPart1(t *testing.T) {
	Assert(t, 136, Part1(u.Linewisefile_chan("testinput")))
}