package main

import (
	"container/heap"
	"strconv"
	"strings"
)

type Point [3]uint16

func PointFromCSV(s string) (p Point) {
	for i, cos := range strings.Split(s, ",") {
		conv, err := strconv.ParseUint(cos, 10, 16)
		if err != nil {
			panic(err)
		}
		p[i] = uint16(conv)
	}
	return
}

func (p *Point) X() uint16 {
	return p[0]
}

func (p *Point) Y() uint16 {
	return p[1]
}

func (p *Point) Z() uint16 {
	return p[2]
}

func (p *Point) EqXY(oth *Point) bool {
	return p[0] == oth[0] && p[1] == oth[1]
}

// / A Brick in 3D space.
// / The Start point will always be lower than the End point if they aren't the same height
type Brick struct {
	Start Point
	End   Point
}

func ParseBrick(s string) Brick {
	sp := strings.Split(s, "~")

	return Brick{
		Start: PointFromCSV(sp[0]),
		End:   PointFromCSV(sp[1]),
	}
}

// / Determines if a brick is horizontally spread or not
func (b *Brick) IsHori() bool {
	return b.Start.Z() == b.End.Z()
}

func (b *Brick) Bottom() uint16 {
	return b.Start.Z()
}

func (b *Brick) Top() uint16 {
	return b.End.Z()
}

// / Returns the cells that are taken up by this brick
func (b *Brick) Cubes() (o []Point) {
	o = append(o, b.Start)
	var i int
	for i = 0; i < 3; i++ {
		if b.Start[i] != b.End[i] {
			break
		}
	}
	if i == 3 {
		// Brick is a cube
		return
	}
	var p = b.Start
	for p[i] = p[i] + 1; p[i] <= b.End[i]; p[i]++ {
		o = append(o, p)
	}
	return
}

func (b *Brick) Columns() (o []Point) {
	if b.IsHori() {
		return b.Cubes()
	} else {
		return []Point{b.Start}
	}
}

type BrickStack []*Brick

type Brick_Position struct {
	Position int
	Array    *BrickStack
}

type Brick_s_Positions []Brick_Position
func (bsp Brick_s_Positions) Highest_Position() int {
	var high int = bsp[0].Position
	for _, bp := range bsp {
		if bp.Position > high {
			high = bp.Position
		}
	}
	return high
}

type RestingBrick struct {
	Top uint16
	RestingOn map[*RestingBrick]struct{}
	HoldingUp map[*RestingBrick]struct{}
}

type Brick_Position_Heap struct {
	BPs  map[*Brick]Brick_s_Positions
	Heap []*Brick

}

func NewBrickPositionHeap(bps map[*Brick]Brick_s_Positions) Brick_Position_Heap {
	out := Brick_Position_Heap{
		BPs:  bps,
		Heap: make([]*Brick, 0, len(bps)),
	}
	heap.Init(&out)
	for k := range bps {
		heap.Push(&out, k)
	}
	return out
}

// Len is the number of elements in the collection.
func (bph *Brick_Position_Heap) Len() int {
	return len(bph.Heap)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//   - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//   - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (bph *Brick_Position_Heap) Less(i int, j int) bool {
	a := bph.BPs[bph.Heap[i]].Highest_Position()
	b := bph.BPs[bph.Heap[j]].Highest_Position()
	return a < b
}

// Swap swaps the elements with indexes i and j.
func (bph *Brick_Position_Heap) Swap(i int, j int) {
	bph.Heap[i], bph.Heap[j] = bph.Heap[j], bph.Heap[i]
}

func (bph *Brick_Position_Heap) Push(x any) {
	bph.Heap = append(bph.Heap, x.(*Brick))
}

func (bph *Brick_Position_Heap) Pop() any {
	old := bph.Heap
	n := len(old)
	x := old[n-1]
	bph.Heap = old[0 : n-1]
	return x
}
