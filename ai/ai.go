package ai

import (
	"bytes"
	"fmt"
	"github.com/g3force/go-sokoban2/engine"
	"github.com/op/go-logging"
	"math"
	"sort"
)

var log = logging.MustGetLogger("sokoban")

type AI_Node struct {
	S   engine.State
	V   float64
	pre *AI_Node
}

type AI_Nodes []AI_Node

func (sm *AI_Nodes) Len() int {
	return len(*sm)
}

func (sm *AI_Nodes) Less(i, j int) bool {
	return (*sm)[i].V < (*sm)[j].V
}

func (sm *AI_Nodes) Swap(i, j int) {
	(*sm)[i], (*sm)[j] = (*sm)[j], (*sm)[i]
}

func (ns *AI_Nodes) String() string {
	var buffer bytes.Buffer
	for _, n := range *ns {
		buffer.WriteString(fmt.Sprintf("{%s, %f, %d}\n", n.S.Boxes.String(), n.V, &n.pre))
	}
	return buffer.String()
}

type States []engine.State

func (s *States) Contains(state *engine.State, e *engine.Engine) bool {
	for _, passedBoxSet := range *s {
		if passedBoxSet.Boxes.ContainsAll((*state).Boxes) {
			_, ok := e.FigureShortestPath(passedBoxSet.Figure, (*state).Figure)
			if ok {
				return true
			}
		}
	}
	return false
}

func (s *States) Insert(state *engine.State, e *engine.Engine) bool {
	if !(*s).Contains(state, e) {
		(*s) = append(*s, *state)
		return true
	}
	return false
}

func (states *States) String() string {
	var buffer bytes.Buffer
	for _, s := range *states {
		buffer.WriteString(fmt.Sprintf("{%s, %s}\n", s.Boxes.String(), s.Figure.String()))
	}
	return buffer.String()
}

func Solve(e *engine.Engine) (path []engine.Point, solved bool) {
	dist := calcDistBoxes((*e).CurrentState.Boxes, (*e).Targets)
	initState := engine.State{(*e).CurrentState.Boxes, (*e).CurrentState.Figure}
	fringe := AI_Nodes{AI_Node{initState, dist, nil}}
	passed := States{initState}

	for len(fringe) > 0 {
		log.Debug("Fringe: " + fringe.String())
		log.Debug("Passed: " + passed.String())
		curNode := fringe[0]
		e.PrintState(curNode.S)
		//var in string
		//fmt.Scanf("%s", &in)
		if e.Won() {
			solved = true
			//node := &curNode
			//for (*node).pre != nil {
			//	path = append(path, (*node).P)
			//	node = (*node).pre
			//}
			//path = append(path, (*node).P)
			return
		}

		// update engine
		e.CurrentState = curNode.S
		// and remove it
		fringe = fringe[1:len(fringe)]
		// we passed it
		passed.Insert(&curNode.S, e)
		// expand best node
		nodes := expand(&curNode, e, &passed)

		// update fringe todo inefficient
		for _, node := range *nodes {
			fringe = append(fringe, node)
		}
		sort.Sort(&fringe)
	}
	log.Error("No solution found.")
	return
}

func expand(node *AI_Node, e *engine.Engine, passed *States) (nodes *AI_Nodes) {
	nodes = new(AI_Nodes)
	for box, _ := range (*node).S.Boxes {
		for i := 0; i < 4; i++ {
			pfigure := engine.PointAfterMove(box, (i % 4))
			pbox := engine.PointAfterMove(box, ((i + 2) % 4))

			if PointFreeForFigure(pfigure, e) && PointFreeForBox(pbox, e) {
				figPath, ok := e.FigureShortestPath((*node).S.Figure, pfigure)
				if ok {
					nState := engine.State{(*node).S.Boxes.Clone(), box}
					delete(nState.Boxes, box)
					nState.Boxes[pbox] = true
					if !passed.Contains(&nState, e) {
						dist := calcDistBoxes(nState.Boxes, (*e).Targets)
						dist += float64(len(figPath))
						*nodes = append(*nodes, AI_Node{nState, dist, node})
					}
				}
			}
		}
	}
	return
}

func calcDistBoxes(b1 engine.Points, b2 engine.Points) (dist float64) {
	dist = 0
	for k1, _ := range b1 {
		tmpDist := math.Inf(1)
		for k2, _ := range b2 {
			tmp := calcDistPP(k1, k2)
			if tmp < tmpDist {
				tmpDist = tmp
			}
		}
		dist += tmpDist
		//fmt.Printf("%f\n", tmpDist)
	}
	return
}

func calcDistPP(p1 engine.Point, p2 engine.Point) (dist float64) {
	//dist = math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2))
	dist = math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
	return
}

func PointFreeForFigure(dest engine.Point, e *engine.Engine) (valid bool) {
	valid = false

	if e.Surface[dest.Y][dest.X] == engine.WALL {
		return
	}

	if !e.Surface.In(dest) {
		return
	}

	valid = !e.CurrentState.Boxes.Contains(dest)
	return
}

func PointFreeForBox(dest engine.Point, e *engine.Engine) (valid bool) {
	valid = false

	if e.Surface[dest.Y][dest.X] == engine.WALL {
		return
	}

	if e.Surface[dest.Y][dest.X] == engine.DEAD {
		return
	}

	if !e.Surface.In(dest) {
		return
	}

	valid = !e.CurrentState.Boxes.Contains(dest)
	return
}
