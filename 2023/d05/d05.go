package main

import (
	"fmt"
	"fuaoc2023/day05/u"
	"slices"
	"strings"
)

func main() {
	fmt.Println(part1(u.Linewisefile_chan("input")))
	fmt.Println(part2(u.Linewisefile_chan("input")))
}

func part1(lines <-chan string) int {

	seeds_strings := strings.Split(strings.Split((<-lines), ": ")[1], " ")
	cur_vals := u.ParseNums(seeds_strings)

	// Blank
	<-lines

	// Map Loop
	for {
		_, ok := <-lines
		if !ok {
			break
		}

		var new_vals []int
		for mapline := range lines {
			// Blank line, end of map
			if len(mapline) == 0 {
				break
			}

			// Parse
			ma := u.ParseNums(strings.Split(mapline, " "))
			dest_start := ma[0]
			source_start := ma[1]
			width := ma[2]
			smax := source_start + width
			diff := dest_start - source_start

			// Scan
			var ncurs []int
			for _, v := range cur_vals {
				if v >= source_start && v < smax {
					new_vals = append(new_vals, v+diff)
				} else {
					ncurs = append(ncurs, v)
				}
			}
			if len(ncurs) > 0 {
				cur_vals = ncurs
			} else {
				cur_vals = ncurs
				for len(<-lines) > 0 {
				}
				break
			}
		}
		// Move remaining values with their existing value
		if len(cur_vals) > 0 {
			new_vals = append(new_vals, cur_vals...)
		}
		cur_vals = new_vals
	}
	// Sort and pick smallest
	slices.Sort(cur_vals)
	return cur_vals[0]
}

func part2(lines <-chan string) int {
	seeds_strings := strings.Split(strings.Split((<-lines), ": ")[1], " ")
	seed_range_vals := u.ParseNums(seeds_strings)

	var intervals_head *IntervalNode = nil
	{
		var intervals_cur *IntervalNode
		for i := 0; i < len(seed_range_vals); i += 2 {
			new_node := &IntervalNode{
				Interval: IntervalFromStartAndWidth(seed_range_vals[i], seed_range_vals[i+1]),
				Prev:     intervals_cur,
			}
			if intervals_head == nil {
				intervals_head = new_node
			} else {
				intervals_cur.Next = new_node
			}
			intervals_cur = new_node
		}
	}

	// Blank
	<-lines

	// Map Loop
	for {
		_, ok := <-lines
		if !ok {
			break
		}

		var new_intervals_head *IntervalNode
		var new_intervals_cur *IntervalNode
		for mapline := range lines {
			// Blank line, end of map
			if len(mapline) == 0 {
				break
			}
			if intervals_head == nil {
				for len(<-lines) > 0 {
				}
				break
			}

			// Parse
			ma := u.ParseNums(strings.Split(mapline, " "))
			it := &IntervalTransform{
				Interval: IntervalFromStartAndWidth(ma[1], ma[2]),
				diff:     ma[0] - ma[1],
			}

			cur := intervals_head
			for cur != nil {
				intersection_interval, next, self_remove := cur.Transform(it)
				if self_remove && cur == intervals_head {
					intervals_head = next
				}
				cur = next

				if intersection_interval != nil {
					new_interval := &IntervalNode{
						Interval: *intersection_interval,
						Prev:     new_intervals_cur,
					}
					if new_intervals_head == nil {
						new_intervals_head = new_interval
					} else {
						new_intervals_cur.Next = new_interval
					}
					new_intervals_cur = new_interval
				}
			}
		}

		if intervals_head != nil && new_intervals_head != nil {
			cur := intervals_head
			for cur != nil {
				next := cur.Next
				new_intervals_cur.Next = cur
				cur.Prev = new_intervals_cur
				new_intervals_cur = cur
				cur = next
			}
		}
		if new_intervals_head != nil {
			intervals_head = new_intervals_head
		}
	}

	cur := intervals_head
	min := cur.Start
	cur = cur.Next
	for cur != nil {
		if cur.Start < min {
			min = cur.Start
		}
		cur = cur.Next
	}
	return min
}

type Interval struct {
	Start int
	End   int
}

func IntervalFromStartAndWidth(start int, width int) Interval {
	return Interval{start, start + width - 1}
}

type IntervalNode struct {
	Interval
	Next *IntervalNode
	Prev *IntervalNode
}

type IntervalTransform struct {
	Interval
	diff int
}

func (n *IntervalNode) Transform(it *IntervalTransform) (new_interval *Interval, next *IntervalNode, self_remove bool) {
	next = n.Next
	if max(n.Start, it.Start) <= min(n.End, it.End) {
		new_interval = &Interval{}
		var splinter *IntervalNode = nil
		if n.End > it.End {
			// Splinter
			splinter = &IntervalNode{Interval: Interval{Start: it.End + 1, End: n.End}}
			new_interval.End = it.End + it.diff
		} else {
			new_interval.End = n.End + it.diff
		}

		if n.Start < it.Start {
			new_interval.Start = it.Start + it.diff
			n.End = it.Start - 1
		} else {
			new_interval.Start = n.Start + it.diff
			if splinter != nil {
				// Replace
				n.Start = splinter.Start
				n.End = splinter.End
				splinter = nil
			} else {
				// Remove self
				if n.Prev != nil {
					n.Prev.Next = next
				}
				if next != nil {
					next.Prev = n.Prev
				}

				self_remove = true
				return
			}
		}

		if splinter != nil {
			// Insert between self and next
			n.Next = splinter
			if next != nil {
				splinter.Next = next
				next.Prev = splinter
			}
		}
	}
	return
}
