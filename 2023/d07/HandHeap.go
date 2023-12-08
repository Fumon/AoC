package main

// An HandHeap is a min-heap of ints.
type HandHeap []*Hand

func (h HandHeap) Len() int { return len(h) }
func (h HandHeap) Less(i, j int) bool {
	d := h[i].UniqueFaces - h[j].UniqueFaces
	if d != 0 {
		return d < 0
	}

	d = h[i].Q - h[j].Q
	if d != 0 {
		return d < 0
	}



	var icard int16
	var jcard int16
	for k := 0; k < 5; k++ {
		icard = h[i].Cards[k]
		jcard = h[j].Cards[k]

		if icard != jcard {
			return icard > jcard
		}
	}
	panic("Duplicate Hands")
}

func (h HandHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *HandHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Hand))
}

func (h *HandHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
