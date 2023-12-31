package main

import (
	"fmt"
	"fuaoc2023/day19/u"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {

	delimit := u.LineDelimiter(lines, "")
	var rules []*Rule
	var accepted_sum int

	accfn := func(p *Part, ind uint8, val int) int {
		for _, v := range p {
			accepted_sum += v
		}
		return 0
	}
	accepted := &Rule{fn: accfn}
	rules = append(rules, accepted)

	rejected := &Rule{fn: rejected}
	rules = append(rules, rejected)

	var workflows = make(map[string]Workflow)
	workflows["A"] = accepted
	workflows["R"] = rejected

	{
		var workflow_dangling_backlinks_t = make(map[string][]*Rule)
		var workflow_dangling_backlinks_f = make(map[string][]*Rule)

		// Rules
		for line := range <-delimit {
			comp := strings.Split(line, "{")
			name := comp[0]

			var current *Rule = &Rule{}
			var previous *Rule = nil
			workflows[name] = current

			for _, rule_string := range strings.Split(comp[1][:len(comp[1])-1], ",") {
				ind := strings.IndexByte(rule_string, ':')
				if ind > -1 {
					current.accIndex = uint8(byte_to_part_index[rule_string[0]])
					current.val, _ = strconv.Atoi(rule_string[2:ind])
					switch rule_string[1] {
					case '<':
						current.fn = lt
					case '>':
						current.fn = gt
					}
					workflow_dangling_backlinks_t[rule_string[ind+1:]] = append(workflow_dangling_backlinks_t[rule_string[ind+1:]], current)
					if previous != nil {
						previous.f = current
					}
					previous = current
					current = &Rule{}
				} else {
					// End
					workflow_dangling_backlinks_f[rule_string] = append(workflow_dangling_backlinks_f[rule_string], previous)
				}
			}
		}

		for workflowname, rulelist := range workflow_dangling_backlinks_t {
			wf := workflows[workflowname]
			for _, v := range rulelist {
				v.t = wf
			}
		}
		for workflowname, rulelist := range workflow_dangling_backlinks_f {
			wf := workflows[workflowname]
			for _, v := range rulelist {
				v.f = wf
			}
		}
	}

	// if workflows["in"].f.f.f != rejected {
	// 	for k, v := range workflows {
	// 		if v == workflows["in"].f.f.f {
	// 			fmt.Println("actually linked to ", k)
	// 			break
	// 		}
	// 	}
	// 	panic("Fdafsda")
	// }

	in := workflows["in"]

	// Parts
	for line := range <-delimit {
		part := PartFromLine(line)

		currule := in
		evalloop:
		for {
			switch currule.fn(&part, currule.accIndex, currule.val) {
			case 1:
				currule = currule.t
			case 2:
				currule = currule.f
			case 0:
				break evalloop
			}
		}
	}

	return accepted_sum
}

func Part2(lines <-chan string) int {
	return 0
}

type Rule struct {
	accIndex uint8
	val      int
	fn       func(*Part, uint8, int) int

	f *Rule
	t *Rule
}
type Workflow *Rule

func lt(p *Part, index uint8, val int) int {
	if p[index] < val {
		return 1
	} else {
		return 2
	}
}

func gt(p *Part, index uint8, val int) int {
	if p[index] > val {
		return 1
	} else {
		return 2
	}
}

func rejected(p *Part, index uint8, val int) int {
	return 0
}

var byte_to_part_index = map[byte]int{
	'x': 0,
	'm': 1,
	'a': 2,
	's': 3,
}

type Part [4]int

func PartFromLine(line string) (output Part) {
	for i, v := range strings.Split(strings.Trim(line, "{}"), ",") {
		output[i], _ = strconv.Atoi(v[2:])
	}
	return
}
