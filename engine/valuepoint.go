package engine

type ValuePoint struct {
	P     Point
	Value float32
}

type ValuePoints []ValuePoint

func (sm *ValuePoints) Len() int {
	return len(*sm)
}

func (sm *ValuePoints) Less(i, j int) bool {
	return (*sm)[i].Value > (*sm)[j].Value
}

func (sm *ValuePoints) Swap(i, j int) {
	(*sm)[i], (*sm)[j] = (*sm)[j], (*sm)[i]
}
