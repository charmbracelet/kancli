package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	focus  bool
	status status
	list   list.Model
	height int
	width  int
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status status) column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus, status, defaultList, 0, 0}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/divisor, msg.Height/2)
	case Form:
		return c, c.Set(msg.index, msg.CreateTask())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewForm(task.title, task.description)
				f.index = c.list.Index()
				f.col = c
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Msg {
	if len(c.list.VisibleItems()) > 0 {
		c.list.RemoveItem(c.list.Index())
	}
	return nil
}

func (c *column) Set(i int, t Task) tea.Cmd {
	return func() tea.Msg {
		if i != APPEND {
			c.list.SetItem(i, t)
		}
		return c.list.InsertItem(APPEND, t)
	}
}

func (c *column) setSize(width, height int) {
	// TODO: make height of col depend on num items in list?
	c.height = height / divisor
	c.width = height / width
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

func (c *column) MoveToNext() tea.Msg {
	task := c.list.SelectedItem().(Task)
	c.list.RemoveItem(c.list.Index())
	task.status = c.status.getNext()
	return task
}
