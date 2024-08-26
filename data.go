package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

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

	if _, err := os.Stat("tasks.csv"); os.IsNotExist(err) {
		file, err := os.Create("tasks.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		file.WriteString("title,description,status\n")
	}

	content, err := os.ReadFile("tasks.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(content)))
	r.Comma = ','

	// Skip the header
	_, readErr := r.Read()
	if readErr != nil {
		fmt.Println("Error reading CSV header:", err)
	}

	records, readAllErr := r.ReadAll()
	if readAllErr != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
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

func updateCSV() {
	file, err := os.Create("tasks.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("title,description,status\n")
	w := csv.NewWriter(file)
	for _, col := range board.cols {
		for _, item := range col.list.Items() {
			task := item.(Task)
			w.Write([]string{task.Title(), task.Description(), task.status.String()})
		}
	}
	w.Flush()
}
