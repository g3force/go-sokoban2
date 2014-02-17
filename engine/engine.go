package engine

import (
	"bytes"
	"fmt"
	"io"
	"os"
	//"sort"
	"strings"
)

type Field int8

const (
	EMPTY Field = iota
	WALL  Field = iota
	DEAD  Field = iota
)

type State struct {
	Boxes  Points
	Figure Point
}

type Surface [][]Field

type Engine struct {
	Surface      Surface
	Targets      Points
	History      []State
	CurrentState State
}

// load level from specified file (relative to binary file)
func NewEngine(filename string) (e *Engine) {
	raw, err := readLevelAsString(filename)
	if err != nil {
		panic(err)
	}
	// remove the "\r" from stupid windows files...
	raw = strings.Replace(raw, "\r", "", -1)
	// get single lines in an array
	lines := strings.Split(raw, "\n")

	e = new(Engine)
	e.Surface = Surface{{}}
	e.Targets = Points{}
	e.History = []State{}
	e.CurrentState = State{Points{}, Point{0, 0}}
	y := 0
	maxlen := 0
	var char uint8

	for _, line := range lines {
		if len(line) > 0 && line[0] == '#' && len(line) > maxlen {
			maxlen = len(line)
		}
	}

	for _, line := range lines {
		// filter empty lines and lines that do not start with '#'
		if len(line) == 0 || line[0] != '#' {
			continue
		}
		for x := 0; x < maxlen; x++ {
			char = '#'
			if x < len(line) {
				char = line[x]
			}
			switch char {
			case '#':
				// wall
				e.Surface[y] = append(e.Surface[y], WALL)
			case ' ':
				// empty
				e.Surface[y] = append(e.Surface[y], EMPTY)
			case '$':
				// box, empty
				e.Surface[y] = append(e.Surface[y], EMPTY)
				e.CurrentState.Boxes[NewPoint(x, y)] = true
			case '@':
				// figure, empty
				e.Surface[y] = append(e.Surface[y], EMPTY)
				e.CurrentState.Figure = NewPoint(x, y)
			case '.':
				// target
				e.Surface[y] = append(e.Surface[y], EMPTY)
				e.Targets[NewPoint(x, y)] = true
			case '*':
				// target, box
				e.Surface[y] = append(e.Surface[y], EMPTY)
				e.CurrentState.Boxes[NewPoint(x, y)] = true
			case '+':
				// target, figure
				e.Surface[y] = append(e.Surface[y], EMPTY)
				e.CurrentState.Figure = NewPoint(x, y)
			default:
				panic("Unknown character in level file: " + string(char))
			}
		}
		y++
		e.Surface = append(e.Surface, []Field{})
	}
	// the last sub-array of Surface is always empty, so remove it...
	if len(e.Surface[len(e.Surface)-1]) == 0 {
		e.Surface = e.Surface[:len(e.Surface)-1]
	}
	return
}

// read from the specified file and return whole content in a string
func readLevelAsString(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	return string(result), nil
}

// print a legend of the Surface output
func PrintInfo() {
	fmt.Print("Surface Field association:\n")
	fmt.Print("EMPTY\t\t' '\n")
	fmt.Print("BOX\t\t'$'\n")
	fmt.Print("FIGURE\t\t'x'\n")
	fmt.Print("EMPTY POINT\t'*'\n")
	fmt.Print("BOX POINT\t'%'\n")
	fmt.Print("FIGURE POINT\t'+'\n")
	fmt.Print("WALL\t\t'#'\n")
	fmt.Print("DEAD FIELD\t'☠'\n")
}

/* try moving figure in specified direction.
 * Returns, if figure was moved and if figure moved a box.
 */
func (e *Engine) Move(dir Direction) (success bool) {
	success = false

	dest1 := e.CurrentState.Figure.Add(dir.Point())
	valid, containsBox := e.checkDestination(dest1)
	if !valid {
		return
	}
	var dest2 Point
	if valid && containsBox {
		dest2 = dest1.Add(dir.Point())
		valid, containsSecBox := e.checkDestination(dest2)
		if !valid || containsSecBox {
			return
		}
	}

	success = true
	e.appendState2History(e.CurrentState)

	if containsBox {
		e.moveBox(dest1, dest2)
	}
	e.CurrentState.Figure = dest1
	return
}

// undo the last step (move figure and box to their old positions)
func (e *Engine) UndoStep() {
	l := len(e.History)
	if l == 0 {
		return
	}
	e.CurrentState = e.History[l-1]
	e.History = e.History[0 : l-1]
}

func (e *Engine) moveBox(src Point, dest Point) {
	delete(e.CurrentState.Boxes, src)
	e.CurrentState.Boxes[dest] = true
	//sort.Sort(e.CurrentState.Boxes)
}

func (e *Engine) appendState2History(state State) {
	newState := State{state.Boxes, state.Figure}
	//for key, _ := range state.Boxes {
	//	newState.Boxes[key] = true
	//}
	//newState.Figure = state.Figure
	e.History = append(e.History, newState)
}

func (e *Engine) checkDestination(dest Point) (valid bool, containsBox bool) {
	valid = false
	containsBox = false

	if e.Surface[dest.Y][dest.X] == WALL {
		return
	}

	if !e.Surface.In(dest) {
		return
	}

	valid = true
	containsBox = e.CurrentState.Boxes.Contains(dest)
	return
}

// loop over all points and check, if there is a box. Else return false
func (e *Engine) Won() bool {
	for key, _ := range e.Targets {
		if !e.CurrentState.Boxes[key] {
			return false
		}
	}
	return true
}

// print the current Surface
func (e *Engine) Print() {
	var buffer bytes.Buffer
	for y := 0; y < len(e.Surface); y++ {
		for x := 0; x < len(e.Surface[y]); x++ {
			p := NewPoint(x, y)
			if e.Surface[y][x] == WALL {
				buffer.WriteString("#")
			} else if e.CurrentState.Figure.X == x && e.CurrentState.Figure.Y == y {
				if e.Targets.Contains(p) {
					buffer.WriteString("+")
				} else {
					buffer.WriteString("x")
				}
			} else if e.CurrentState.Boxes.Contains(p) {
				if e.Targets.Contains(p) {
					buffer.WriteString("%%")
				} else {
					buffer.WriteString("$")
				}
			} else {
				if e.Targets.Contains(p) {
					buffer.WriteString("*")
				} else if e.Surface[y][x] == DEAD {
					buffer.WriteString("☠")
				} else {
					buffer.WriteString(" ")
				}
			}
			buffer.WriteString(" ")
		}
		buffer.WriteString("\n")
	}
	fieldnr, deadnr := e.Surface.AmountOfFields()
	buffer.WriteString(fmt.Sprintf("Boxes: %d\n", len(e.CurrentState.Boxes)))
	buffer.WriteString(fmt.Sprintf("Points: %d\n", len(e.Targets)))
	buffer.WriteString(fmt.Sprintf("Fields: %d\n", fieldnr))
	buffer.WriteString(fmt.Sprintf("DeadFields: %d\n", deadnr))

	fmt.Print(buffer.String())
}

// check if the surface border was reached
func (surface *Surface) In(p Point) bool {
	if p.Y < 0 || p.X < 0 || p.Y >= len(*surface) || p.X >= len((*surface)[p.Y]) {
		return false
	}
	return true
}

func (surface *Surface) AmountOfFields() (fields int, dead int8) {
	for y := 0; y < len(*surface); y++ {
		for x := 0; x < len((*surface)[y]); x++ {
			if (*surface)[y][x] != WALL {
				fields++
			}
			if (*surface)[y][x] != DEAD {
				dead++
			}
		}
	}
	return
}
