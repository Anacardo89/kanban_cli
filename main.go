package main

import (
	"log"

	"github.com/Anacardo89/kanboards/logger"
	"github.com/Anacardo89/kanboards/storage"
	"github.com/Anacardo89/kanboards/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := logger.CreateLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if storage.DBExists() {
		storage.OpenDB()
	} else {
		storage.CreateDBfile()
		storage.CreateDBTables()
	}
	defer storage.DB.Close()

	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		logger.Error.Fatal(err)
	}
}
