package ui

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputField struct {
	field textinput.Model
	data  string
}

type Menu struct {
	witdh    int
	height   int
	menu     *kanban.Menu
	cursor   int
	selected *dll.Node
	Input    InputField
}

func NewMenu() Menu {
	return Menu{
		cursor: 0,
		menu:   kanban.StartMenu(),
		Input:  InputField{field: textinput.New()},
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.witdh = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.Input.field.Focused() {
			m.handleInput(msg.String())
			m.Input.field, cmd = m.Input.field.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.handleMoveUp()
			return m, nil
		case "down":
			m.handleMoveDown()
			return m, nil
		case "n":
			m.setInput()
			return m, m.Input.field.Focus()
		}
	}
	return m, nil
}

func (m Menu) View() string {
	output := ""
	if m.menu.Projects.GetLength() == 0 && !m.Input.field.Focused() {
		output += "No projects.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		style := EmptyStyle()
		return lipgloss.Place(m.witdh, m.height, lipgloss.Center, lipgloss.Center, style.Render(output))
	}
	if m.Input.field.Focused() {
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

func (m *Menu) handleMoveUp() {
	if m.menu.Projects.GetLength() == 0 {
		return
	}
	if m.cursor == 0 {
		m.cursor = m.menu.Projects.GetLength() - 1
		m.selected, _ = m.menu.Projects.WalkTo(m.cursor)
		return
	}
	m.cursor--
	m.selected, _ = m.selected.Prev()
}

func (m *Menu) handleMoveDown() {
	if m.menu.Projects.GetLength() == 0 {
		return
	}
	if m.cursor == m.menu.Projects.GetLength()-1 {
		m.cursor = 0
		m.selected, _ = m.menu.Projects.WalkTo(m.cursor)
		return
	}
	m.cursor++
	m.selected, _ = m.selected.Next()
}

func (m *Menu) setInput() {
	m.Input.field.Prompt = ": "
	m.Input.field.CharLimit = 120
	m.Input.field.Placeholder = "Project Title"
}

func (m *Menu) handleInput(key string) {
	switch key {
	case "esc":
		m.Input.field.Blur()
		m.Input.field.SetValue("")
		return
	case "enter":
		m.Input.data = m.Input.field.Value()
		m.menu.AddProject(m.Input.data)
		m.Input.data = ""
		m.Input.field.Blur()
		return
	}
}
