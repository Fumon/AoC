package ek

import "fmt"

type EdgeKey string

var edge_key_memo = map[[2]int]EdgeKey{}

func New(a, b int) EdgeKey {
	input_key := [2]int{a, b}
	if v, ok := edge_key_memo[input_key]; ok {
		return v
	}
	if b < a {
		a, b = b, a
	}
	
	out := EdgeKey(fmt.Sprintf("%d-%d", a, b))
	edge_key_memo[input_key] = out
	return out
}

func (ek EdgeKey) Recover_indicies() [2]int {
	var first, second int
	fmt.Sscanf(string(ek), "%d-%d", &first, &second)
	return [2]int{first, second}
}