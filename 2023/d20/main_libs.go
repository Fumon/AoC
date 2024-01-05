package main

import "fmt"

type Pulse struct {
	source string
	destination string
	high bool
}

func (p Pulse) String() string {
	hstring := "low"
	if p.high {
		hstring = "high"
	}
	return fmt.Sprintf("%s -%s-> %s", p.source, hstring, p.destination)
}


type PulseQueue []Pulse
func (pq *PulseQueue) Pop() (Pulse, bool){
	if len(*pq) < 1 {
		return Pulse{}, false
	}

	out := (*pq)[0]
	*pq = (*pq)[1:]
	// fmt.Println(out)
	return out, true
}
func (pq *PulseQueue) Push(np Pulse) {
	*pq = append(*pq, np)
}

type Conjunction struct {
	in_connections map[string]uint64
	statemask uint64
}

type Flipflop struct {
	statemask uint64
}