package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var maxcolmap = map[string]int {
	"red": 12,
	"green": 13,
	"blue": 14,
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}

	ls := bufio.NewScanner(f)

	var codeSumCount int
	ScanLoop:
	for ls.Scan() {
		line := ls.Text()
		gameandcounts := strings.Split(line, ":")
		gameid, err :=  strconv.Atoi(strings.Split(gameandcounts[0], "Game ")[1])
		if err != nil {
			log.Panic(err)
		}
		for _, group := range strings.Split(gameandcounts[1], "; ") {
			for _, count := range strings.Split(strings.Trim(group, " "), ", ") {
				foo := strings.Split(count, " ")
				counti, err := strconv.Atoi(foo[0])
				if err != nil {
					log.Panic(err)
				}
				if counti > maxcolmap[foo[1]] {
					continue ScanLoop
				}
			}
		}
		codeSumCount += gameid
	}
	fmt.Println(codeSumCount)
}