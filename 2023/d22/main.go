package main

import (
	"cmp"
	"container/heap"
	"fmt"
	"fuaoc2023/day22/u"
	"slices"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {
	var bricks []Brick
	var numbricks int
	var max_x, max_y uint16
	var brick_columns = map[[2]uint16]BrickStack{}

	for line := range lines {
		new_brick := ParseBrick(line)
		max_x = max(new_brick.Start[0], new_brick.End[0], max_x)
		max_y = max(new_brick.Start[1], new_brick.End[1], max_y)

		bricks = append(bricks, new_brick)
		brick_ref := &bricks[numbricks]
		numbricks++

		var entries []Point
		if brick_ref.IsHori() {
			entries = brick_ref.Cubes()
		} else {
			entries = []Point{brick_ref.Start}
		}

		for _, p := range entries {
			coord := [2]uint16{p.X(), p.Y()}
			brick_columns[coord] = append(brick_columns[coord], brick_ref)
		}
	}

	// Sort by Z
	var per_brick_positions = make(map[*Brick]Brick_s_Positions, len(bricks))
	for coord := range brick_columns {
		brickstack := brick_columns[coord]
		slices.SortFunc(brickstack, func(a, b *Brick) int { return cmp.Compare(a.Bottom(), b.Bottom()) })
		for position, brickref := range brickstack {
			per_brick_positions[brickref] = append(per_brick_positions[brickref], Brick_Position{
				Position: position,
				Array:    &brickstack,
			})
		}
	}

	// Min heap over
	// Max Position index over All brick positions for that brick
	var position_heap = NewBrickPositionHeap(per_brick_positions)

	var ordered_layers = map[int][]*Brick{}
	for position_heap.Len() > 0 {
		nbrick := heap.Pop(&position_heap).(*Brick)
		nl := position_heap.BPs[nbrick].Highest_Position()
		ordered_layers[nl] = append(ordered_layers[nl], nbrick)
	}

	var floorbrick = Brick{
		Start: [3]uint16{0, 0, 0},
		End:   [3]uint16{max_x, max_y, 0},
	}
	var floor = &RestingBrick{
		Top:       0,
		RestingOn: map[*RestingBrick]struct{}{},
		HoldingUp: map[*RestingBrick]struct{}{},
	}
	var resting_bricks = make(map[*Brick]*RestingBrick, len(bricks)+1)
	resting_bricks[&floorbrick] = floor

	var toplayer = make([][]*RestingBrick, max_x+1)
	for x := 0; x < int(max_x+1); x++ {
		yrow := make([]*RestingBrick, max_y+1)
		for y := 0; y < int(max_y); y++ {
			yrow[y] = floor
		}
		toplayer[x] = yrow
	}

	for _, brickref := range ordered_layers[0] {
		var ntop uint16 = 1
		if !brickref.IsHori() {
			ntop += brickref.Top() - brickref.Bottom()
		}
		nrestingbrick := &RestingBrick{
			Top:       ntop,
			RestingOn: map[*RestingBrick]struct{}{floor: {}},
			HoldingUp: make(map[*RestingBrick]struct{}),
		}
		floor.HoldingUp[nrestingbrick] = struct{}{}
		resting_bricks[brickref] = nrestingbrick
		for _, v := range brickref.Columns() {
			toplayer[v.X()][v.Y()] = nrestingbrick
		}
	}

	var carryover []*Brick
	for i := 1; ; i++ {
		vvv, ok := ordered_layers[i]
		if !ok {
			break
		}

		var unresolved = append(carryover, vvv...)
		var unresolved_len = len(unresolved)
		for len(unresolved) > 0 {
			var layer = unresolved
			unresolved = []*Brick{}
		brickref_loop:
			for _, brickref := range layer {
				bps := position_heap.BPs[brickref]
				var highest_top uint16
				var resting_on = map[*RestingBrick]struct{}{}
				for _, bp := range bps {
					var underbrick *Brick
					if bp.Position == 0 {
						underbrick = &floorbrick
					} else {
						underbrick = (*bp.Array)[bp.Position-1]
					}
					rb, ok := resting_bricks[underbrick]
					if !ok {
						unresolved = append(unresolved, brickref)
						continue brickref_loop
					}

					if rb.Top > highest_top {
						highest_top = rb.Top
						resting_on = map[*RestingBrick]struct{}{rb: {}}
					} else if rb.Top == highest_top{
						resting_on[rb] = struct{}{}
					}
				}

				var ntop = highest_top + 1
				if !brickref.IsHori() {
					ntop += brickref.Top() - brickref.Bottom()
				}
				// Insert!
				nrestingbrick := &RestingBrick{
					Top:       ntop,
					RestingOn: resting_on,
					HoldingUp: make(map[*RestingBrick]struct{}),
				}
				resting_bricks[brickref] = nrestingbrick

				for r := range resting_on {
					r.HoldingUp[nrestingbrick] = struct{}{}
				}

				for _, p := range brickref.Columns() {
					toplayer[p.X()][p.Y()] = nrestingbrick
				}
			}
			if len(unresolved) == unresolved_len {
				carryover = unresolved
				break
			} else {
				unresolved_len = len(unresolved)
			}
		}
	}


	// Count
	var output int
	for _, resting_brick := range resting_bricks {
		if resting_brick == floor {
			continue
		}
		if len(resting_brick.HoldingUp) == 0 {
			output += 1
			continue
		}

		var invalid bool
		for hu := range resting_brick.HoldingUp {
			if len(hu.RestingOn) < 2 {
				invalid = true
				break
			}
		}
		if !invalid {
			output += 1
		}
	}

	return output
}

func Part2(lines <-chan string) int {
	return 0
}
