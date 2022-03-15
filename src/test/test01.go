package main

import (
	"fmt"
)

type hasX interface {
	getX() int
}

type point struct {
	x int
	y int
}

func (p *point) getX() int {
	return p.x
}

type rect struct {
	x      int
	y      int
	width  int
	height int
}

func (r *rect) getX() int {
	return r.x
}

func main() {
	fmt.Println("开始。")
	var v1 hasX = &point{3, 4}
	var v2 hasX = &rect{5, 6, 7, 8}
	y1 := v1.(*point).y
	fmt.Printf("y1=%d\n", y1)
	r1 := v2.(*point)
	y2 := r1.y
	fmt.Printf("y2=%d\n", y2)
}
