package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	help         help.Model
	keys         keyMap
	loaded       bool
	focused      status
	lists        []list.Model
	quitting     bool
	editingIndex int
}

func New() *Model {
	return &Model{keys: keys, help: help.New(), editingIndex: noEdit}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil { // will happen if board is empty
		return nil
	}
	selectedTask := selectedItem.(Task)
	m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
	return nil
}

func (m *Model) DeleteCurrent() tea.Msg {
	if len(m.lists[m.focused].VisibleItems()) > 0 {
		selectedTask := m.lists[m.focused].SelectedItem().(Task)
		m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
	}
	return nil
}

func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists() {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	// Init To Do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) calculateHeight(height int) {
	columnStyle.Height(height - height/divisor)
	focusedStyle.Height(height - height/divisor)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		columnStyle.Width(msg.Width / divisor)
		focusedStyle.Width(msg.Width / divisor)
		m.help.Width = msg.Width
		m.calculateHeight(msg.Height)
		for i, list := range m.lists {
			list.SetSize(msg.Width/divisor, msg.Height/2)
			m.lists[i], _ = list.Update(msg)
		}
		m.loaded = true
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Left):
			m.Prev()
		case key.Matches(msg, m.keys.Right):
			m.Next()
		case key.Matches(msg, m.keys.Enter):
			return m, m.MoveToNext
		case key.Matches(msg, m.keys.New):
			models[board] = m // save the state of the current model
			models[form] = NewForm(m.focused)
			return models[form].Update(nil)
		case key.Matches(msg, m.keys.Edit):
			list := m.lists[m.focused]
			if len(list.VisibleItems()) == 0 {
				return m, nil
			}
			task := list.SelectedItem().(Task)
			editForm := NewForm(m.focused)
			editForm.title.SetValue(task.title)
			editForm.description.SetValue(task.description)
			m.editingIndex = list.Index()
			models[board] = m // save the state of the current model
			models[form] = editForm
			return models[form].Update(nil)
		case key.Matches(msg, m.keys.Delete):
			return m, m.DeleteCurrent
		case key.Matches(msg, m.keys.Help):
			if m.help.ShowAll {
				m.help.ShowAll = false
			} else {
				m.help.ShowAll = true
			}
		}
	case Task:
		task := msg
		list := &m.lists[task.status]
		// if edit, replace existing task in list
		if m.editingIndex != noEdit {
			index := m.editingIndex
			m.editingIndex = noEdit
			return m, list.SetItem(index, task)
		}
		// add task to end of list
		return m, list.InsertItem(len(list.Items()), task)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	todoView := m.lists[todo].View()
	inProgView := m.lists[inProgress].View()
	doneView := m.lists[done].View()
	var board string
	switch m.focused {
	case inProgress:
		board = lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgView),
			columnStyle.Render(doneView),
		)
	case done:
		board = lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgView),
			focusedStyle.Render(doneView),
		)
	default:
		board = lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgView),
			columnStyle.Render(doneView),
		)
	}
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(m.keys))
}
