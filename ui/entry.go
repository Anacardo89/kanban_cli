package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int

const (
	menu modelState = iota
	project
	card
)

type WindowSize struct {
	width  int
	height int
}

var ws WindowSize

type Model struct {
	state   modelState
	menu    Menu
	project Project
	sp      *kanban.Project
	card    tea.Model
}

func New() Model {
	return Model{
		state: menu,
		menu:  NewMenu(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case modelState:
		m.state = msg
	}
	switch m.state {
	case menu:
		updatedMenu, cmd := m.menu.Update(msg)
		m.menu = updatedMenu.(Menu)
		return m, cmd
	case project:
		m.sp = m.menu.selected
		m.project = OpenProject(m.sp)
		updatedProject, cmd := m.project.Update(msg)
		m.project = updatedProject.(Project)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return m.menu.View()
	case project:
		return m.project.View()
	}
	return ""
}
