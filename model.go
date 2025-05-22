package main

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Model struct {
	// false is board, true is form
	modifying bool
	board     *Board
	form      Form
}

func newModel() Model {
	return Model{
		board: NewBoard().initLists(),
		form:  newDefaultForm(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ReturnBoardMsg:
		m.modifying = false
	case NewFormMsg:
		m.modifying = true
		m.form = newDefaultForm()
		m.form.index = APPEND
		m.form.col = msg.column
	case EditFormMsg:
		m.modifying = true
		m.form = NewForm(msg.title, msg.description)
		m.form.index = msg.index
		m.form.col = msg.column
	}

	if m.modifying {
		m.form, cmd = m.form.Update(msg)
		return m, cmd
	}
	m.board, cmd = m.board.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.modifying {
		return m.form.View()
	}
	return m.board.View()
}
