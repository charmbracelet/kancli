package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
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

const margin = 4

var board *Board

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

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
