package main

import "fmt"

type Hailstone struct {
	Start [3]float64
	SpeedVector [3]float64
}

func (h Hailstone) String() string {
	return fmt.Sprintf("%v, %v, %v @ %v, %v, %v", h.Start[0], h.Start[1], h.Start[2], h.SpeedVector[0], h.SpeedVector[1], h.SpeedVector[2])
}

func parseHailstone(line string) Hailstone {
	var h Hailstone
	fmt.Sscanf(line, "%f, %f, %f @ %f, %f, %f", &h.Start[0], &h.Start[1], &h.Start[2], &h.SpeedVector[0], &h.SpeedVector[1], &h.SpeedVector[2])
	return h
}