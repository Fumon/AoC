package u

import (
	"bufio"
	"os"
)

func Linewisefile(filename string) <-chan string {
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