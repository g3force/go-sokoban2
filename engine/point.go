package engine

import ()

// simple Point type
type Point struct {
	X int
	Y int
}

type Points map[Point]bool

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

func (points *Points) Contains(p Point) bool {
	return (*points)[p]
}

// Methods required by sort.Interface.
//func (s Points) Len() int {
//	return len(s)
//}
//func (s Points) Less(i, j int) bool {
//	if s[i].X != s[j].X {
//		return s[i].X < s[j].X
//	}
//	return s[i].Y < s[j].Y
//}
//func (s Points) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}

// Method for printing - sorts the elements before printing.
//func (s Points) String() string {
//	str := "["
//	for i, elem := range s {
//		if i > 0 {
//			str += " "
//		}
//		str += fmt.Sprint(elem)
//	}
//	return str + "]"
//}
