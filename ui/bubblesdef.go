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
		err    error
		boards []list.Model
		node   *dll.Node
	)
	node, err = p.project.Boards.HeadNode()
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < p.project.Boards.Length(); i++ {
		b := list.New([]list.Item{}, list.NewDefaultDelegate(), ws.width/3, ws.height-9)
		b.SetShowHelp(false)
		board := node.Val().(*kanban.Board)
		b.Title = board.Title
		b.InfiniteScrolling = true
		node, err = node.Next()
		if err != nil {
			log.Println(err)
		}
		boards = append(boards, b)
	}
	p.boards = boards
}
