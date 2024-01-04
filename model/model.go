package model

import (
	"image/color"

	"github.com/Anacardo89/ds/lists/dll"
)

type Label struct {
	Title string
	Color color.Color
}

type CheckItem struct {
	Title string
	Check bool
}

type Card struct {
	Title       string
	Description string
	CheckList   dll.DLL
	Labels      dll.DLL
}

type List struct {
	Title string
	Cards dll.DLL
}

type Board struct {
	Title string
	Lists dll.DLL
}

func (b *Board) StartBoard(title string) {
	b.Title = title
	b.Lists = dll.New()
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

func (c *CheckItem) RenameCheckItem(title string) {
	c.Title = title
}

func (c *CheckItem) CheckCheckItem() {
	c.Check = !c.Check
}

func (c *Card) AddLabel(title string, color color.Color) {
	label := Label{
		Title: title,
		Color: color,
	}
	c.Labels.Append(label)
}

func (c *Card) RemoveLabel(label dll.DLL) error {
	_, err := c.CheckList.Remove(label)
	return err
}
