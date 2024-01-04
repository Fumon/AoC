package main

import (
	"fmt"
	"fuaoc2023/day19/u"
	"sort"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(Part1(u.Linewisefile_chan("input")))
	fmt.Println(Part2(u.Linewisefile_chan("input")))
}

func Part1(lines <-chan string) int {

	delimit := u.LineDelimiter(lines, "")

	var accepted_sum int
	accfn := func(p *Part, ind uint8, val int) int {
		for _, v := range p {
			accepted_sum += v
		}
		return 0
	}
	accepted := &Rule{fn: accfn}
	rejected := &Rule{fn: rejected}

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

	in := workflows["in"]
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

	delimit := u.LineDelimiter(lines, "")
	accepted := &Rule2{}
	rejected := &Rule2{}

	var workflows = make(map[string]Workflow2)
	workflows["A"] = accepted
	workflows["R"] = rejected

	{
		var workflow_dangling_backlinks_t = make(map[string][]*Rule2)
		var workflow_dangling_backlinks_f = make(map[string][]*Rule2)

		// Rule2s
		for line := range <-delimit {
			comp := strings.Split(line, "{")
			name := comp[0]

			var current *Rule2 = &Rule2{}
			var previous *Rule2 = nil
			workflows[name] = current

			for _, rule_string := range strings.Split(comp[1][:len(comp[1])-1], ",") {
				ind := strings.IndexByte(rule_string, ':')
				if ind > -1 {
					current.accIndex = uint8(byte_to_part_index[rule_string[0]])
					current.val, _ = strconv.Atoi(rule_string[2:ind])
					if rule_string[1] == '>' {
						current.gt = true
					}
					workflow_dangling_backlinks_t[rule_string[ind+1:]] = append(workflow_dangling_backlinks_t[rule_string[ind+1:]], current)
					if previous != nil {
						previous.f = current
					}
					previous = current
					current = &Rule2{}
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

	// Drain
	for range <-delimit {
	}

	in := workflows["in"]
	var accepted_paths, rejected_paths [][]*Rule2
	var accepted_combos, rejected_combos int
	const attr_min, attr_max = 1, 4000

	var b func([]*Rule2, *Rule2, [4]Interval, int)
	b = func(path []*Rule2, next *Rule2, intervals [4]Interval, total_combos int) {
		if next == accepted {
			accepted_paths = append(accepted_paths, path)
			accepted_combos += total_combos
			return
		} else if next == rejected {
			rejected_paths = append(rejected_paths, path)
			rejected_combos += total_combos
			return
		}

		cur_interval := intervals[next.accIndex]

		var true_combos int
		var false_combos int
		var true_intervals [4]Interval = intervals
		var false_intervals [4]Interval = intervals
		if next.gt {
			true_start := next.val + 1
			if true_start > cur_interval.end {
				false_intervals[next.accIndex] = cur_interval
				false_combos = total_combos
			} else if true_start <= cur_interval.start {
				true_intervals[next.accIndex] = cur_interval
				true_combos = total_combos
			} else {
				true_intervals[next.accIndex] = Interval{true_start, cur_interval.end}
				false_intervals[next.accIndex] = Interval{cur_interval.start, true_start - 1}
				true_combos = (total_combos * true_intervals[next.accIndex].size()) / cur_interval.size()
				false_combos = total_combos - true_combos
			}
		} else {
			true_end := next.val - 1
			if true_end < cur_interval.start {
				false_intervals[next.accIndex] = cur_interval
				false_combos = total_combos
			} else if true_end >= cur_interval.end {
				true_intervals[next.accIndex] = cur_interval
				true_combos = total_combos
			} else {
				true_intervals[next.accIndex] = Interval{cur_interval.start, true_end}
				false_intervals[next.accIndex] = Interval{true_end + 1, cur_interval.end}
				true_combos = (total_combos * true_intervals[next.accIndex].size()) / cur_interval.size()
				false_combos = total_combos - true_combos
			}
		}

		if true_combos > 0 {
			tpath := append(path, next)
			b(tpath, next.t, true_intervals, true_combos)
		}

		if false_combos > 0 {
			frule := *next
			if next.gt {
				frule.gt = false
				frule.val++
			} else {
				frule.gt = true
				frule.val--
			}
			fpath := append(path, &frule)
			b(fpath, next.f, false_intervals, false_combos)
		}
	}

	b([]*Rule2{}, in, [4]Interval{{attr_min, attr_max}, {attr_min, attr_max}, {attr_min, attr_max}, {attr_min, attr_max}}, 4000*4000*4000*4000)

	// fmt.Println("Accepted:", len(accepted_paths), "Rejected:", len(rejected_paths))

	// var hyperrectangles []Hyperrectangle
	// for _, path := range accepted_paths {
	// 	var part_intervals [4]*Interval
	// 	for i := range part_intervals {
	// 		np := Interval{attr_min, attr_max}
	// 		part_intervals[i] = &np
	// 	}

	// 	for _, r := range path {
	// 		rinterval := r.to_interval(attr_min, attr_max)
	// 		part_intervals[r.accIndex].intersect(&rinterval)
	// 	}

	// 	// fmt.Print("Ranges ", i, ":")
	// 	// for i, iv := range part_intervals {
	// 	// 	fmt.Printf("\t%d\t>=\t%c\t<=\t%d", iv.start, part_index_to_char[uint8(i)], iv.end)
	// 	// }
	// 	hyperrectangles = append(hyperrectangles, Hyperrectangle{*part_intervals[0], *part_intervals[1], *part_intervals[2], *part_intervals[3]})

	// 	// fmt.Print("\n")
	// }

	// for i, path := range accepted_paths {
	// 	fmt.Printf("Path %3d", i)
	// 	for _, r := range path {
	// 		fmt.Print("\t", r)
	// 	}
	// 	fmt.Print("\n")
	// 	fmt.Printf("  v: %13d\t", hyperrectangles[i].volume())
	// 	fmt.Print(hyperrectangles[i].String(), "\n")
	// }

	// // Detect overlaps
	// var overlaps [4]int
	// var hroverlaps [2]Hyperrectangle
	// for i := 0; i < len(hyperrectangles)-1; i++ {
	// 	ih := hyperrectangles[i]
	// 	for j := i + 1; j < len(hyperrectangles); j++ {
	// 		jh := hyperrectangles[j]
	// 		if ih.X.does_overlap_not_same(&jh.X) {
	// 			overlaps[0]++
	// 			hroverlaps[0] = ih
	// 			hroverlaps[1] = jh
	// 		}
	// 		if ih.Y.does_overlap_not_same(&jh.Y) {
	// 			overlaps[1]++
	// 			hroverlaps[0] = ih
	// 			hroverlaps[1] = jh
	// 		}
	// 		if ih.Z.does_overlap_not_same(&jh.Z) {
	// 			overlaps[2]++
	// 			hroverlaps[0] = ih
	// 			hroverlaps[1] = jh
	// 		}
	// 		if ih.W.does_overlap_not_same(&jh.W) {
	// 			overlaps[3]++
	// 			hroverlaps[0] = ih
	// 			hroverlaps[1] = jh
	// 		}
	// 	}
	// }
	// fmt.Println("Overlaps:", overlaps)
	// fmt.Println("Overlaps:", hroverlaps[0].String(), hroverlaps[1].String())

	// for i, path := range paths {
	// 	var part_attr_inequality_lists [4][]*Rule2
	// 	for _, r := range path {
	// 		part_attr_inequality_lists[r.accIndex] = append(part_attr_inequality_lists[r.accIndex], r)
	// 	}

	// 	fmt.Print("Per attribute ", i, ":\t")
	// 	for i, rlist := range part_attr_inequality_lists {
	// 		fmt.Printf("| %c: ", part_index_to_char[uint8(i)])
	// 		for _, r_in_list := range rlist {
	// 			fmt.Print(r_in_list, " ")
	// 		}
	// 	}
	// 	fmt.Print("\n")
	// }

	// var brute_force_total int
	// for x := attr_min; x <= attr_max; x++ {
	// 	for y := attr_min; y <= attr_max; y++ {
	// 		for z := attr_min; z <= attr_max; z++ {
	// 			for w := attr_min; w <= attr_max; w++ {
	// 				point := [4]int{x, y, z, w}
	// 				for _, hyperrectangle := range hyperrectangles {
	// 					if hyperrectangle.contains(point) {
	// 						brute_force_total += 1
	// 						break
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return accepted_combos
}

type Rule struct {
	accIndex uint8
	val      int
	fn       func(*Part, uint8, int) int

	f *Rule
	t *Rule
}
type Workflow *Rule
type Rule2 struct {
	accIndex uint8
	val      int
	gt       bool

	f *Rule2
	t *Rule2
}
type Workflow2 *Rule2

func (r *Rule2) String() string {
	var inequalsym = '<'
	if r.gt {
		inequalsym = '>'
	}
	return fmt.Sprintf("%c%c%d", part_index_to_char[r.accIndex], inequalsym, r.val)
}

func (r *Rule2) to_interval(min, max int) Interval {
	if r.gt {
		return Interval{r.val + 1, max}
	} else {
		return Interval{min, r.val - 1}
	}
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func intMin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

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
var part_index_to_char = map[uint8]byte{
	0: 'x',
	1: 'm',
	2: 'a',
	3: 's',
}

type Part [4]int

func PartFromLine(line string) (output Part) {
	for i, v := range strings.Split(strings.Trim(line, "{}"), ",") {
		output[i], _ = strconv.Atoi(v[2:])
	}
	return
}

type Interval struct {
	start, end int
}

func (i *Interval) contains(x int) bool {
	return x >= i.start && x <= i.end
}

func (i *Interval) size() int {
	s := i.end - i.start + 1
	if s < 0 {
		return 0
	} else {
		return s
	}
}

func (i *Interval) intersect(other *Interval) {
	st := other.start - i.start
	if st > 0 {
		i.start += st
	}

	ed := other.end - i.end
	if ed < 0 {
		i.end += ed
	}
}

func (i *Interval) does_overlap(other *Interval) bool {
	return (i.start <= other.end) && (other.start <= i.end)
}

func (i *Interval) does_overlap_not_same(other *Interval) bool {
	return (i.start <= other.end) && (other.start <= i.end) && (i.start != other.start || i.end != other.end)
}

// Implement sort.Interface
type Intervals []Interval

func (a Intervals) Len() int           { return len(a) }
func (a Intervals) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Intervals) Less(i, j int) bool { return a[i].start < a[j].start }

type Hyperrectangle struct {
	X, Y, Z, W Interval
}

func (hr *Hyperrectangle) contains(x [4]int) bool {
	return hr.X.contains(x[0]) && hr.Y.contains(x[1]) && hr.Z.contains(x[2]) && hr.W.contains(x[3])
}

func (hr *Hyperrectangle) volume() int {
	return hr.X.size() * hr.Y.size() * hr.Z.size() * hr.W.size()
}

func (hr *Hyperrectangle) String() string {
	return fmt.Sprintf("%4d >= x <= %4d %4d >= m <= %4d %4d >= a <= %4d %4d >= s <= %4d",
		hr.X.start, hr.X.end, hr.Y.start, hr.Y.end, hr.Z.start, hr.Z.end, hr.W.start, hr.W.end)
}

func union(intervals Intervals) Intervals {
	if len(intervals) == 0 {
		return nil
	}
	sort.Sort(intervals)

	var result Intervals
	var current = intervals[0]

	for _, interval := range intervals {
		if interval.start <= current.end {
			if interval.end > current.end {
				current.end = interval.end
			}
		} else {
			result = append(result, current)
			current = interval
		}
	}
	result = append(result, current)

	return result
}
