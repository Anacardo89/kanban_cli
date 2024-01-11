package main

import (
	"log"

	"github.com/Anacardo89/kanban_cli/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("log.log", "error")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	menu := TestData()
// 	sm := menu.ToStorage()
// 	fmt.Println(sm.ToYAML())
// }
