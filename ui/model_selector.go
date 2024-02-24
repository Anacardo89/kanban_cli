package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int
type updateFlag int

const (
	menuState modelState = iota
	projectState
	cardState
	labelState
)

const (
	upNone updateFlag = iota
	upMenu
	upProject
	upLabel
	upCard
)

// Implements tea.Model
type Selector struct {
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

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmd = nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
	case updateFlag:
		s.update = msg
	case modelState:
		s.state = msg
	}
	if s.update != upNone {
		cmd = s.checkUpFlag()
		return s, cmd
	}
	cmd = s.checkState(msg)
	return s, cmd
}

func (s Selector) View() string {
	switch s.state {
	case menuState:
		return s.menu.View()
	case projectState:
		return s.project.View()
	case labelState:
		return s.label.View()
	case cardState:
		return s.card.View()
	}
	return ""
}

// **************************************
func Test() Selector {
	return Selector{
		state: menuState,
		menu:  TestData(),
	}
}

// **************************************

// called by main
func New() Selector {
	return Selector{
		state: menuState,
		menu:  NewMenu(),
		sp:    nil,
		sc:    nil,
		sl:    nil,
	}
}

// Update
func (s *Selector) checkUpFlag() tea.Cmd {
	switch s.update {
	case upMenu:
		s.update = upNone
		s.menu.UpdateMenu()
		return func() tea.Msg { return menuState }
	case upProject:
		s.update = upNone
		s.project.UpdateProject()
		return func() tea.Msg { return projectState }
	case upLabel:
		if s.sc == nil {
			return func() tea.Msg { return labelState }
		}
		s.update = upNone
		s.sl = s.label.getLabel()
		s.sc.AddLabel(s.sl)
		s.sl = nil
		return func() tea.Msg { return upCard }
	case upCard:
		s.update = upNone
		s.card.UpdateCard()
		return func() tea.Msg { return cardState }
	}
	return nil
}

func (s *Selector) checkState(msg tea.Msg) tea.Cmd {
	switch s.state {
	case menuState:
		s.sp = nil
		s.sc = nil
		updatedMenu, cmd := s.menu.Update(msg)
		s.menu = updatedMenu.(Menu)
		return cmd
	case projectState:
		s.sc = nil
		if s.sp != s.menu.getProject() {
			s.sp = s.menu.getProject()
			s.project = OpenProject(s.sp)
			s.label = OpenLabels(s.sp)
		}
		updatedProject, cmd := s.project.Update(msg)
		s.project = updatedProject.(Project)
		return cmd
	case labelState:
		updatedLabel, cmd := s.label.Update(msg)
		s.label = updatedLabel.(Label)
		return cmd
	case cardState:
		if s.sc != s.project.getCard() {
			s.sc = s.project.getCard()
			s.card = OpenCard(s.sc)
		}
		updatedCard, cmd := s.card.Update(msg)
		s.card = updatedCard.(Card)
		return cmd
	}
	return nil
}
