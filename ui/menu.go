package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Implements tea.Model
type Menu struct {
	menu   *kanban.Menu
	cursor int
	list   list.Model
	Input  InputField
}

func NewMenu() Menu {
	m := Menu{
		menu:  kanban.StartMenu(),
		Input: InputField{field: textinput.New()},
	}
	setMenuItemDelegate()
	m.setupList()
	return m
}

func (m *Menu) UpdateMenu() {
	m.setupList()
}

func TestData() Menu {
	return Menu{
		cursor: 0,
		menu:   kanban.TestData(),
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
		updateWindowSize(msg)
		m.setupList()
		return m, nil
	case tea.KeyMsg:
		if m.Input.field.Focused() {
			m.handleInput(msg.String())
			m.Input.field, cmd = m.Input.field.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			p := m.getProject()
			if p == nil {
				return m, nil
			}
			return m, func() tea.Msg { return project }
		case "n":
			m.setInput()
			return m, m.Input.field.Focus()
		case "d":
			p := m.getProject()
			if p == nil {
				return m, nil
			}
			m.deleteProject()
			return m, nil
		}
	}
	m.list, cmd = m.list.Update(msg)
	m.cursor = m.list.Cursor()
	return m, cmd
}

func (m Menu) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	var (
		bottomLines    = ""
		emptyTxtStyled = ""
		inputStyled    = ""
		menuStyled     = ""
		output         = ""
	)

	if m.menu.Projects.Length() == 0 {
		emptyTxt := "No projects.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		emptyTxtStyled = EmptyStyle.Render(emptyTxt)
		if m.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2; i++ {
				bottomLines += "\n"
			}
			inputStyled = InputFieldStyle.Render(m.Input.field.View())
		}
		output = lipgloss.Place(
			ws.width,
			ws.height,
			lipgloss.Center,
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Center,
				emptyTxtStyled,
				bottomLines,
				inputStyled,
			))
		return output
	}

	menuStyled = ListStyle.Render(m.list.View())
	if m.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(m.Input.field.View())
	}
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			menuStyled,
			bottomLines,
			inputStyled,
		))
	return output
}

// action
func (m *Menu) getProject() *kanban.Project {
	if m.menu.Projects.Length() == 0 {
		return nil
	}
	node, err := m.menu.Projects.WalkTo(m.cursor)
	if err != nil {
		log.Println(err)
	}
	return node.Val().(*kanban.Project)
}

func (m *Menu) deleteProject() {
	var err error
	if m.menu.Projects.Length() == 0 {
		return
	}
	node, err := m.menu.Projects.WalkTo(m.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	p := node.Val().(*kanban.Project)
	err = m.menu.RemoveProject(p)
	if err != nil {
		log.Println(err)
	}
	m.setupList()
	m.cursor = m.list.Cursor()
}
