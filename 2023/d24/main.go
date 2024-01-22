package main

import (
	"fmt"
	"fuaoc2023/day24/u"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input"), [2]int{200000000000000, 400000000000000}))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string, XY_min_max [2]int) int {

	var hailstones []Hailstone
	for line := range lines {
		hailstones = append(hailstones, parseHailstone(line))
	}

	var intersection_count int
	for i := 0; i < len(hailstones) - 1; i++ {
		A := hailstones[i]
		for j := i + 1; j < len(hailstones); j++ {
			B := hailstones[j]

			// fmt.Println("Hailstone A:", A)
			// fmt.Println("Hailstone B:", B)

			// Check if A and B's X and Y path intersect in the future using linear algebra
			// If they do collide, ensure the collision is within the XYmax bounds
			if A.SpeedVector[0] == B.SpeedVector[0] && A.SpeedVector[1] == B.SpeedVector[1] {
				// If A and B are already colliding, increment collision_count
				if A.Start[0] == B.Start[0] && A.Start[1] == B.Start[1] && A.Start[2] == B.Start[2] {
					intersection_count++
					continue
				} else {
					continue
				}
			}
			
			ls := u.Linear_System{
				A: u.Two_by_two_matrix{
					{A.SpeedVector[0], -B.SpeedVector[0]},
					{A.SpeedVector[1], -B.SpeedVector[1]},
				},
				B: u.Two_by_one_matrix{
					{B.Start[0] - A.Start[0]},
					{B.Start[1] - A.Start[1]},
				},
			}

			// If the system is not solvable, continue
			if ls.A.Det() == 0 {
				// fmt.Println("Hailstones' paths are parallel; they never intersect.")
				continue
			}

			// If the system is solvable, find the solution
			// If the solution is not within the XYmax bounds, continue
			// If the solution is within the XYmax bounds, increment collision_count
			solution := ls.Solve()
			var a_in_past bool = solution[0][0] < 0
			var b_in_past bool = solution[1][0] < 0
			if a_in_past && b_in_past {
				// fmt.Println("Hailstones' paths crosed in the past for both hailstones.")
				continue
			} else if a_in_past {
				// fmt.Println("Hailstones' paths crosed in the past for hailstone A.")
				continue
			} else if b_in_past {
				// fmt.Println("Hailstones' paths crosed in the past for hailstone B.")
				continue
			}

			// Get X and Y coordinates of the solution
			var sol_x = solution[0][0]*A.SpeedVector[0] + A.Start[0]
			var sol_y = solution[0][0]*A.SpeedVector[1] + A.Start[1]
			
			var inside_XY_min_max bool = sol_x >= float64(XY_min_max[0]) && sol_x <= float64(XY_min_max[1]) && sol_y >= float64(XY_min_max[0]) && sol_y <= float64(XY_min_max[1])

			if inside_XY_min_max {
				// fmt.Printf("Hailstones' paths will cross inside the test area (at x=%v, y=%v).\n", sol_x, sol_y)
				intersection_count++
			} // else {
			// 	fmt.Printf("Hailstones' paths will cross outside the test area (at x=%v, y=%v).\n", sol_x, sol_y)
			// }
		}
	}

	return intersection_count
}

func Part2(lines <-chan string) int {
	return 0
}