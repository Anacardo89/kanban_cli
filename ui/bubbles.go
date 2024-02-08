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
	var node *dll.Node
	var err error
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
			node, err = p.project.Boards.HeadNode()
			if err != nil {
				log.Println(err)
			}
			p.sb = node
		case add:
			board := p.sb.Val().(*kanban.Board)
			board.AddCard(p.Input.data)
			boardItems = p.boards[p.hcursor].Items()
			boardItem := Item{
				title: p.Input.data,
			}
			boardItems = append(boardItems, boardItem)
			p.boards[p.hcursor].SetItems(boardItems)
			node, _ = board.Cards.HeadNode()
			if err != nil {
				log.Println(err)
			}
			p.sc = node
		}
		p.Input.data = ""
		p.Input.field.SetValue("")
		p.Input.field.Blur()
		p.vcursor = 0
		return
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

// Menu
var (
	menuItems    []list.Item
	menuDelegate list.DefaultDelegate
)

func setMenuItemDelegate() {

	menuDelegate = list.NewDefaultDelegate()
	menuDelegate.ShowDescription = false
	menuDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	menuDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).Border(lipgloss.HiddenBorder(), false, false, false, true)
}

func (m *Menu) setupMenuList() {
	var (
		node  *dll.Node
		items []list.Item
	)
	l := list.New([]list.Item{}, menuDelegate, ws.width/3, ws.height-9)
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
var (
	boardItems    []list.Item
	boardDelegate list.DefaultDelegate
	moveDelegate  list.DefaultDelegate
)

func setBoardItemDelegate() {
	boardDelegate = list.NewDefaultDelegate()
	boardDelegate.ShowDescription = true
	boardDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	boardDelegate.Styles.SelectedTitle.Foreground(SelectedListItemColor).Border(lipgloss.HiddenBorder(), false, false, false, true)
	boardDelegate.Styles.SelectedDesc = boardDelegate.Styles.SelectedTitle.Copy()
}

func setMoveDelegate() {
	moveDelegate = list.NewDefaultDelegate()
	moveDelegate.Styles.NormalTitle.Foreground(ListItemColor)
	moveDelegate.Styles.SelectedTitle.
		Foreground(SelectedListItemColor).
		BorderLeft(false).
		BorderTop(true).
		BorderForeground(MoveBarColor)
	moveDelegate.Styles.SelectedDesc = boardDelegate.Styles.SelectedTitle.Copy().
		BorderTop(false)
}

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
		b := list.New([]list.Item{}, boardDelegate, ws.width/3, ws.height-9)
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
