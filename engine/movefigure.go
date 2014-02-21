package engine

import (
	"math"
	"sort"
)

type MF_Node struct {
	P   Point
	V   float32
	pre *MF_Node
}

type MF_Nodes []MF_Node

func (sm *MF_Nodes) Len() int {
	return len(*sm)
}

func (sm *MF_Nodes) Less(i, j int) bool {
	return (*sm)[i].V < (*sm)[j].V
}

func (sm *MF_Nodes) Swap(i, j int) {
	(*sm)[i], (*sm)[j] = (*sm)[j], (*sm)[i]
}

func (e *Engine) FigureShortestPath(src Point, dest Point) (path []Point, exists bool) {
	exists = false
	path = []Point{}
	passed := Points{}
	dist := calcDist(src, dest)
	fringe := MF_Nodes{MF_Node{src, dist, nil}}

	for len(fringe) > 0 {
		curNode := fringe[0]
		if curNode.P == dest {
			exists = true
			node := &curNode
			for (*node).pre != nil {
				path = append(path, (*node).P)
				node = (*node).pre
			}
			path = append(path, (*node).P)
			return
		}
		// expand best node
		nodes := figureShortestPathExpand(curNode.P, &passed, e)
		// and remove it
		fringe = fringe[1:len(fringe)]

		// add expanded nodes to fringe
		for _, node := range nodes {
			dist = calcDist(node, dest)
			fringe = append(fringe, MF_Node{node, dist, &curNode})
		}
		sort.Sort(&fringe)
	}
	//log.Error("Could not move figure to given destination")
	return
}

func calcDist(p1 Point, p2 Point) (dist float32) {
	dist = float32(math.Sqrt(math.Pow(float64(p1.X-p2.X), 2) + math.Pow(float64(p1.Y-p2.Y), 2)))
	return
}

func figureShortestPathExpand(node Point, passed *Points, e *Engine) (nodes []Point) {
	nodes = []Point{}
	for dir := 0; dir < 4; dir++ {
		p := PointAfterMove(node, dir)
		valid, box := e.CheckDestination(p)
		if valid && !box && !(*passed).Contains(p) {
			nodes = append(nodes, p)
			(*passed)[p] = true
		}
	}
	return
}
