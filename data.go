package main

import (
	"github.com/charmbracelet/bubbles/list"
)

func stringToStatus(s string) status {
	switch s {
	case "todo":
		return todo
	case "inProgress":
		return inProgress
	case "done":
		return done
	default:
		return todo
	}
}

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	records := readCSV()

	listItems := [][]list.Item{
		todo:       {},
		inProgress: {},
		done:       {},
	}

	// Extract all with status todo
	for _, record := range records {
		if len(record) >= 3 {
			s := stringToStatus(record[2])
			listItems[s] = append(listItems[s], Task{s, record[0], record[1]})
		}
	}

	b.cols[todo].list.Title = "To Do"
	b.cols[todo].list.SetItems(listItems[todo])

	b.cols[inProgress].list.Title = "In Progress"
	b.cols[inProgress].list.SetItems(listItems[inProgress])

	b.cols[done].list.Title = "Done"
	b.cols[done].list.SetItems(listItems[done])
}
