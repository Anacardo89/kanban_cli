package main

import (
	"log"

	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/storage"
	"github.com/Anacardo89/kanban_cli/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := logger.CreateLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if storage.DBExists() {
		storage.SetDB()
	} else {
		storage.CreateDBfile()
		storage.CreateDB()
	}

	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}
