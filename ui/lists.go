package ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// List
// Implement list.Item interface

type Meta struct {
	initial string
	color   string
}

type Item struct {
	title       string
	description string
	meta        []Meta
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) Meta() []Meta        { return i.meta }
func (i Item) FilterValue() string { return i.title }

// Define item delegates
var (
	NoDescDelegate   list.DefaultDelegate
	DescDelegate     list.DefaultDelegate
	LabelDelegate    LabelListDelegate
	TopWhiteDelegate list.DefaultDelegate
)

func setMenuItemDelegate() {
	NoDescDelegate = list.NewDefaultDelegate()
	NoDescDelegate.ShowDescription = false
	NoDescDelegate.SetSpacing(0)
	NoDescDelegate.Styles.NormalTitle.Foreground(WHITE)
	NoDescDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
}

func setBoardItemDelegate() {
	DescDelegate = list.NewDefaultDelegate()
	DescDelegate.ShowDescription = true
	DescDelegate.SetSpacing(0)
	DescDelegate.Styles.NormalTitle.Foreground(WHITE)
	DescDelegate.Styles.NormalDesc = DescDelegate.Styles.NormalTitle.Copy()
	DescDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
	DescDelegate.Styles.SelectedDesc = DescDelegate.Styles.SelectedTitle.Copy()
}

func setLabelItemDelegate() {
	LabelDelegate = NewLabelListDelegate()
	LabelDelegate.ShowDescription = true
	LabelDelegate.Styles.NormalTitle.Foreground(WHITE)
	LabelDelegate.Styles.NormalDesc = LabelDelegate.Styles.NormalTitle.Copy()
	LabelDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
	LabelDelegate.Styles.SelectedDesc = LabelDelegate.Styles.SelectedTitle.Copy()
}

func setMoveDelegate() {
	TopWhiteDelegate = list.NewDefaultDelegate()
	TopWhiteDelegate.Styles.NormalTitle.Foreground(WHITE)
	TopWhiteDelegate.Styles.SelectedTitle.
		Foreground(SelectedListItemColor).
		Padding(0, 0, 0, 2).
		BorderLeft(false).
		BorderTop(true).
		BorderForeground(WHITE)
	TopWhiteDelegate.Styles.SelectedDesc = TopWhiteDelegate.Styles.SelectedTitle.Copy().
		BorderTop(false)
}

// Menu
func (m *Menu) setupList() {
	var (
		node      *dll.Node
		menuItems []list.Item
	)
	l := list.New([]list.Item{}, NoDescDelegate, ws.width/3, ws.height-6)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	for i := 0; i < m.menu.Projects.Length(); i++ {
		node, _ = m.menu.Projects.WalkTo(i)
		project := node.Val().(*kanban.Project)
		item := Item{
			title: project.Title,
		}
		menuItems = append(menuItems, item)
	}
	l.SetItems(menuItems)
	m.list = l
}

// Project
func (p *Project) setupBoards() {
	var (
		err         error
		boards      []list.Model
		node        *dll.Node
		cardNode    *dll.Node
		checkitem   *kanban.CheckItem
		label       *kanban.Label
		items       []list.Item
		placeholder interface{}
		metaSlice   []Meta
	)
	node, _ = p.project.Boards.HeadNode()
	if node == nil {
		return
	}
	board := node.Val().(*kanban.Board)
	for i := 0; i < p.project.Boards.Length(); i++ {
		b := list.New([]list.Item{}, NewCardDelegate(), ws.width/3, ws.height-9)
		b.SetShowHelp(false)
		b.Title = board.Title
		b.InfiniteScrolling = true
		for j := 0; j < board.Cards.Length(); j++ {
			cardNode, err = board.Cards.WalkTo(j)
			if err != nil {
				log.Println(err)
				err = nil
			}
			c := cardNode.Val().(*kanban.Card)
			checkTotal := c.CheckList.Length()
			checkDone := 0
			for k := 0; k < checkTotal; k++ {
				placeholder, err = c.CheckList.GetAt(k)
				if err != nil {
					log.Println(err)
					err = nil
				}
				checkitem = placeholder.(*kanban.CheckItem)
				if checkitem.Check {
					checkDone++
				}
			}

			labelLen := c.CardLabels.Length()

			for k := 0; k < labelLen; k++ {
				placeholder, err = c.CardLabels.GetAt(k)
				if err != nil {
					log.Println(err)
					err = nil
				}
				label = placeholder.(*kanban.Label)
				meta := Meta{
					initial: string(label.Title[0]),
					color:   label.Color,
				}
				metaSlice = append(metaSlice, meta)
			}

			item := Item{
				title: c.Title,
				description: fmt.Sprintf("[✓]%s/%s %sL",
					strconv.Itoa(checkDone),
					strconv.Itoa(checkTotal),
					strconv.Itoa(labelLen),
				),
				meta: metaSlice,
			}
			items = append(items, item)
			b.SetItems(items)
		}
		items = []list.Item{}
		boards = append(boards, b)
		node, _ = node.Next()
		if node != nil {
			board = node.Val().(*kanban.Board)
		}
	}
	p.boards = boards
}

// Label
func (l *Label) setupList() {
	var (
		node       *dll.Node
		labelItems []list.Item
	)
	lst := list.New([]list.Item{}, NewLabelListDelegate(), ws.width/3, ws.height-9)
	lst.SetShowHelp(false)
	lst.Title = "Labels"
	lst.InfiniteScrolling = true
	for i := 0; i < l.project.Labels.Length(); i++ {
		node, _ = l.project.Labels.WalkTo(i)
		label := node.Val().(*kanban.Label)
		item := Item{
			title:       label.Title,
			description: label.Color,
		}
		labelItems = append(labelItems, item)
	}
	lst.SetItems(labelItems)
	l.list = lst
}

// Card
func (c *Card) setupLists() {
	var (
		node           *dll.Node
		checklistItems []list.Item
		labelItems     []list.Item
	)

	cl := list.New([]list.Item{}, NewCheckListDelegate(), ws.width/2, ws.height/3+1)
	cl.SetShowHelp(false)
	cl.Title = "Checklist"
	cl.InfiniteScrolling = true
	for i := 0; i < c.card.CheckList.Length(); i++ {
		node, _ = c.card.CheckList.WalkTo(i)
		checkItem := node.Val().(*kanban.CheckItem)
		item := Item{
			title: checkItem.Title,
		}
		if checkItem.Check {
			item.description = "1"
		} else {
			item.description = "0"
		}
		checklistItems = append(checklistItems, item)
	}
	cl.SetItems(checklistItems)
	c.checklist = cl

	ll := list.New([]list.Item{}, NewLabelListDelegate(), ws.width/2, ws.height/3+1)
	ll.SetShowHelp(false)
	ll.Title = "Card Labels"
	ll.InfiniteScrolling = true
	for i := 0; i < c.card.CardLabels.Length(); i++ {
		node, _ = c.card.CardLabels.WalkTo(i)
		label := node.Val().(*kanban.Label)
		item := Item{
			title: label.Title,
		}
		labelItems = append(labelItems, item)
	}
	ll.SetItems(labelItems)
	c.labels = ll
}
