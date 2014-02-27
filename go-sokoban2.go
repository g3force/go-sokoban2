package main

import (
	"fmt"
	"github.com/g3force/go-sokoban2/ai"
	"github.com/g3force/go-sokoban2/engine"
	curses "github.com/lye/curses"
	"github.com/op/go-logging"
	stdlog "log"
	"os"
	"strconv"
)

const (
	LOGGER_NAME  = "sokoban"
	windowHeight = 20
	windowWidth  = 40
)

var log = logging.MustGetLogger(LOGGER_NAME)

type CursesLogBackend struct {
	xOffset int
	cursor  *int
}

func (b CursesLogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	maxy, _ := curses.Getmaxyx()
	s := rec.Formatted()
	curses.Mvaddstr(*b.cursor, b.xOffset, s)
	curses.Mvaddstr((*b.cursor)+1, b.xOffset, "---------------------------------")
	(*b.cursor)++
	if *b.cursor > maxy {
		(*b.cursor) = 0
	}
	return nil
}

func main() {
	runmode := false
	single := true
	level := "res/alevel"
	straightAhead := false
	outputFreq := int32(50000)
	printSurface := false
	threads := 1

	if len(os.Args) > 1 {
		for i, _ := range os.Args {
			switch os.Args[i] {
			case "-r":
				runmode = true
			case "-l":
				if len(os.Args) > i+1 {
					level = os.Args[i+1]
				}
			case "-i":
				engine.PrintInfo()
			case "-m":
				single = false
			case "-d":
				if len(os.Args) > i+1 {
					debuglevel, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						switch debuglevel {
						case 0:
							logging.SetLevel(logging.CRITICAL, LOGGER_NAME)
						case 1:
							logging.SetLevel(logging.ERROR, LOGGER_NAME)
						case 2:
							logging.SetLevel(logging.WARNING, LOGGER_NAME)
						case 3:
							logging.SetLevel(logging.INFO, LOGGER_NAME)
						case 4:
							logging.SetLevel(logging.DEBUG, LOGGER_NAME)
						}
					}
				}
			case "-s":
				straightAhead = true
			case "-f":
				if len(os.Args) > i+1 {
					of, err := strconv.Atoi(os.Args[i+1])
					if err == nil {
						outputFreq = int32(of)
					}
				}
			case "-p":
				printSurface = true
			case "-t":
				if len(os.Args) > i+1 {
					t, err := strconv.Atoi(os.Args[i+1])
					if err != nil {
						panic(err)
					} else {
						threads = t
					}
				}
			}
		}
	}

	//logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} 0x%{id:x} %{message}"))
	logging.SetFormatter(logging.MustStringFormatter("%{level:.1s} %{message}"))

	fo, err := os.Create("sokoban.log")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	// Setup one stdout and one syslog backend.
	//logBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	logBackend := logging.NewLogBackend(fo, "", stdlog.Lshortfile)
	logBackend.Color = true

	// Combine them both into one logging backend.
	logging.SetBackend(logBackend)

	e := engine.NewEngine(level)
	ai.MarkDeadFields(e)
	log.Info("Level: " + level)

	_ = runmode // todo
	_ = single
	_ = straightAhead
	_ = outputFreq
	_ = printSurface
	_ = threads
	//if runmode {
	//	ai.Run(e, single, outputFreq, printSurface, straightAhead, threads)
	//	return
	//}

	fmt.Println("Press m for manual or r for run.")
	for {
		var choice string
		fmt.Scanf("%s", &choice)

		if choice == "r" {
			//ai.Run(e, single, outputFreq, printSurface, straightAhead, threads)
			e.Print()
			ai.Solve(e)
			return
		} else if choice == "m" {
			break
		}
	}

	// init curses GUI
	curses.Initscr()
	defer curses.End()
	curses.Cbreak()
	curses.Noecho()
	curses.Stdscr.Keypad(true)

	//cursor := new(int)
	//clb := CursesLogBackend{, cursor}
	//xOffset := 7 + 2*len(e.Surface[0])
	//logging.SetBackend(clb)

	var input string
	for {
		curses.Mvaddstr(0, 0, e.SurfaceToStr(e.CurrentState))
		curses.Mvaddstr(len(e.Surface)+5, 0, "w,a,s,d to control, q to quit")
		curses.Refresh()
		log.Debug("blubb")
		input = string(curses.Getch())
		switch input {
		case "d":
			e.Move(0)
		case "s":
			e.Move(1)
		case "a":
			e.Move(2)
		case "w":
			e.Move(3)
		case "q":
			return
		default:
			e.UndoStep()
		}
	}
}
