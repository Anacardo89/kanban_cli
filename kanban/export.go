package kanban

import (
	"github.com/Anacardo89/kanboards/storage"
)

func (m *Menu) Export() {
	projects := m.projectsToStorage()
	sm := storage.Menu{
		Projects: projects,
	}
	storage.ToFile(sm.ToYAML())
}

func (m *Menu) projectsToStorage() []storage.Project {
	projects := []storage.Project{}
	for i := 0; i < m.Projects.Length(); i++ {
		projectNode, _ := m.Projects.GetNodeAt(i)
		projectVal := projectNode.Val().(*Project)
		project := storage.Project{
			Id:     projectVal.Id,
			Title:  projectVal.Title,
			Labels: projectVal.labelsToStorage(),
			Boards: projectVal.boardsToStorage(),
		}
		projects = append(projects, project)
	}
	return projects
}

func (p *Project) boardsToStorage() []storage.Board {
	boards := []storage.Board{}
	for i := 0; i < p.Boards.Length(); i++ {
		boardNode, _ := p.Boards.GetNodeAt(i)
		boardVal := boardNode.Val().(*Board)
		board := storage.Board{
			Id:    boardVal.Id,
			Pos:   boardVal.Pos,
			Title: boardVal.Title,
			Cards: boardVal.cardsToStorage(),
		}
		boards = append(boards, board)
	}
	return boards
}

func (p *Project) labelsToStorage() []storage.Label {
	labels := []storage.Label{}
	for i := 0; i < p.Labels.Length(); i++ {
		labelNode, _ := p.Labels.GetNodeAt(i)
		labelVal := labelNode.Val().(*Label)
		label := storage.Label{
			Id:    labelVal.Id,
			Title: labelVal.Title,
			Color: labelVal.Color,
		}
		labels = append(labels, label)
	}
	return labels
}

func (b *Board) cardsToStorage() []storage.Card {
	cards := []storage.Card{}
	for i := 0; i < b.Cards.Length(); i++ {
		cardNode, _ := b.Cards.GetNodeAt(i)
		cardVal := cardNode.Val().(*Card)
		card := storage.Card{
			Id:          cardVal.Id,
			Title:       cardVal.Title,
			Description: cardVal.Description,
			CheckList:   cardVal.checkListToStorage(),
			CardLabels:  cardVal.cardLabelsToStorage(),
		}
		cards = append(cards, card)
	}
	return cards
}

func (c *Card) checkListToStorage() []storage.CheckItem {
	checkList := []storage.CheckItem{}
	for i := 0; i < c.CheckList.Length(); i++ {
		checkNode, _ := c.CheckList.GetNodeAt(i)
		checkVal := checkNode.Val().(*CheckItem)
		checkItem := storage.CheckItem{
			Id:    checkVal.Id,
			Title: checkVal.Title,
			Check: checkVal.Check,
		}
		checkList = append(checkList, checkItem)
	}
	return checkList
}

func (c *Card) cardLabelsToStorage() []storage.Label {
	cardLabels := []storage.Label{}
	for i := 0; i < c.CardLabels.Length(); i++ {
		cardLabelNode, _ := c.CardLabels.GetNodeAt(i)
		cardLabelVal := cardLabelNode.Val().(*Label)
		cardLabel := storage.Label{
			Id:    cardLabelVal.Id,
			Title: cardLabelVal.Title,
			Color: cardLabelVal.Color,
		}
		cardLabels = append(cardLabels, cardLabel)
	}
	return cardLabels
}
