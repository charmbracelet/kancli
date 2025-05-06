package main

import (
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// TODOs
// - clean up textinput & area placeholder text (form could use some razzle dazzle) - use huh for this instead?
// - get accurate help menu options (not showing options properly)
// - move board to its own file

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
		// need to trigger what used to happen when a Form was received as msg.
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

/* Board Stuff */

type EditFormMsg struct {
	title       string
	description string
	index       int
	column      column
}

type NewFormMsg struct {
	index  int
	column column
}

type Board struct {
	help     help.Model
	loaded   bool
	focused  status
	cols     []column
	quitting bool
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (m *Board) Update(msg tea.Msg) (*Board, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width - margin
		for i := range m.cols {
			m.cols[i], cmd = m.cols[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case ReturnBoardMsg:
		if !msg.abort {
			return m, m.cols[m.focused].Set(msg.form.index, msg.form.CreateTask())
		}
	case moveMsg:
		return m, m.cols[m.focused.getNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getPrev()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	}
	m.cols[m.focused], cmd = m.cols[m.focused].Update(msg)
	return m, cmd
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}
