package ui

import (
	"github.com/Anacardo89/kanboards/kanban"
	"github.com/Anacardo89/kanboards/storage"
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
	state  modelState
	update updateFlag
	m      Menu
	p      Project
	sp     *kanban.Project
	c      Card
	sc     *kanban.Card
	l      Label
	sl     *kanban.Label
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
		return s.m.View()
	case projectState:
		return s.p.View()
	case labelState:
		return s.l.View()
	case cardState:
		return s.c.View()
	}
	return ""
}

// called by main
func New() Selector {
	s := Selector{
		state: menuState,
		m:     NewMenu(),
		sp:    nil,
		sc:    nil,
		sl:    nil,
	}
	s.loadKanban()
	s.m.UpdateMenu()
	return s
}

// Update
func (s *Selector) checkUpFlag() tea.Cmd {
	switch s.update {
	case upMenu:
		s.update = upNone
		s.m.UpdateMenu()
		return func() tea.Msg { return menuState }
	case upProject:
		s.update = upNone
		s.p.UpdateProject()
		return func() tea.Msg { return projectState }
	case upLabel:
		if s.sc == nil {
			return func() tea.Msg { return labelState }
		}
		s.update = upNone
		s.sl = s.l.getLabel()
		cl, _ := s.sc.CardLabels.HeadNode()
		if cl != nil {
			storage.CreateCardLabel(s.sc.Id, s.sl.Id)
			s.sc.AddLabel(s.sl)
			s.sl = nil
			return func() tea.Msg { return upCard }
		}
		for i := 0; i < s.sc.CardLabels.Length(); i++ {
			clVal := cl.Val().(*kanban.Label)
			if s.sl.Id == clVal.Id {
				s.sl = nil
				return func() tea.Msg { return upCard }
			}
			cl, _ = cl.Next()
		}
		storage.CreateCardLabel(s.sc.Id, s.sl.Id)
		s.sc.AddLabel(s.sl)
		s.sl = nil
		return func() tea.Msg { return upCard }
	case upCard:
		s.update = upNone
		s.c.UpdateCard()
		return func() tea.Msg { return cardState }
	}
	return nil
}

func (s *Selector) checkState(msg tea.Msg) tea.Cmd {
	switch s.state {
	case menuState:
		s.sp = nil
		s.sc = nil
		updatedMenu, cmd := s.m.Update(msg)
		s.m = updatedMenu.(Menu)
		return cmd
	case projectState:
		s.sc = nil
		if s.sp != s.m.getProject() {
			s.sp = s.m.getProject()
			s.p = OpenProject(s.sp)
			s.l = OpenLabels(s.sp)
		}
		updatedProject, cmd := s.p.Update(msg)
		s.p = updatedProject.(Project)
		return cmd
	case labelState:
		updatedLabel, cmd := s.l.Update(msg)
		s.l = updatedLabel.(Label)
		return cmd
	case cardState:
		if s.sc != s.p.getCard() {
			s.sc = s.p.getCard()
			s.c = OpenCard(s.sc)
		}
		updatedCard, cmd := s.c.Update(msg)
		s.c = updatedCard.(Card)
		return cmd
	}
	return nil
}
