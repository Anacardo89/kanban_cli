package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/storage"
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
	return s
}

func (s *Selector) loadKanban() {
	s.loadProjects()
	s.loadLabels()
	s.loadBoards()
	s.loadCards()
	s.loadCheckItems()
	s.loadCardLabels()
}

func (s *Selector) loadProjects() {
	for i := 0; i < len(storage.ProjectsSql); i++ {
		s.m.menu.AddProject(
			storage.ProjectsSql[i].Id,
			storage.ProjectsSql[i].Title,
		)
	}
}

func (s *Selector) loadLabels() {
	for i := 0; i < s.m.menu.Projects.Length(); i++ {
		project, err := s.m.menu.Projects.GetAt(i)
		if err != nil {
			logger.Error.Fatal(err)
		}
		p := project.(*kanban.Project)
		for j := 0; j < len(storage.LabelsSql); j++ {
			if p.Id == storage.LabelsSql[j].ProjectId {
				p.AddLabel(
					storage.LabelsSql[j].Id,
					storage.LabelsSql[j].Title,
					storage.LabelsSql[j].Color,
				)
			}
		}
	}
}

func (s *Selector) loadBoards() {
	for i := 0; i < s.m.menu.Projects.Length(); i++ {
		project, err := s.m.menu.Projects.GetAt(i)
		if err != nil {
			logger.Error.Fatal(err)
		}
		p := project.(*kanban.Project)
		for j := 0; j < len(storage.BoardsSql); j++ {
			if p.Id == storage.BoardsSql[j].ProjectId {
				p.AddBoard(
					storage.BoardsSql[j].Id,
					storage.BoardsSql[j].Title,
				)
			}
		}
	}
}

func (s *Selector) loadCards() {
	for i := 0; i < s.m.menu.Projects.Length(); i++ {
		project, err := s.m.menu.Projects.GetAt(i)
		if err != nil {
			logger.Error.Fatal(err)
		}
		p := project.(*kanban.Project)
		for j := 0; j < p.Boards.Length(); j++ {
			board, err := p.Boards.GetAt(j)
			if err != nil {
				logger.Error.Fatal(err)
			}
			b := board.(*kanban.Board)
			for k := 0; k < len(storage.CardsSql); k++ {
				if b.Id == storage.CardsSql[k].BoardId {
					b.AddCard(
						storage.CardsSql[k].Id,
						storage.CardsSql[k].Title,
						storage.CardsSql[k].Desc.String,
					)
				}
			}
		}
	}
}

func (s *Selector) loadCheckItems() {
	for i := 0; i < s.m.menu.Projects.Length(); i++ {
		project, err := s.m.menu.Projects.GetAt(i)
		if err != nil {
			logger.Error.Fatal(err)
		}
		p := project.(*kanban.Project)
		for j := 0; j < p.Boards.Length(); j++ {
			board, err := p.Boards.GetAt(i)
			if err != nil {
				logger.Error.Fatal(err)
			}
			b := board.(*kanban.Board)
			for k := 0; k < b.Cards.Length(); k++ {
				card, err := b.Cards.GetAt(i)
				if err != nil {
					logger.Error.Fatal(err)
				}
				c := card.(*kanban.Card)
				if c.Id == storage.CheckItemsSql[k].CardId {
					done := false
					if storage.CheckItemsSql[k].Done == 1 {
						done = true
					}
					c.AddCheckItem(
						storage.CheckItemsSql[k].Id,
						storage.CheckItemsSql[k].Title,
						done,
					)
				}
			}
		}
	}
}

// jesus man, just make a tree ALREADY
func (s *Selector) loadCardLabels() {
	for i := 0; i < s.m.menu.Projects.Length(); i++ {
		project, err := s.m.menu.Projects.GetAt(i)
		if err != nil {
			logger.Error.Fatal(err)
		}
		p := project.(*kanban.Project)
		for j := 0; j < p.Boards.Length(); j++ {
			board, err := p.Boards.GetAt(i)
			if err != nil {
				logger.Error.Fatal(err)
			}
			b := board.(*kanban.Board)
			for k := 0; k < b.Cards.Length(); k++ {
				card, err := b.Cards.GetAt(i)
				if err != nil {
					logger.Error.Fatal(err)
				}
				c := card.(*kanban.Card)
				for l := 0; l < len(storage.CardLabelsSql); l++ {
					if c.Id == storage.CardLabelsSql[l].CardId {
						for m := 0; m < p.Labels.Length(); m++ {
							lbel, err := p.Labels.GetAt(m)
							if err != nil {
								logger.Error.Fatal(err)
							}
							label := lbel.(*kanban.Label)
							if label.Id == storage.CardLabelsSql[l].LabelId {
								c.AddLabel(label)
							}
						}
					}
				}
			}
		}
	}
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
