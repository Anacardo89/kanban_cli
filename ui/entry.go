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
	label
)

type inputFlag string

const (
	none   inputFlag = "none"
	new    inputFlag = "n"
	add    inputFlag = "a"
	rename inputFlag = "r"
	move   inputFlag = "m"
	delete inputFlag = "d"
)

type Model struct {
	state   modelState
	menu    Menu
	project Project
	sp      *kanban.Project
	card    Card
	sc      *kanban.Card
	label   Label
}

func New() Model {
	return Model{
		state: menu,
		menu:  NewMenu(),
	}
}

func Test() Model {
	return Model{
		state: menu,
		menu:  TestData(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
	case modelState:
		m.state = msg
	}
	switch m.state {
	case menu:
		updatedMenu, cmd := m.menu.Update(msg)
		m.menu = updatedMenu.(Menu)
		return m, cmd
	case project:
		if m.sp != m.menu.getProject() {
			m.sp = m.menu.getProject()
			m.project = OpenProject(m.sp)
			m.label = OpenLabel(m.sp)
		}
		updatedProject, cmd := m.project.Update(msg)
		m.project = updatedProject.(Project)
		return m, cmd
	case label:
		updatedLabel, cmd := m.label.Update(msg)
		m.label = updatedLabel.(Label)
		return m, cmd
	case card:
		if m.sc != m.project.getCard() {
			m.sc = m.project.getCard()
			m.card = OpenCard(m.sc)
		}
		updatedCard, cmd := m.card.Update(msg)
		m.card = updatedCard.(Card)
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
	case label:
		return m.label.View()
	case card:
		return m.card.View()
	}
	return ""
}
