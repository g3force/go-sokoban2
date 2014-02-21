package main

import (
	"fmt"
	"github.com/g3force/go-sokoban2/ai"
	"github.com/g3force/go-sokoban2/engine"
	"github.com/lye/curses"
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

func main() {
	logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} 0x%{id:x} %{message}"))

	// Setup one stdout and one syslog backend.
	logBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	logBackend.Color = true

	syslogBackend, err := logging.NewSyslogBackend("")
	if err != nil {
		log.Fatal(err)
	}

	// Combine them both into one logging backend.
	logging.SetBackend(logBackend, syslogBackend)

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

	var input string
	for {
		curses.Mvaddstr(0, 0, e.SurfaceToStr(e.CurrentState))
		curses.Mvaddstr(len(e.Surface)+5, 0, "w,a,s,d to control, q to quit")
		curses.Refresh()
		input = string(curses.Getch())
		//fmt.Scanf("%s", &input)
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
