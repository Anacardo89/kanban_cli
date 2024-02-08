package ui

import (
	"log"

	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Input Field
type InputField struct {
	field textinput.Model
	data  string
}

// Menu
func (m *Menu) setInput() {
	m.Input.field.Prompt = ": "
	m.Input.field.CharLimit = 120
	m.Input.field.Placeholder = "Project Title"
}

func (m *Menu) handleInput(key string) {
	switch key {
	case "esc":
		m.Input.field.SetValue("")
		m.Input.data = ""
		m.Input.field.Blur()
		return
	case "enter":
		m.Input.data = m.Input.field.Value()
		if m.Input.data == "" {
			return
		}
		menuItem := Item{
			title: m.Input.data,
		}
		menuItems = append(menuItems, menuItem)
		m.list.SetItems(menuItems)
		m.menu.AddProject(m.Input.data)
		m.Input.data = ""
		m.Input.field.SetValue("")
		m.Input.field.Blur()
		m.cursor = 0
		return
	}
}

// Project
func (p *Project) setInput() {
	p.Input.field.Prompt = ": "
	p.Input.field.CharLimit = 120
	switch p.inputFlag {
	case new:
		p.Input.field.Placeholder = "Board Title"
	case add:
		p.Input.field.Placeholder = "Card Title"
	}
}

func (p *Project) handleInput(key string) {
	switch key {
	case "esc":
		p.Input.field.SetValue("")
		p.Input.data = ""
		p.Input.field.Blur()
		return
	case "enter":
		p.Input.data = p.Input.field.Value()
		if p.Input.data == "" {
			return
		}
		switch p.inputFlag {
		case new:
			p.project.AddBoard(p.Input.data)
			p.setupBoards()
			p.hcursor = 0
		case add:
			node, err := p.project.Boards.WalkTo(p.hcursor)
			if err != nil {
				log.Println(err)
				err = nil
				return
			}
			board := node.Val().(*kanban.Board)
			board.AddCard(p.Input.data)
			boardItems = p.boards[p.hcursor].Items()
			boardItem := Item{
				title: p.Input.data,
			}
			boardItems = append(boardItems, boardItem)
			p.boards[p.hcursor].SetItems(boardItems)
		}
		p.Input.data = ""
		p.Input.field.SetValue("")
		p.Input.field.Blur()
		return
	}
}

// Label
func (l *Label) setInput() {
	l.Input.field.Prompt = ": "
	l.Input.field.CharLimit = 120
	l.Input.field.Placeholder = "Label Ttile"
}

func (l *Label) handleInput(key string) {
	switch key {
	case "esc":
		l.Input.field.SetValue("")
		l.Input.data = ""
		l.Input.field.Blur()
		return
	case "enter":
		if labelName == "" {
			l.Input.data = l.Input.field.Value()
			if l.Input.data == "" {
				return
			}
			labelName = l.Input.data
			l.Input.data = ""
			l.Input.field.SetValue("")
			l.Input.field.Placeholder = "Label Hex Color"
			return
		}
		l.Input.data = l.Input.field.Value()
		if l.Input.data[0] != '#' {
			l.Input.data = string('#') + l.Input.data
		}
		if len(l.Input.data) != 7 {
			return
		}
		l.project.AddLabel(labelName, lipgloss.Color(l.Input.data))
		labelItem := Item{
			title: l.Input.data,
		}
		labelItems = append(labelItems, labelItem)
		l.list.SetItems(labelItems)
		l.setupList()
		l.Input.data = ""
		l.Input.field.SetValue("")
		l.Input.field.Blur()
	}
}

// List
// Implement list.Item interface
type Item struct {
	title       string
	description string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) FilterValue() string { return i.title }

// Define item delegates
var (
	NoDescDelegate   list.DefaultDelegate
	DescDelegate     list.DefaultDelegate
	LabelDelegate    list.DefaultDelegate
	TopWhiteDelegate list.DefaultDelegate
)

func setMenuItemDelegate() {
	NoDescDelegate = list.NewDefaultDelegate()
	NoDescDelegate.ShowDescription = false
	NoDescDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	NoDescDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
}

func setBoardItemDelegate() {
	DescDelegate = list.NewDefaultDelegate()
	DescDelegate.ShowDescription = true
	DescDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	DescDelegate.Styles.NormalDesc = DescDelegate.Styles.NormalTitle.Copy()
	DescDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
	DescDelegate.Styles.SelectedDesc = DescDelegate.Styles.SelectedTitle.Copy()
}

func setLabelItemDelegate() {
	LabelDelegate = list.NewDefaultDelegate()
	LabelDelegate.ShowDescription = true
	LabelDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	LabelDelegate.Styles.NormalDesc = LabelDelegate.Styles.NormalTitle.Copy()
	LabelDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).
		Border(lipgloss.HiddenBorder(), false, false, false, true)
	LabelDelegate.Styles.SelectedDesc = LabelDelegate.Styles.SelectedTitle.Copy()
}

func setMoveDelegate() {
	TopWhiteDelegate = list.NewDefaultDelegate()
	TopWhiteDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	TopWhiteDelegate.Styles.SelectedTitle.
		Foreground(SelectedListItemColor).
		BorderLeft(false).
		BorderTop(true).
		BorderForeground(WHITE)
	TopWhiteDelegate.Styles.SelectedDesc = TopWhiteDelegate.Styles.SelectedTitle.Copy().
		BorderTop(false)
}

// Menu
var (
	menuItems []list.Item
)

func (m *Menu) setupList() {
	var (
		node  *dll.Node
		items []list.Item
	)
	l := list.New([]list.Item{}, NoDescDelegate, ws.width/3, ws.height-9)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	for i := 0; i < m.menu.Projects.Length(); i++ {
		node, _ = m.menu.Projects.WalkTo(i)
		project := node.Val().(*kanban.Project)
		item := Item{
			title: project.Title,
		}
		items = append(items, item)
	}
	l.SetItems(items)
	m.list = l
}

// Project
var boardItems []list.Item

func (p *Project) setupBoards() {
	var (
		err      error
		boards   []list.Model
		node     *dll.Node
		cardNode *dll.Node
		items    []list.Item
	)
	node, _ = p.project.Boards.HeadNode()
	if node == nil {
		return
	}
	board := node.Val().(*kanban.Board)
	for i := 0; i < p.project.Boards.Length(); i++ {
		b := list.New([]list.Item{}, DescDelegate, ws.width/3, ws.height-9)
		b.SetShowHelp(false)
		b.Title = board.Title
		b.InfiniteScrolling = true
		for j := 0; j < board.Cards.Length(); j++ {
			cardNode, err = board.Cards.WalkTo(j)
			if err != nil {
				log.Println(err)
			}
			c := cardNode.Val().(*kanban.Card)
			item := Item{
				title: c.Title,
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

// List
var labelItems []list.Item

func (l *Label) setupList() {
	var (
		node  *dll.Node
		items []list.Item
	)
	lst := list.New([]list.Item{}, LabelDelegate, ws.width/3, ws.height-9)
	lst.SetShowHelp(false)
	lst.Title = "Labels"
	lst.InfiniteScrolling = true
	for i := 0; i < l.project.Labels.Length(); i++ {
		node, _ = l.project.Labels.WalkTo(i)
		label := node.Val().(*kanban.Label)
		item := Item{
			title: label.Title,
		}
		LabelDelegate.Styles.NormalDesc.Background(label.Color)
		LabelDelegate.Styles.SelectedDesc.Background(label.Color)
		labelItems = append(labelItems, item)
	}
	lst.SetItems(items)
	l.list = lst
	l.cursor = l.list.Cursor()
}
