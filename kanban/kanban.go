/*
Menu
  |_Project
    |_Label
	|_List
	  |_Card
		|_CheckList
		|_CardLabels
*/

package kanban

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	Projects dll.DLL
}

type Project struct {
	Title  string
	Lists  dll.DLL
	Labels dll.DLL
}

type List struct {
	Title string
	Cards dll.DLL
}

type Label struct {
	Title string
	Color lipgloss.Color
}

type Card struct {
	Title       string
	Description string
	CheckList   dll.DLL
	CardLabels  dll.DLL
}

type CheckItem struct {
	Title string
	Check bool
}

// Menu
func StartMenu() *Menu {
	return &Menu{
		Projects: dll.New(),
	}
}

func (m *Menu) AddProject(title string) {
	project := &Project{
		Title: title,
		Lists: dll.New(),
	}
	m.Projects.Append(project)
}

func (m *Menu) RemoveProject(project dll.DLL) error {
	_, err := m.Projects.Remove(project)
	return err
}

// Project
func (p *Project) RenameProject(title string) {
	p.Title = title
}

func (p *Project) AddList(title string) {
	list := &List{
		Title: title,
		Cards: dll.New(),
	}
	p.Lists.Append(list)
}

func (p *Project) RemoveList(list dll.DLL) error {
	_, err := p.Lists.Remove(list)
	return err
}

func (p *Project) AddLabel(title string, color lipgloss.Color) {
	label := &Label{
		Title: title,
		Color: color,
	}
	p.Labels.Append(label)
}

func (p *Project) RemoveLabel(label dll.DLL) error {
	_, err := p.Labels.Remove(label)
	return err
}

// Label
func (l *Label) RenameLabel(title string) {
	l.Title = title
}

func (l *Label) ChangeColor(color lipgloss.Color) {
	l.Color = color
}

// List
func (l *List) RenameList(title string) {
	l.Title = title
}

func (l *List) AddCard(title string) {
	card := &Card{
		Title:      title,
		CheckList:  dll.New(),
		CardLabels: dll.New(),
	}
	l.Cards.Append(card)
}

func (l *List) RemoveCard(card dll.DLL) error {
	_, err := l.Cards.Remove(card)
	return err
}

// Card
func (c *Card) RenameCard(title string) {
	c.Title = title
}

func (c *Card) AddDescription(description string) {
	c.Description = description
}

func (c *Card) AddCheckItem(title string) {
	checkItem := &CheckItem{
		Title: title,
		Check: false,
	}
	c.CheckList.Append(checkItem)
}

func (c *Card) RemoveCheckItem(checkItem dll.DLL) error {
	_, err := c.CheckList.Remove(checkItem)
	return err
}

func (c *Card) AddLabel(label Label) {
	c.CardLabels.Append(label)
}

func (c *Card) RemoveLabel(label dll.DLL) error {
	_, err := c.CardLabels.Remove(label)
	return err
}

// CheckItem
func (c *CheckItem) RenameCheckItem(title string) {
	c.Title = title
}

func (c *CheckItem) CheckCheckItem() {
	c.Check = !c.Check
}
