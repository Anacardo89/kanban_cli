package ui

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputField struct {
	field  textinput.Model
	data   string
	active bool
}

type MainModel struct {
	witdh   int
	height  int
	menu    *kanban.Menu
	cursor  int
	Input   InputField
	current *dll.Node
}

func New() MainModel {
	return MainModel{
		cursor: 0,
		menu:   kanban.StartMenu(),
		Input: InputField{
			field:  textinput.New(),
			active: false,
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.menu.Boards.GetLength() > 0 {
		m.current, _ = m.menu.Boards.WalkTo(m.cursor)
	}
	if m.Input.field.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.Input.field.Blur()
				m.Input.field.SetValue("")
				m.Input.active = false
				return m, nil
			case "enter":
				m.Input.data = m.Input.field.Value()
				m.menu.AddBoard(m.Input.data)
				m.Input.data = ""
				m.Input.field.Blur()
				return m, nil
			}
		}
		m.Input.field, cmd = m.Input.field.Update(msg)
		return m, cmd
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
			m.Input.field.Placeholder = "Project Title"
			m.Input.active = true
			return m, m.Input.field.Focus()
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
	if m.menu.Boards.GetLength() == 0 && !m.Input.active {
		output += "No projects currently.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		style := EmptyStyle()
		return lipgloss.Place(m.witdh, m.height, lipgloss.Center, lipgloss.Center, style.Render(output))
	}

	if m.Input.active {
		style := InputStyle()
		output += lipgloss.Place(
			m.witdh, m.height,
			lipgloss.Left, lipgloss.Bottom,
			lipgloss.JoinVertical(
				lipgloss.Left,
				style.Render(m.Input.field.View()),
			),
		)
		return output
	}
	return output
}
