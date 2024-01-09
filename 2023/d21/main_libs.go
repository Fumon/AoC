package main

type point [2]int

func (p *point) around() (out [4]point) {
	out[0] = point{p[0]+1, p[1]}
	out[1] = point{p[0]-1, p[1]}
	out[2] = point{p[0], p[1]+1}
	out[3] = point{p[0], p[1]-1}
	return
}

type Bounds struct {
	width, height int
}
func (b Bounds) check(p point) bool {
	return (p[0] >= 0 && p[0] < b.width) && (p[1] >= 0 && p[1] < b.height)
}
func (b Bounds) wrap(p point) point {
	return point{
		((p[0] % b.height) + b.height) % b.height,
		((p[1] % b.width) + b.width) % b.width,
	}
}