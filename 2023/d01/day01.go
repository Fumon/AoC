package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

var istoint = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
	"0":     0,
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
}

func main() {
	rp := regexp.MustCompile(`([[:digit:]]|one|two|three|four|five|six|seven|eight|nine)`)

	file, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	linescan := bufio.NewScanner(file)
	var result int
	for linescan.Scan() {
		line := linescan.Text()
		mat := rp.FindAllString(line, -1)
		f, l := mat[0], mat[len(mat) - 1]

		fmt.Println(line, " -> ", mat, "\n\t", f, l, "->", 10*istoint[mat[0]] + istoint[mat[len(mat) - 1]])

		result += 10*istoint[mat[0]] + istoint[mat[len(mat) - 1]]
	}
	fmt.Println(result)
}
