package ui

import (
	"github.com/Anacardo89/ds/lists/queue"
	"github.com/Anacardo89/kanboards/kanban"
	"github.com/Anacardo89/kanboards/logger"
	"github.com/Anacardo89/kanboards/storage"
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
		boards := storage.GetBoardsWithParentOrdered(project.Id)
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
	len := q.Length()
	for i := 0; i < len; i++ {
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
			p.AddLabel(v.Id, v.Title, v.Color)
		case storage.BoardSql:
			c = nil
			p.AddBoard(v.Id, v.Title)
			b = p.GetBoardById(v.Id)
		case storage.CardSql:
			desc := v.Desc.String
			b.AddCard(v.Id, v.Title, desc)
			c = b.GetCardById(v.Id)
		case storage.CheckItemSql:
			done := false
			if v.Done == 1 {
				done = true
			}
			c.AddCheckItem(v.Id, v.Title, done)
		case storage.CardLabelSql:
			l := p.GetLabelById(v.LabelId)
			c.AddLabel(l)
		}
	}
}
