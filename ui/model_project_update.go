package ui

import (
	"strings"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/storage"
	tea "github.com/charmbracelet/bubbletea"
)

// called by project.Update()
func (p *Project) inputFocused(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		p.textinput.SetValue("")
		p.textinput.Blur()
		p.flag = none
	case "enter":
		p.txtInputEnter()
		p.flag = none
	}
	p.textinput, cmd = p.textinput.Update(msg)
	return cmd
}

func (p *Project) checkFlag(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch p.flag {
	case new:
		cmd = p.flagNew(msg)
		return cmd
	case move:
		cmd = p.flagMove(msg)
		return cmd
	case mvBoard:
		cmd = p.flagMvBoard(msg)
		return cmd
	case mvCard:
		cmd = p.flagMvCard(msg)
		return cmd
	case rename:
		cmd = p.flagRename(msg)
		return cmd
	case delete:
		cmd = p.flagDelete(msg)
		return cmd
	case dBoard:
		cmd = p.flagDBoard(msg)
		return cmd
	case dCard:
		cmd = p.flagDCard(msg)
		return cmd
	}
	p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	return cmd
}

func (p *Project) keyPress(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "k", "j", "up", "down":
		if p.empty {
			return nil
		}
	case "h", "left":
		if !p.empty {
			p.moveLeft()
		}
		return nil
	case "l", "right":
		if !p.empty {
			p.moveRight()
		}
		return nil
	case "i":
		return func() tea.Msg { return labelState }
	case "enter":
		cmd = p.keyPressEnter()
		return cmd
	case "esc":
		return func() tea.Msg { return upMenu }
	case "n":
		cmd = p.keyPressN()
		return cmd
	case "m":
		if !p.empty {
			p.flag = move
		}
		return nil
	case "r":
		if !p.empty {
			p.flag = rename
		}
		return nil
	case "d":
		if !p.empty {
			p.flag = delete
		}
		return nil
	}
	p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	return cmd
}

// textinput
func (p *Project) txtInputEnter() {
	if p.textinput.Value() == "" {
		return
	}
	switch p.flag {
	case nBoard:
		res := storage.CreateBoard(p.textinput.Value(), p.project.Id)
		id64, err := res.LastInsertId()
		if err != nil {
			logger.Error.Fatal(err)
		}
		id := int(id64)
		p.project.AddBoard(id, p.textinput.Value())
		if p.empty {
			p.empty = false
		}
		p.emptyBoard = append(p.emptyBoard, true)
		p.cursor = 0
	case nCard:
		b := p.getBoard()
		res := storage.CreateCard(p.textinput.Value(), b.Id)
		id64, err := res.LastInsertId()
		if err != nil {
			logger.Error.Fatal(err)
		}
		id := int(id64)
		b.AddCard(id, p.textinput.Value(), "")
		p.emptyBoard[p.cursor] = false
	case rename:
		if strings.Contains(p.textinput.Placeholder, "Project") {
			storage.UpdateProject(p.project.Id, p.textinput.Value())
			p.project.RenameProject(p.textinput.Value())
		} else {
			b := p.getBoard()
			storage.UpdateBoard(b.Id, p.textinput.Value())
			b.RenameBoard(p.textinput.Value())
		}
	}
	p.flag = none
	p.setLists()
	p.textinput.SetValue("")
	p.textinput.Blur()
}

// actionFlag
func (p *Project) flagNew(msg tea.KeyMsg) tea.Cmd {
	if p.empty {
		p.textinput.Placeholder = BOARD_TITLE
		return p.textinput.Focus()
	} else {
		switch msg.String() {
		case "b":
			p.flag = nBoard
			p.textinput.Placeholder = BOARD_TITLE
			return p.textinput.Focus()
		case "c":
			p.flag = nCard
			p.textinput.Placeholder = CARD_TITLE
			return p.textinput.Focus()
		}
	}
	return nil
}

func (p *Project) flagMove(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "b":
		p.flag = mvBoard
	case "c":
		if p.emptyBoard[p.cursor] {
			p.flag = none
			return nil
		}
		p.flag = mvCard
		p.moveFrom = []int{p.cursor, p.boards[p.cursor].Cursor()}
	}
	return nil
}

func (p *Project) flagMvBoard(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "h", "left":
		p.moveBoardLeft()
	case "l", "right":
		p.moveBoardRight()
	case "enter", "esc":
		p.flag = none
		p.moveFrom = []int{-1, 0}
		return nil
	}
	p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	return cmd
}

func (p *Project) flagMvCard(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "h", "left":
		p.moveLeft()
	case "l", "right":
		p.moveRight()
	case "esc":
		p.flag = none
		p.moveFrom = []int{-1, 0}
		return nil
	case "enter":
		p.moveCard()
		p.moveFrom = []int{-1, 0}
		p.flag = none
	}
	p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	return cmd
}

func (p *Project) flagRename(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "p":
		p.textinput.Placeholder = PROJECT_TITLE
		return p.textinput.Focus()
	case "b":
		p.textinput.Placeholder = BOARD_TITLE
		return p.textinput.Focus()
	}
	return nil
}

func (p *Project) flagDelete(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "b":
		p.flag = dBoard
	case "c":
		if !p.emptyBoard[p.cursor] {
			p.flag = dCard
		} else {
			p.flag = none
		}
	}
	return nil
}

func (p *Project) flagDBoard(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "n", "enter", "esc":
		p.flag = none
	case "y":
		p.deleteBoard()
		p.flag = none
	}
	return nil
}

func (p *Project) flagDCard(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "n", "enter", "esc":
		p.flag = none
	case "y":
		p.deleteCard()
		p.flag = none
	}
	return nil
}

// key presses
func (p *Project) keyPressEnter() tea.Cmd {
	if p.empty {
		return nil
	}
	if p.emptyBoard[p.cursor] {
		return nil
	}
	return func() tea.Msg { return cardState }
}

// actions
// movement
func (p *Project) moveLeft() {
	if p.empty {
		return
	}
	if p.cursor == 0 {
		p.cursor = p.project.Boards.Length() - 1
	} else {
		p.cursor--
	}
}

func (p *Project) moveRight() {
	if p.empty {
		return
	}
	if p.cursor == p.project.Boards.Length()-1 {
		p.cursor = 0
	} else {
		p.cursor++
	}
}

// new
func (p *Project) keyPressN() tea.Cmd {
	if p.empty {
		p.flag = nBoard
		return p.textinput.Focus()
	} else {
		p.flag = new
		return nil
	}
}

// move
func (p *Project) moveBoardLeft() {
	b := p.getBoard()
	bVal := *b
	p.project.Boards.RemoveAt(p.cursor)
	if p.cursor == 0 {
		p.project.Boards.Append(&bVal)
		p.cursor = p.project.Boards.Length() - 1
	} else {
		p.project.Boards.InsertAt(p.cursor-1, &bVal)
		p.cursor--
	}
	p.setLists()
}

func (p *Project) moveBoardRight() {
	b := p.getBoard()
	bVal := *b
	p.project.Boards.RemoveAt(p.cursor)
	if p.cursor == p.project.Boards.Length() {
		p.project.Boards.Prepend(&bVal)
		p.cursor = 0
	} else {
		p.project.Boards.InsertAt(p.cursor+1, &bVal)
		p.cursor++
	}
	p.setLists()
}

func (p *Project) moveCard() {
	boardFrom, err := p.project.Boards.GetAt(p.moveFrom[0])
	if err != nil {
		logger.Error.Fatal(err)
		return
	}
	boardTo := p.getBoard()
	c, err := boardFrom.(*kanban.Board).Cards.GetAt(p.moveFrom[1])
	if err != nil {
		logger.Error.Fatal(err)
		return
	}
	cardVal := *c.(*kanban.Card)
	storage.UpdateCardParent(cardVal.Id, boardTo.Id)
	boardFrom.(*kanban.Board).Cards.RemoveAt(p.moveFrom[1])
	boardTo.Cards.Append(&cardVal)
	p.setLists()
}

// delete
func (p *Project) deleteBoard() {
	if p.empty {
		return
	}
	b := p.getBoard()
	storage.DeleteBoard(b.Id)
	err := p.project.RemoveBoard(b)
	if err != nil {
		logger.Error.Fatal(err)
	}
	if p.project.Boards.Length() == 0 {
		p.empty = true
	}
	p.cursor = 0
	p.setLists()
}

func (p *Project) deleteCard() {
	b := p.getBoard()
	c := p.getCard()
	storage.DeleteCard(c.Id)
	err := b.RemoveCard(c)
	if err != nil {
		logger.Error.Fatal(err)
	}
	p.emptyBoard[p.cursor] = true
	p.setLists()
}
