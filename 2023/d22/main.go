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
	// Sort by Z
	// Min heap over
	// Max Position index over All brick positions for that brick
	// Insert!
	floor, resting_bricks := process_bricks(lines)

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

func process_bricks(lines <-chan string) (*RestingBrick, map[*Brick]*RestingBrick) {
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
					} else if rb.Top == highest_top {
						resting_on[rb] = struct{}{}
					}
				}

				var ntop = highest_top + 1
				if !brickref.IsHori() {
					ntop += brickref.Top() - brickref.Bottom()
				}

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
	return floor, resting_bricks
}

func Part2(lines <-chan string) int {
	floor, resting_bricks := process_bricks(lines)

	// var sum_map = make(map[*RestingBrick]int, len(resting_bricks))
	// var get_disintegrate_sum func(*RestingBrick) int
	// get_disintegrate_sum = func(rb *RestingBrick) (total int) {
	// 	if val, ok := sum_map[rb]; ok {
	// 		return val
	// 	}

	// 	for oth := range rb.HoldingUp {
	// 		if len(oth.RestingOn) == 1 {
	// 			total += get_disintegrate_sum(oth) + 1
	// 		} else {
	// 			get_disintegrate_sum(oth)
	// 		}
	// 	}
	// 	sum_map[rb] = total
	// 	return
	// }

	// get_disintegrate_sum(floor)

	// var total int
	// for k, v := range sum_map {
	// 	if k != floor {
	// 		total += v
	// 	}
	// 	println(k, ": ", v)
	// }

	var total int
	for _, rb := range resting_bricks {
		if rb == floor {
			continue
		}
		frontier := rb.HoldingUp
		supports := map[*RestingBrick]struct{}{rb: {}}
		var disintegrate_sum int
		for len(frontier) > 0 {
			next_frontier := map[*RestingBrick]struct{}{}
		frontier_loop:
			for f := range frontier {
				if _, ok := supports[f]; ok {
					continue frontier_loop
				}
				for support := range f.RestingOn {
					if _, ok := supports[support]; !ok {
						continue frontier_loop
					}
				}
				disintegrate_sum += 1
				// Add to supports
				supports[f] = struct{}{}
				// Add to frontier
				for hu := range f.HoldingUp {
					next_frontier[hu] = struct{}{}
				}
			}
			frontier = next_frontier
		}

		total += disintegrate_sum
	}

	return total
}
