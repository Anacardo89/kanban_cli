/*
Menu
  |_Board
    |_Label
	|_List
	  |_Card
		|_CheckList

*/

package kanban

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	Boards dll.DLL
}

func StartMenu() *Menu {
	return &Menu{
		Boards: dll.New(),
	}
}

func (m *Menu) AddBoard(title string) {
	board := Board{
		Title: title,
		Lists: dll.New(),
	}
	m.Boards.Append(board)
}

func (m *Menu) RemoveBoard(board dll.DLL) error {
	_, err := m.Boards.Remove(board)
	return err
}

type Board struct {
	Title  string
	Lists  dll.DLL
	Labels dll.DLL
}

func (b *Board) RenameBoard(title string) {
	b.Title = title
}

func (b *Board) AddList(title string) {
	list := List{
		Title: title,
		Cards: dll.New(),
	}
	b.Lists.Append(list)
}

func (b *Board) RemoveList(list dll.DLL) error {
	_, err := b.Lists.Remove(list)
	return err
}

func (b *Board) AddLabel(title string, color lipgloss.Color) {
	label := Label{
		Title: title,
		Color: color,
	}
	b.Labels.Append(label)
}

func (b *Board) RemoveLabel(label dll.DLL) error {
	_, err := b.Labels.Remove(label)
	return err
}

type Label struct {
	Title string
	Color lipgloss.Color
}

func (l *Label) RenameLabel(title string) {
	l.Title = title
}

func (l *Label) ChangeColor(color lipgloss.Color) {
	l.Color = color
}

type List struct {
	Title string
	Cards dll.DLL
}

func (l *List) RenameList(title string) {
	l.Title = title
}

func (l *List) AddCard(title string) {
	card := Card{
		Title:     title,
		CheckList: dll.New(),
		Labels:    dll.New(),
	}
	l.Cards.Append(card)
}

func (l *List) RemoveCard(card dll.DLL) error {
	_, err := l.Cards.Remove(card)
	return err
}

type Card struct {
	Title       string
	Description string
	CheckList   dll.DLL
	Labels      dll.DLL
}

func (c *Card) RenameCard(title string) {
	c.Title = title
}

func (c *Card) AddDescription(description string) {
	c.Description = description
}

func (c *Card) AddCheckItem(title string) {
	checkItem := CheckItem{
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
	c.Labels.Append(label)
}

func (c *Card) RemoveLabel(label dll.DLL) error {
	_, err := c.Labels.Remove(label)
	return err
}

type CheckItem struct {
	Title string
	Check bool
}

func (c *CheckItem) RenameCheckItem(title string) {
	c.Title = title
}

func (c *CheckItem) CheckCheckItem() {
	c.Check = !c.Check
}
