package u

import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

func Linewisefile_chan(filename string) <-chan string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	outchan := make(chan string, 3)

	go func() {
		for scanner.Scan() {
			outchan <- scanner.Text()
		}
		close(outchan)
		f.Close()
	}()

	return outchan
}

func Linewisefile_slice(filename string) (out []string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	f.Close()
	return
}

type ASCIINumberHashSet struct {
	set [256]bool
}

func (a *ASCIINumberHashSet) Insert(s []byte) {
	a.set[(s[0]<<4)|(s[1]&0xF)] = true
}

func (a *ASCIINumberHashSet) Exists(s []byte) bool {
	return a.set[(s[0]<<4)|(s[1]&0xF)]
}

func ParseNums(sarr []string) (out []int) {
	for _, s := range sarr {
		if len(s) > 0 {
			v, err := strconv.Atoi(s)
			if err != nil {panic(err)}
			out = append(out, v)
		}
	}
	return
}

func ParseNum(s string) int {
	if len(s) == 0 {
		panic("NO EMPTY STRINGS")
	}
	v, err := strconv.Atoi(s)
	if err != nil {panic(err)}
	return v
}

func Assert[T comparable](t *testing.T, e, r T) {
	if e != r {
		t.Error("Expected ", e, " but got ", r)
	}
}

func LineDelimiter(lines <-chan string, linedelimiter string) <-chan chan string {
	var output = make(chan chan string, 1)
	go func () {
		var currentchan = make(chan string)
		output <- currentchan
		for line := range lines {
			if line == linedelimiter {
				close(currentchan)
				currentchan = make(chan string)
				output <- currentchan
			} else {
				currentchan <- line
			}
		}
		close(output)
		close(currentchan)
	}()
	return output
}

func PeekAtFirstLineWidth(lines <-chan string) (<-chan string, int) {
	line1 := <-lines
	nc := make(chan string, 3)
	nc <- line1
	go func() {
		for line := range lines {
			nc <- line
		}
		close(nc)
	}()

	return nc, len(line1)
}
