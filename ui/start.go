package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int
type updateFlag int
type inputFlag string

const (
	menu modelState = iota
	project
	card
	label
)

const (
	upNone updateFlag = iota
	upMenu
	upProject
	upLabel
	upCard
)

const (
	none   inputFlag = "none"
	new    inputFlag = "n"
	add    inputFlag = "a"
	rename inputFlag = "r"
	move   inputFlag = "m"
	delete inputFlag = "d"
	color  inputFlag = "color"
	title  inputFlag = "title"
)

type Model struct {
	state   modelState
	update  updateFlag
	menu    Menu
	project Project
	sp      *kanban.Project
	card    Card
	sc      *kanban.Card
	label   Label
	sl      *kanban.Label
}

func New() Model {
	return Model{
		state: menu,
		menu:  NewMenu(),
		sp:    nil,
		sc:    nil,
		sl:    nil,
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
	case updateFlag:
		m.update = msg
	case modelState:
		m.state = msg
	}

	switch m.update {
	case upMenu:
		m.update = upNone
		m.menu.UpdateMenu()
		return m, func() tea.Msg { return menu }
	case upProject:
		m.update = upNone
		m.project.UpdateProject()
		return m, func() tea.Msg { return project }
	case upLabel:
		if m.sc == nil {
			return m, nil
		}
		m.update = upNone
		m.sl = m.label.getLabel()
		m.sc.AddLabel(m.sl)
		m.sl = nil
		return m, func() tea.Msg { return upCard }
	case upCard:
		m.update = upNone
		m.card.UpdateCard()
		return m, func() tea.Msg { return card }
	}

	switch m.state {
	case menu:
		m.sp = nil
		m.sc = nil
		updatedMenu, cmd := m.menu.Update(msg)
		m.menu = updatedMenu.(Menu)
		return m, cmd
	case project:
		m.sc = nil
		if m.sp != m.menu.getProject() {
			m.sp = m.menu.getProject()
			m.project = OpenProject(m.sp)
			m.label = OpenLabels(m.sp)
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
