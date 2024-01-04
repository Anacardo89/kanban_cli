package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
	cursor int
	menu   *kanban.Menu
	input  textinput.Model
}

func New() MainModel {
	return MainModel{
		cursor: 0,
		menu:   kanban.StartMenu(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor == 0 {
				m.cursor = m.menu.Boards.GetLength() - 1
				return m, nil
			}
			m.cursor--
			return m, nil
		case "down":
			if m.cursor == m.menu.Boards.GetLength()-1 {
				m.cursor = 0
				return m, nil
			}
			m.cursor++
			return m, nil
		case "n":
		}
	}
	return m, nil
}

func (m MainModel) View() string {
	return "test"
}
