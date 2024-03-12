package kanban

import (
	"github.com/Anacardo89/kanboards/logger"
	"github.com/Anacardo89/kanboards/storage"
)

func (m *Menu) Import() {
	data := storage.FromFile()
	menu := storage.FromYAML(data)
	m.projectsFromStorage(menu.Projects)
}

func (m *Menu) projectsFromStorage(sp []storage.Project) {
	for i := 0; i < len(sp); i++ {
		res := storage.CreateProject(sp[i].Title)
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error.Println(err)
		}
		m.AddProject(id, sp[i].Title)
		p := m.GetProjectById(sp[i].Id)
		p.labelsFromStorage(sp[i].Labels)
		p.boardsFromStorage(sp[i].Boards)
	}
}

func (p *Project) labelsFromStorage(sl []storage.Label) {
	for i := 0; i < len(sl); i++ {
		res := storage.CreateLabel(sl[i].Title, sl[i].Color, p.Id)
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error.Println(err)
		}
		p.AddLabel(id, sl[i].Title, sl[i].Color)
	}
}

func (p *Project) boardsFromStorage(sb []storage.Board) {
	for i := 0; i < len(sb); i++ {
		res := storage.CreateBoard(sb[i].Title, p.Id)
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error.Println(err)
		}
		p.AddBoard(id, sb[i].Title)
		b := p.GetBoardById(sb[i].Id)
		b.cardsFromStorage(sb[i].Cards, p)
	}
}

func (b *Board) cardsFromStorage(sc []storage.Card, p *Project) {
	for i := 0; i < len(sc); i++ {
		res := storage.CreateCard(sc[i].Title, b.Id)
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error.Println(err)
		}
		storage.UpdateCardDesc(id, sc[i].Description)
		b.AddCard(id, sc[i].Title, sc[i].Description)
		c := b.GetCardById(sc[i].Id)
		c.checkListFromStorage(sc[i].CheckList)
		c.cardLabelsFromStorage(sc[i].CardLabels, p)
	}
}

func (c *Card) checkListFromStorage(sci []storage.CheckItem) {
	for i := 0; i < len(sci); i++ {
		done := 0
		if sci[i].Check {
			done = 1
		}
		res := storage.CreateCheckItem(sci[i].Title, done, c.Id)
		id, err := res.LastInsertId()
		if err != nil {
			logger.Error.Println(err)
		}
		c.AddCheckItem(id, sci[i].Title, sci[i].Check)
	}
}

func (c *Card) cardLabelsFromStorage(scl []storage.Label, p *Project) {
	for i := 0; i < len(scl); i++ {
		storage.CreateCardLabel(c.Id, scl[i].Id)
		l := p.GetLabelById(scl[i].Id)
		c.AddLabel(l)
	}
}
