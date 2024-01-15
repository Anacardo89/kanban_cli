package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputField struct {
	field textinput.Model
	data  string
}

type Menu struct {
	witdh  int
	height int
	menu   *kanban.Menu
	list   list.Model
	cursor int
	Input  InputField
}

func NewMenu() Menu {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 20)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	return Menu{
		cursor: 0,
		menu:   kanban.StartMenu(),
		Input:  InputField{field: textinput.New()},
		list:   l,
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
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		case "down":
			m.handleMoveDown()
			m.list, cmd = m.list.Update(msg)
			return m, cmd
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
		output = "No projects.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		return lipgloss.Place(m.witdh, m.height, lipgloss.Center, lipgloss.Center, emptyStyle.Render(output))
	}
	if m.menu.Projects.GetLength() > 0 {
		output = lipgloss.Place(0, 0, lipgloss.Left, lipgloss.Top, menuListStyle.Render(m.list.View()))
	}
	if m.Input.field.Focused() {
		output = lipgloss.Place(m.witdh, m.height, lipgloss.Left, lipgloss.Bottom,
			lipgloss.JoinVertical(
				lipgloss.Left,
				inputStyle.Render(m.Input.field.View()),
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
		return
	}
	m.cursor--
}

func (m *Menu) handleMoveDown() {
	if m.menu.Projects.GetLength() == 0 {
		return
	}
	if m.cursor == m.menu.Projects.GetLength()-1 {
		m.cursor = 0
		return
	}
	m.cursor++
}

func (m *Menu) setInput() {
	m.Input.field.Prompt = ": "
	m.Input.field.CharLimit = 120
	m.Input.field.Placeholder = "Project Title"
}

func (m *Menu) handleInput(key string) {
	switch key {
	case "esc":
		m.Input.field.SetValue("")
		m.Input.data = ""
		m.Input.field.Blur()
		return
	case "enter":
		log.Println(m.Input.field.Value())
		m.Input.data = m.Input.field.Value()
		menuItem := menuItem{
			title: m.Input.data,
		}
		menuItems = append(menuItems, menuItem)
		m.list.SetItems(menuItems)
		m.menu.AddProject(m.Input.data)
		m.Input.data = ""
		m.Input.field.Blur()
		return
	}
}
