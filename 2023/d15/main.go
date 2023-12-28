package main

import (
	"fmt"
	"fuaoc2023/day15/u"
	"strings"
)

const DEBUG = true

func main() {
	fmt.Println("Part1:", Part1(u.Linewisefile_chan("input")))
	fmt.Println("Part2:", Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) (output int) {
	for line := range lines {
		for _, str := range strings.Split(line, ",") {
			output += int(HASH(str))
		}
	}
	return
}

func Part2(lines <-chan string) (output int) {
	var lens_boxes ArrLensBox

	for line := range lines {
		for _, step := range strings.Split(line, ",") {
			command := ParseStep(step)
			lens_boxes.Exec(command)
			fmt.Print("After \"", step, "\":\n", lens_boxes.String(), "\n")
		}
	}

	for i, b := range lens_boxes {
		slot_num := 1
		cur := b
		for cur != nil {
			output += (i + 1) * slot_num * int(cur.focal_length)
			slot_num++
			cur = cur.next
		}
	}

	return
}

func HASH(str string) (result uint8) {
	for _, ch := range []byte(str) {
		result += ch
		result *= 17
	}
	return
}

type Command struct {
	box   uint8
	label string
	mode  bool // false for remove, true for place
	focal uint8
}

func ParseStep(step string) (c Command) {
	splitpoint := strings.IndexAny(step, "-=")
	c.label = step[:splitpoint]
	c.box = HASH(c.label)
	if step[splitpoint] == '=' {
		c.mode = true
		c.focal = []byte(step)[splitpoint+1] & 0x0F
	}
	return
}

type ArrLensBox [256]*LensInBox

func (alb *ArrLensBox) Exec(cmd Command) {
	var prev *LensInBox
	var cur *LensInBox = alb[cmd.box]
	for cur != nil {
		if cur.label == cmd.label {
			break
		}
		prev = cur
		cur = cur.next
	}

	if cmd.mode {
		if cur == nil {
			newbox := LensInBox{
				label:        cmd.label,
				focal_length: cmd.focal,
			}
			if prev == nil {
				alb[cmd.box] = &newbox
			} else {
				prev.next = &newbox
			}
		} else {
			cur.focal_length = cmd.focal
		}
	} else if cur != nil {
		if prev == nil {
			alb[cmd.box] = cur.next
		} else {
			prev.next = cur.next
		}
	}
}

func (alb *ArrLensBox) String() string {
	var bild strings.Builder
	for i, b := range alb {
		if b == nil { continue }
		fmt.Fprint(&bild, "Box ", i, ":")
		cur := b
		for cur != nil {
			fmt.Fprint(&bild, " ", cur)
			cur = cur.next
		}
		if DEBUG {
			fmt.Fprint(&bild, "\n")
		}
	}

	return bild.String()
}

type LensInBox struct {
	label        string
	focal_length uint8
	next         *LensInBox
}

func (lib *LensInBox) String() string {
	return fmt.Sprint("[", lib.label, " ", lib.focal_length, "]")
}
