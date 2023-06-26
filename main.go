package kancli

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

const margin = 4

const (
	todo status = iota
	inProgress
	done
)

// How to Use:
// func main() {
// 	f, err := tea.LogToFile("debug.log", "debug")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	defer f.Close()
//
// 	board = NewBoard()
// 	board.initLists()
// 	p := tea.NewProgram(board)
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }
