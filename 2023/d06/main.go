package main

import (
	"fmt"
	"fuaoc2023/day06/u"
	"math"
	"strings"
)

func main() {
    fmt.Println(part1(form_input(u.Linewisefile_chan("input"))))
	fmt.Println(part2(form_input2(u.Linewisefile_chan("input"))))
}

func form_input(lines <-chan string) (output [][]int) {
	times := u.ParseNums(strings.Split(strings.Split((<-lines), ":")[1], " "))
	distances := u.ParseNums(strings.Split(strings.Split((<-lines), ":")[1], " "))
	for i, t := range times {
		output = append(output, []int{t, distances[i]})
	}
	return
}

func form_input2(lines <-chan string) (output [][]int) {
	times := u.ParseNums([]string{strings.ReplaceAll(strings.Split((<-lines), ":")[1], " ", "")})
	distances := u.ParseNums([]string{strings.ReplaceAll(strings.Split((<-lines), ":")[1], " ", "")})
	for i, t := range times {
		output = append(output, []int{t, distances[i]})
	}
	return
}

func part1(races [][]int) (ways int) {
	ways = 1
	for _, race := range races {
		min, max := SolveQ(race[0], race[1])
		
		mway := (max - min) + 1
		ways *= mway
	}
	return
}

func part2(races [][]int) (ways int) {
	ways = 1
	for _, race := range races {
		min, max := SolveQ(race[0], race[1])
		
		mway := (max - min) + 1
		ways *= mway
	}
	return
}

func SolveQ(tlimit int, distancerecord int) (min, max int) {
	discrimsqrt := math.Sqrt(float64((tlimit * tlimit) - (4*distancerecord)))
	return int(math.Ceil(0.5 * (float64(tlimit) - discrimsqrt) + 0.001)),  int(math.Floor(0.5 * (float64(tlimit) + discrimsqrt) - 0.001))
}