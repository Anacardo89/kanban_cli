package ui

import (
	"log"

	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type InputField struct {
	field textinput.Model
	data  string
}

var (
	menuItems  []list.Item
	boardItems []list.Item
)

type Item struct {
	title       string
	description string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }
func (i Item) FilterValue() string { return i.title }

func (m *Menu) SetupList() {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), ws.width/3, ws.height-9)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	m.list = l
}

func (p *Project) SetupBoards() {
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
		b := list.New([]list.Item{}, list.NewDefaultDelegate(), ws.width/3, ws.height-9)
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
