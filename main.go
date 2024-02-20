package main

import (
	"log"

	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.CreateLogger()
	f, err := tea.LogToFile("ui.log", "error")
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
