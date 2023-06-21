package main

import "github.com/charmbracelet/bubbles/list"

// Provides the mock data to fill the kanban board

func (m *Board) initLists() {
	m.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	// Init To Do
	m.cols[todo].list.Title = "To Do"
	m.cols[todo].list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	m.cols[inProgress].list.Title = "In Progress"
	m.cols[inProgress].list.SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	m.cols[done].list.Title = "Done"
	m.cols[done].list.SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}
