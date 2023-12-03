package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}

	ls := bufio.NewScanner(f)

	var powersum int
	for ls.Scan() {
		line := ls.Text()
		gameandcounts := strings.Split(line, ":")
		var maxcount = map[string]int {
			"red": 0,
			"green": 0,
			"blue": 0,
		}
		for _, group := range strings.Split(gameandcounts[1], "; ") {
			for _, count := range strings.Split(strings.Trim(group, " "), ", ") {
				foo := strings.Split(count, " ")
				counti, err := strconv.Atoi(foo[0])
				if err != nil {
					log.Panic(err)
				}

				val, ok := maxcount[foo[1]]
				if !ok {
					maxcount[foo[1]] = counti
				} else {
					maxcount[foo[1]] = max(val, counti)
				}
			}
		}
		powersum += maxcount["red"] * maxcount["green"] * maxcount["blue"]
	}
	fmt.Println(powersum)
}