package u

import (
	"bufio"
	"os"
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
