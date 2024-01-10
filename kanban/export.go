package kanban

import (
	"github.com/Anacardo89/kanban_cli/storage"
)

func (m *Menu) ToStorage() *storage.Menu {
	projects := m.projectsToStorage()
	return &storage.Menu{
		Projects: projects,
	}
}

func (m *Menu) projectsToStorage() []storage.Project {
	projects := []storage.Project{}
	for i := 0; i < m.Projects.GetLength(); i++ {
		projectNode, _ := m.Projects.WalkTo(i)
		projectVal := projectNode.GetVal().(*Project)
		project := storage.Project{
			Title:  projectVal.Title,
			Lists:  projectVal.listsToStorage(),
			Labels: projectVal.labelsToStorage(),
		}
		projects = append(projects, project)
	}
	return projects
}

func (p *Project) listsToStorage() []storage.List {
	lists := []storage.List{}
	for i := 0; i < p.Lists.GetLength(); i++ {
		listNode, _ := p.Lists.WalkTo(i)
		listVal := listNode.GetVal().(*List)
		list := storage.List{
			Title: listVal.Title,
			Cards: listVal.cardsToStorage(),
		}
		lists = append(lists, list)
	}
	return lists
}

func (p *Project) labelsToStorage() []storage.Label {
	labels := []storage.Label{}
	for i := 0; i < p.Labels.GetLength(); i++ {
		labelNode, _ := p.Labels.WalkTo(i)
		labelVal := labelNode.GetVal().(*Label)
		label := storage.Label{
			Title: labelVal.Title,
			Color: labelVal.Color,
		}
		labels = append(labels, label)
	}
	return labels
}

func (l *List) cardsToStorage() []storage.Card {
	cards := []storage.Card{}
	for i := 0; i < l.Cards.GetLength(); i++ {
		cardNode, _ := l.Cards.WalkTo(i)
		cardVal := cardNode.GetVal().(*Card)
		card := storage.Card{
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
	for i := 0; i < c.CheckList.GetLength(); i++ {
		checkNode, _ := c.CheckList.WalkTo(i)
		checkVal := checkNode.GetVal().(*CheckItem)
		checkItem := storage.CheckItem{
			Title: checkVal.Title,
			Check: checkVal.Check,
		}
		checkList = append(checkList, checkItem)
	}
	return checkList
}

func (c *Card) cardLabelsToStorage() []storage.Label {
	cardLabels := []storage.Label{}
	for i := 0; i < c.CardLabels.GetLength(); i++ {
		cardLabelNode, _ := c.CardLabels.WalkTo(i)
		cardLabelVal := cardLabelNode.GetVal().(*Label)
		cardLabel := storage.Label{
			Title: cardLabelVal.Title,
			Color: cardLabelVal.Color,
		}
		cardLabels = append(cardLabels, cardLabel)
	}
	return cardLabels
}
