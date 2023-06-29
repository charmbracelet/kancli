package kancli

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

const margin = 4

type Status interface {
	Next() int
	Prev() int
	String() string
}

type Column struct {
	focus  bool
	List   list.Model
	status Status
	height int
	width  int
	board  *Board
}

func (c *Column) Focus() {
	c.focus = true
}

func (c *Column) Blur() {
	c.focus = false
}

func (c *Column) Focused() bool {
	return c.focus
}

// NewColumn creates a new column from a list.
func NewColumn(l []list.Item, status Status, focus bool) Column {
	defaultList := list.New(l, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return Column{focus: focus, status: status, List: defaultList}
}

// Init does initial setup for the column.
func (c Column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case list.Item:
		return c, c.Set(APPEND, msg)
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.List.SetSize(msg.Width/margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext()
		}
	}
	c.List, cmd = c.List.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.getStyle().Render(c.List.View())
}

func (c *Column) DeleteCurrent() tea.Cmd {
	if len(c.List.VisibleItems()) > 0 {
		c.List.RemoveItem(c.List.Index())
	}

	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)
	return cmd
}

// Set adds an item to a column.
func (c *Column) Set(i int, item list.Item) tea.Cmd {
	if i != APPEND {
		return c.List.SetItem(i, item)
	}
	return c.List.InsertItem(APPEND, item)
}

func (c *Column) setSize(width, height int) {
	c.width = width / margin
}

func (c *Column) getStyle() lipgloss.Style {
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

// MoveMsg can be handled by the lib user to update the status of their items.
type MoveMsg struct {
	i    int
	item list.Item
}

// MoveToNext returns the new column index for the selected item.
func (c *Column) MoveToNext() tea.Cmd {
	// If nothing is selected, the SelectedItem will return Nil.
	item := c.List.SelectedItem()
	if item == nil {
		return nil
	}
	// move item
	c.List.RemoveItem(c.List.Index())

	// refresh list
	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return MoveMsg{c.status.Next(), item} })
}
