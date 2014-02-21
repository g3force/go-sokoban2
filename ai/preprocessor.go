package ai

import (
	"github.com/g3force/go-sokoban2/engine"
	//"github.com/op/go-logging"
)

// check, if given point is a dead corner
func DeadCorner(e *engine.Engine, point engine.Point) (found bool, x int8) {
	var p engine.Point
	hit := false
	x = 0
	found = false
	if (*e).Targets.Contains(point) {
		return
	}

	// check clockwise, if there is a wall or not.
	// If there is a wall two times together, corner is dead
	for i := 0; i < 5; i++ {
		x = int8(i % 4)
		p = point.Add(engine.Direction(x).Point())
		if !(*e).Surface.In(p) || (*e).Surface[p.Y][p.X] == engine.WALL {
			if hit {
				found = true
				return
			} else {
				hit = true
			}
		} else {
			hit = false
		}
	}
	return
}

func MarkDeadFields(e *engine.Engine) {
	for y := 0; y < len((*e).Surface); y++ {
		for x := 0; x < len((*e).Surface[y]); x++ {
			thisPoint := engine.NewPoint(x, y)
			// walls can't be dead fields
			if !(*e).Surface.In(thisPoint) || (*e).Surface[y][x] == engine.WALL {
				continue
			}
			dead, dir1 := DeadCorner(e, thisPoint)
			if dead {
				(*e).Surface[y][x] = engine.DEAD
				dir1 = (dir1 + 2) % 4 //dir1, dir2 are the directions of a possible dead wall
				dir2 := (dir1 - 1) % 4
				deadWall, p := checkForDeadWall(e, thisPoint, dir1, (dir2+2)%4)
				if deadWall {
					markDeadWall(e, thisPoint, p)
				}
				deadWall, p = checkForDeadWall(e, thisPoint, dir2, (dir1+2)%4)
				if deadWall {
					markDeadWall(e, thisPoint, p)
				}
			}
		}
	}
}

//deadEdge: first dead Edge to star
//dir: direction where the wall will go on
//wallDir: direction of the wall, left or right of the dir???
func checkForDeadWall(e *engine.Engine, deadEdge engine.Point, dir int8, wallDir int8) (bool, engine.Point) {
	possDead := deadEdge
	for {
		possDead = possDead.Add(engine.Direction(dir).Point())
		if !(*e).Surface.In(possDead) {
			return false, possDead
		}
		possField := (*e).Surface[possDead.Y][possDead.X]
		possWallPos := possDead.Add(engine.Direction(wallDir).Point())
		if !(*e).Surface.In(possWallPos) {
			return false, possDead
		}
		possWall := (*e).Surface[possWallPos.Y][possWallPos.X]
		if possField == engine.WALL || (*e).Targets.Contains(possDead) || possWall != engine.WALL {
			return false, possDead
		} else {
			dead, _ := DeadCorner(e, possDead)
			if dead {
				return true, possDead
			}
		}
	}
	log.Error("checkForDeadWall: end of For loop")
	return false, possDead
}

func markDeadWall(e *engine.Engine, start engine.Point, end engine.Point) {
	if start.X == end.X && start.Y != end.Y {
		if start.Y < end.Y {
			for i := start.Y; i <= end.Y; i++ {
				((*e).Surface)[i][start.X] = engine.DEAD
			}
		} else {
			for i := start.Y; i >= end.Y; i-- {
				((*e).Surface)[i][start.X] = engine.DEAD
			}
		}
	} else if start.Y == end.Y && start.X != end.X {
		if start.X < end.X {
			for i := start.X; i <= end.X; i++ {
				((*e).Surface)[start.Y][i] = engine.DEAD
			}
		} else {
			for i := start.X; i >= end.X; i-- {
				((*e).Surface)[start.Y][i] = engine.DEAD
			}
		}
	} else {
		log.Debug("Solo dead end")
	}
	return
}
