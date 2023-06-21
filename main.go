package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: fix adding an item to list

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

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

func main() {
	board := NewBoard()
	board.initLists()
	p := tea.NewProgram(board)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
