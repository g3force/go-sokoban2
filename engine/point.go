package engine

import (
	"strconv"
)

// simple Point type
type Point struct {
	X int
	Y int
}

type Points map[Point]bool
type PointList []Point

// add to points (their x and y)
func (p1 *Point) Add(p2 Point) Point {
	var newP Point
	newP.X = p1.X + p2.X
	newP.Y = p1.Y + p2.Y
	return newP
}

func NewPoint(x int, y int) Point {
	return Point{x, y}
}

func (p *Point) Clone() Point {
	return NewPoint((*p).X, (*p).Y)
}

func (ps *Points) Clone() (nps Points) {
	nps = make(map[Point]bool, len(*ps))
	for ps, _ := range *ps {
		nps[ps] = true
	}
	return
}

func (points *Points) Contains(p Point) bool {
	return (*points)[p]
}

func (ps1 *Points) ContainsAll(ps2 Points) bool {
	for p, _ := range ps2 {
		if !(*ps1)[p] {
			return false
		}
	}
	return true
}

func (ps *Points) String() (result string) {
	for p, _ := range *ps {
		result += p.String()
	}
	return
}

func (p *Point) String() (result string) {
	result = "[" + strconv.Itoa(p.X) + " " + strconv.Itoa(p.Y) + "]"
	return
}

// Methods required by sort.Interface.
func (s PointList) Len() int {
	return len(s)
}
func (s PointList) Less(i, j int) bool {
	if s[i].X != s[j].X {
		return s[i].X < s[j].X
	}
	return s[i].Y < s[j].Y
}
func (s PointList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
