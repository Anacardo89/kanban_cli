/*
Menu
	|-Project
		|-Label
		|-Board
			|-Card
				|-CheckList
				|-CardLabels
*/

package kanban

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanboards/logger"
)

type Menu struct {
	Projects dll.DLL
}

type Project struct {
	Id     int64
	Title  string
	Boards dll.DLL
	Labels dll.DLL
}

type Board struct {
	Id    int64
	Pos   int
	Title string
	Cards dll.DLL
}

type Label struct {
	Id    int64
	Title string
	Color string
}

type Card struct {
	Id          int64
	Title       string
	Description string
	CheckList   dll.DLL
	CardLabels  dll.DLL
}

type CheckItem struct {
	Id    int64
	Title string
	Check bool
}

// Menu
func StartMenu() *Menu {
	return &Menu{
		Projects: dll.New(),
	}
}

func (m *Menu) AddProject(id int64, title string) {
	project := &Project{
		Id:     id,
		Title:  title,
		Boards: dll.New(),
		Labels: dll.New(),
	}
	m.Projects.Append(project)
}

func (m *Menu) GetProjectById(id int64) *Project {
	for i := 0; i < m.Projects.Length(); i++ {
		curr, err := m.Projects.GetAt(i)
		if err != nil {
			logger.Error.Println(err)
		}
		c := curr.(*Project)
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (m *Menu) RemoveProject(project *Project) error {
	err := m.Projects.Remove(project)
	return err
}

// Project
func (p *Project) RenameProject(title string) {
	p.Title = title
}

func (p *Project) AddBoard(id int64, title string) {
	board := &Board{
		Id:    id,
		Pos:   p.Boards.Length(),
		Title: title,
		Cards: dll.New(),
	}
	p.Boards.Append(board)
}

func (p *Project) GetBoardById(id int64) *Board {
	for i := 0; i < p.Boards.Length(); i++ {
		curr, err := p.Boards.GetAt(i)
		if err != nil {
			logger.Error.Println(err)
		}
		c := curr.(*Board)
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (p *Project) RemoveBoard(board *Board) error {
	err := p.Boards.Remove(board)
	return err
}

func (p *Project) AddLabel(id int64, title string, color string) {
	label := &Label{
		Id:    id,
		Title: title,
		Color: color,
	}
	p.Labels.Append(label)
}

func (p *Project) GetLabelById(id int64) *Label {
	for i := 0; i < p.Labels.Length(); i++ {
		curr, err := p.Labels.GetAt(i)
		if err != nil {
			logger.Error.Println(err)
		}
		c := curr.(*Label)
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (p *Project) RemoveLabel(label *Label) error {
	err := p.Labels.Remove(label)
	return err
}

// Label
func (l *Label) RenameLabel(title string) {
	l.Title = title
}

func (l *Label) ChangeColor(color string) {
	l.Color = color
}

// Board
func (b *Board) RenameBoard(title string) {
	b.Title = title
}

func (b *Board) AddCard(id int64, title string, desc string) {
	card := &Card{
		Id:          id,
		Title:       title,
		Description: desc,
		CheckList:   dll.New(),
		CardLabels:  dll.New(),
	}
	b.Cards.Append(card)
}

func (b *Board) GetCardById(id int64) *Card {
	for i := 0; i < b.Cards.Length(); i++ {
		curr, err := b.Cards.GetAt(i)
		if err != nil {
			logger.Error.Println(err)
		}
		c := curr.(*Card)
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (b *Board) RemoveCard(card *Card) error {
	err := b.Cards.Remove(card)
	return err
}

// Card
func (c *Card) RenameCard(title string) {
	c.Title = title
}

func (c *Card) AddDescription(description string) {
	c.Description = description
}

func (c *Card) AddCheckItem(id int64, title string, done bool) {
	checkItem := &CheckItem{
		Id:    id,
		Title: title,
		Check: done,
	}
	c.CheckList.Append(checkItem)
}

func (c *Card) GetCheckItemById(id int64) *CheckItem {
	for i := 0; i < c.CheckList.Length(); i++ {
		curr, err := c.CheckList.GetAt(i)
		if err != nil {
			logger.Error.Println(err)
		}
		c := curr.(*CheckItem)
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (c *Card) RemoveCheckItem(checkItem *CheckItem) error {
	err := c.CheckList.Remove(checkItem)
	return err
}

func (c *Card) AddLabel(label *Label) {
	c.CardLabels.Append(label)
}

func (c *Card) RemoveLabel(label *Label) error {
	err := c.CardLabels.Remove(label)
	return err
}

// CheckItem
func (c *CheckItem) RenameCheckItem(title string) {
	c.Title = title
}

func (c *CheckItem) CheckCheckItem() {
	c.Check = !c.Check
}
