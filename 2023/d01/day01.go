package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	linescan := bufio.NewScanner(file)
	var result int
	for linescan.Scan() {
		line := linescan.Text()

		var f, e rune
		var fInd int
		for ind, ch := range line {
			if isdigit(ch) {
				f = ch
				fInd = ind
				break
			}
		}
		rline := []rune(line)
		for i := len(rline) - 1; i >= fInd; i-- {
			if isdigit(rline[i]) {
				e = rline[i]
				break
			}
		}
		out, err := strconv.Atoi(string([]rune{f, e}))
		if err != nil {
			log.Fatal(err)
		}

		result += out
	}
	fmt.Println(result)
}

func isdigit(r rune) bool {
	return r >= '0' && r <= '9'
}
