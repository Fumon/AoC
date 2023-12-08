package main

import (
	"container/heap"
	"fmt"
	"fuaoc2023/day07/u"
	"strings"
)

func main()  {
	fmt.Println(part1(u.Linewisefile_chan("input")))
}

type Hand struct {
	Hand        string
	Cards       [5]int16
	UniqueFaces int
	Q           int
	Bid         int
}

func cardToInt(r byte) (out int16, rank int) {
	switch r {
	case '2':
		return 1, 0
	case '3':
		return 2, 1
	case '4':
		return 4, 2
	case '5':
		return 8, 3
	case '6':
		return 16, 4
	case '7':
		return 32, 5
	case '8':
		return 64, 6
	case '9':
		return 128, 7
	case 'T':
		return 256, 8
	case 'J':
		return 512, 9
	case 'Q':
		return 1024, 10
	case 'K':
		return 2048, 11
	case 'A':
		return 4096, 12
	default:
		panic(fmt.Sprintf("Unknown card %v", string(r)))
	}
}

func part1(lines <-chan string) int {
	hands := &HandHeap{}
	heap.Init(hands)

	for line := range lines {
		g := strings.Split(line, " ")
		var cards [5]int16
		var faces_proxy int16
		var ranks [13]int
		for i, c := range []byte(g[0]) {
			ci, rank := cardToInt(c)
			cards[i] = ci
			faces_proxy |= ci
			ranks[rank] += 1
		}
		var qsig = 1
		for _, r := range ranks {
			if r > 0 {
				qsig *= r
			}
		}

		var unique_faces int
		for ; faces_proxy != 0; faces_proxy >>= 1 {
			unique_faces += int(faces_proxy & 1)
		}

		bid := u.ParseNum(g[1])

		heap.Push(hands, &Hand{
			Hand:        g[0],
			Cards:       cards,
			UniqueFaces: unique_faces,
			Q:           qsig,
			Bid:         bid,
		})
	}

	var total int
	for rank := hands.Len(); rank > 0; rank-- {
		hand := heap.Pop(hands).(*Hand)
		fmt.Println(total, " += ", hand.Bid * rank, ": ", *hand)
		total += hand.Bid * rank
	}
	return total
}
