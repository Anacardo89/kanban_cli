package kanban

import (
	"github.com/Anacardo89/kanboards/storage"
)

func (m *Menu) Import() {
	data := storage.FromFile()
	menu := storage.FromYAML(data)
	m.projectsFromStorage(menu.Projects)
}

func (m *Menu) projectsFromStorage(sp []storage.Project) {
	for i := 0; i < len(sp); i++ {
		m.AddProject(sp[i].Id, sp[i].Title)
		p := m.GetProjectById(sp[i].Id)
		p.labelsFromStorage(sp[i].Labels)
		p.boardsFromStorage(sp[i].Boards)
	}
}

func (p *Project) labelsFromStorage(sl []storage.Label) {
	for i := 0; i < len(sl); i++ {
		p.AddLabel(sl[i].Id, sl[i].Title, sl[i].Color)
	}
}

func (p *Project) boardsFromStorage(sb []storage.Board) {
	for i := 0; i < len(sb); i++ {
		p.AddBoard(sb[i].Id, sb[i].Title)
		b := p.GetBoardById(sb[i].Id)
		b.cardsFromStorage(sb[i].Cards, p)
	}
}

func (b *Board) cardsFromStorage(sc []storage.Card, p *Project) {
	for i := 0; i < len(sc); i++ {
		b.AddCard(sc[i].Id, sc[i].Title, sc[i].Description)
		c := b.GetCardById(sc[i].Id)
		c.checkListFromStorage(sc[i].CheckList)
		c.cardLabelsFromStorage(sc[i].CardLabels, p)
	}
}

func (c *Card) checkListFromStorage(sci []storage.CheckItem) {
	for i := 0; i < len(sci); i++ {
		c.AddCheckItem(sci[i].Id, sci[i].Title, sci[i].Check)
	}
}

func (c *Card) cardLabelsFromStorage(scl []storage.Label, p *Project) {
	for i := 0; i < len(scl); i++ {
		l := p.GetLabelById(scl[i].Id)
		c.AddLabel(l)
	}
}
