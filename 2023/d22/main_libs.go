package main

import (
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

/// A Brick in 3D space.
/// The Start point will always be lower than the End point if they aren't the same height
type Brick struct {
	Start Point
	End Point
}

func ParseBrick(s string) Brick {
	sp := strings.Split(s, "~")

	return Brick{
		Start: PointFromCSV(sp[0]),
		End:   PointFromCSV(sp[1]),
	}
}

/// Determines if a brick is horizontally spread or not
func (b *Brick) IsHori() bool {
	return b.Start.Z() == b.End.Z()
}

func (b *Brick) Bottom() uint16 {
	return b.Start.Z()
}

func (b *Brick) Top() uint16 {
	return b.End.Z()
}

/// Returns the cells that are taken up by this brick
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

func (b *Brick) IsColliding(oth *Brick) bool {
	var b_candidates []Point
	if b.IsHori() {
		b_candidates = append(b_candidates, )
	} else {
		b_candidates = append(b_candidates, b.)
	}

	return true
}

type RestingBrick struct {
	Brick
	Holding map[*RestingBrick]struct{}
	RestingOn map[*RestingBrick]struct{}
}

func NewRestingBrick(b Brick) RestingBrick {
	return RestingBrick{
		Brick:     b,
		Holding:   make(map[*RestingBrick]struct{}, 3),
		RestingOn: make(map[*RestingBrick]struct{}, 3),
	}
}

type RestingBrickStack struct {
	Collision_map map[uint16][]*RestingBrick
	Resting_stack []RestingBrick
	resting_count int
}

func (r *RestingBrickStack) InsertBrick(b Brick) {
	r.Resting_stack = append(r.Resting_stack, NewRestingBrick(b))
	rbrick_ref := &r.Resting_stack[r.resting_count]
	r.resting_count++

	r.Collision_map[rbrick_ref.Top()] = append(r.Collision_map[rbrick_ref.Top()], rbrick_ref)
}

func (r *RestingBrickStack) Insert_On_Collision(b *Brick) (was_insert bool) {
	collision_candidates := r.Collision_map[b.Bottom()]
	if len(collision_candidates) < 1 {
		return false
	}

	

	return
} 