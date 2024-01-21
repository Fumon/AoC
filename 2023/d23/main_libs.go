package main

type Vertex struct {
	coords [2]int
	out []*Edge
}

type Edge struct {
	From, To *Vertex
	Length int
}