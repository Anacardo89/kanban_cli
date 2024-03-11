package ui

import (
	"github.com/Anacardo89/ds/lists/queue"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/storage"
)

func (s *Selector) loadKanban() {
	q := DBToQueue()
	s.QueueToKanban(q)
}

func DBToQueue() queue.Queue {
	q := queue.New()
	projects := storage.GetAllProjects()
	for _, project := range projects {
		q.Enqueue(project)
		labels := storage.GetLabelsWithParent(project.Id)
		for _, label := range labels {
			q.Enqueue(label)
		}
		boards := storage.GetBoardsWithParent(project.Id)
		for _, board := range boards {
			q.Enqueue(board)
			cards := storage.GetCardsWithParent(board.Id)
			for _, card := range cards {
				q.Enqueue(card)
				checklist := storage.GetCheckItemsWithParent(card.Id)
				for _, ci := range checklist {
					q.Enqueue(ci)
				}
				cardLabels := storage.GetLabelsInCard(card.Id)
				for _, cl := range cardLabels {
					q.Enqueue(cl)
				}
			}
		}
	}
	return q
}

func (s *Selector) QueueToKanban(q queue.Queue) {
	var (
		p *kanban.Project
		b *kanban.Board
		c *kanban.Card
	)
	for i := 0; i < q.Length(); i++ {
		val, err := q.Dequeue()
		if err != nil {
			logger.Error.Println("Could not read from load queue:", err)
		}
		switch v := val.(type) {
		case storage.ProjectSql:
			b, c = nil, nil
			s.m.menu.AddProject(v.Id, v.Title)
			p = s.m.menu.GetProjectById(v.Id)
		case storage.LabelSql:
			if c == nil {
				p.AddLabel(v.Id, v.Title, v.Color)
			} else {
				l := p.GetLabelById(v.Id)
				c.AddLabel(l)
			}
		case *kanban.Board:
			c = nil
			p.AddBoard(v.Id, v.Title)
			b = p.GetBoardById(v.Id)
		case *kanban.Card:
			b.AddCard(v.Id, v.Title, v.Description)
			c = b.GetCardById(v.Id)
		case *kanban.CheckItem:
			c.AddCheckItem(v.Id, v.Title, v.Check)

		}
	}
}
