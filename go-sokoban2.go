package main

import (
	"fmt"
	"github.com/g3force/go-sokoban2/ai"
	"github.com/g3force/go-sokoban2/engine"
	stdlog "log"
	"os"
	"strconv"

	"github.com/op/go-logging"
)

const (
	LOGGER_NAME = "sokoban"
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

	// surface
	e.Print()

	var choice string
	for {
		choice = ""
		fmt.Println("Press m for manual or r for run: ")
		fmt.Scanf("%s", &choice)
		if choice == "r" {
			//ai.Run(e, single, outputFreq, printSurface, straightAhead, threads)
			break
		} else if choice == "m" {
			fmt.Println("Manual mode\n")
			var input string
			for {
				fmt.Scanf("%s", &input)
				switch input {
				case "0":
					e.Move(0)
				case "1":
					e.Move(1)
				case "2":
					e.Move(2)
				case "3":
					e.Move(3)
				default:
					e.UndoStep()
				}
				e.Print()
			}
			break
		}
	}
}
