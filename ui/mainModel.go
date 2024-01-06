package ui

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
	witdh      int
	height     int
	menu       *kanban.Menu
	cursor     int
	inputField textinput.Model
	input      string
	current    *dll.Node
}

func New() MainModel {
	return MainModel{
		cursor:     0,
		menu:       kanban.StartMenu(),
		inputField: textinput.New(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.Boards.GetLength() > 0 {
		m.current, _ = m.menu.Boards.WalkTo(m.cursor)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.witdh = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.menu.Boards.GetLength() == 0 {
				return m, nil
			}
			if m.cursor == 0 {
				m.cursor = m.menu.Boards.GetLength() - 1
				m.current, _ = m.menu.Boards.WalkTo(m.cursor)
				return m, nil
			}
			m.cursor--
			m.current, _ = m.menu.Boards.WalkTo(m.cursor)
			return m, nil
		case "down":
			if m.menu.Boards.GetLength() == 0 {
				return m, nil
			}
			if m.cursor == m.menu.Boards.GetLength()-1 {
				m.cursor = 0
				m.current, _ = m.menu.Boards.WalkTo(m.cursor)
				return m, nil
			}
			m.cursor++
			m.current, _ = m.menu.Boards.WalkTo(m.cursor)
			return m, nil
		case "n":
			m.inputField.Placeholder = "Name the Project: "
			m.inputField.Focus()
		}
	}
	return m, nil
}

func (m MainModel) View() string {
	output := ""
	if m.witdh == 0 {
		output += "Loading..."
		return output
	}
	if m.menu.Boards.GetLength() == 0 {
		output += "No projects currently.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		style := EmptyStyle()
		output = style.Render(output)
		return lipgloss.Place(m.witdh, m.height, lipgloss.Center, lipgloss.Center, output)
	}

	return output
}
