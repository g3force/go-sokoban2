package engine

type Direction int8 // 0-3, -1 for invalid
const (
	NO_DIRECTION = Direction(-1)
	RIGHT        = 0
	DOWN         = 1
	LEFT         = 2
	UP           = 3
)

// convert direction from int to Point
func (dir Direction) Point() (p Point) {
	dir = dir % 4
	switch dir {
	case RIGHT: // right
		p.X = 1
		p.Y = 0
	case DOWN: // down
		p.X = 0
		p.Y = 1
	case LEFT: // left
		p.X = -1
		p.Y = 0
	case UP: // up
		p.X = 0
		p.Y = -1
	}
	return p
}
