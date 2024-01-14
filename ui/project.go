package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
)

type Project struct {
	witdh   int
	height  int
	project *kanban.Project
	list    []list.Model
	label   list.Model
	cursor  int
	Input   InputField
}
