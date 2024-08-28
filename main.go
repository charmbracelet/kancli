package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type status int

func (s status) getNext() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == todo {
		return done
	}
	return s - 1
}

func (s status) String() string {
	switch s {
	case todo:
		return "todo"
	case inProgress:
		return "inProgress"
	case done:
		return "done"
	default:
		return "unknown"
	}
}

const margin = 4

var (
	board   *Board
	csvFile string
)

const (
	todo status = iota
	inProgress
	done
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	config := readConfig()
	csvFile = config.DbPath

	board = NewBoard()
	board.initLists()
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
