package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/textinput"
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
		m.menu.AddProject(m.Input.data)
		m.setupList()
		m.Input.data = ""
		m.Input.field.SetValue("")
		m.Input.field.Blur()
		m.cursor = m.list.Cursor()
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
	l.Input.field.Placeholder = "Label Title"
}

func (l *Label) handleInput(key string) {
	switch key {
	case "esc":
		l.Input.field.SetValue("")
		l.Input.data = ""
		l.Input.field.Blur()
		return
	case "enter":
		if l.inputFlag == title {
			l.Input.data = l.Input.field.Value()
			if l.Input.data == "" {
				return
			}
			labelName = l.Input.data
			l.Input.data = ""
			l.Input.field.SetValue("")
			l.Input.field.Placeholder = "Label Hex Color"
			l.inputFlag = color
			return
		}
		if l.inputFlag == color {
			l.Input.data = l.Input.field.Value()
			if l.Input.data[0] != '#' {
				l.Input.data = string('#') + l.Input.data
			}
			if len(l.Input.data) != 7 {
				return
			}
			l.project.AddLabel(labelName, l.Input.data)
			l.setupList()
			l.Input.data = ""
			l.Input.field.SetValue("")
			l.Input.field.Blur()
			l.inputFlag = none
		}
	}
}
