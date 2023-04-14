package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

/* MODEL MANAGEMENT */
var models []tea.Model

// edit constant
const (
	noEdit = -1
)

const (
	board status = iota
	form
)

/* STYLING */
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder())
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
)

/* CUSTOM ITEM */

type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

/* MAIN MODEL */

type Model struct {
	loaded       bool
	focused      status
	lists        []list.Model
	quitting     bool
	editingIndex int
}

func New() *Model {
	return &Model{editingIndex: noEdit}
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		columnStyle.Width(msg.Width / divisor)
		focusedStyle.Width(msg.Width / divisor)
		columnStyle.Height(msg.Height - divisor)
		focusedStyle.Height(msg.Height - divisor)
		for i, list := range m.lists {
			list.SetSize(msg.Width/divisor, msg.Height/2)
			m.lists[i], _ = list.Update(msg)
		}
		m.loaded = true
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		case "n":
			models[board] = m // save the state of the current model
			models[form] = NewForm(m.focused)
			return models[form].Update(nil)
		case "e":
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
		case "d":
			return m, m.DeleteCurrent
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
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}
	} else {
		return "loading..."
	}
}

/* FORM MODEL */
type Form struct {
	focused     status
	title       textinput.Model
	description textarea.Model
}

func NewForm(focused status) *Form {
	form := Form{
		focused:     focused,
		title:       textinput.New(),
		description: textarea.New(),
	}

	form.title.Focus()
	return &form
}

func (m Form) CreateTask() tea.Msg {
	return Task{
		m.focused, m.title.Value(), m.description.Value(),
	}
}

func (m Form) Init() tea.Cmd {
	return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				models[form] = m
				return models[board], m.CreateTask
			}
		}
	}
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}

func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}

func main() {
	boardView := Model{}
	boardView.initLists()
	models = []tea.Model{&boardView, NewForm(todo)}
	m := models[board]
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
