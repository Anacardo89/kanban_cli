package ui

import (
	"log"

	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// implements tea.Model
type Project struct {
	project   *kanban.Project
	boards    []list.Model
	hcursor   int
	vcursor   int
	moveFrom  []int
	Input     InputField
	inputFlag inputFlag
}

func OpenProject(kp *kanban.Project) Project {
	p := Project{
		Input:   InputField{field: textinput.New()},
		project: kp,
		hcursor: 0,
	}
	setBoardItemDelegate()
	setMoveDelegate()
	p.setupBoards()
	return p
}

func (p Project) Init() tea.Cmd {
	return nil
}

func (p Project) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		p.setupBoards()
		return p, nil
	case tea.KeyMsg:
		if p.Input.field.Focused() {
			p.handleInput(msg.String())
			p.Input.field, cmd = p.Input.field.Update(msg)
			return p, cmd
		}
		if p.inputFlag == delete {
			p.handleDelete(msg.String())
			return p, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		case "right":
			p.handleMoveRight()
			p.boards[p.hcursor], cmd = p.boards[p.hcursor].Update(msg)
			return p, cmd
		case "left":
			p.handleMoveLeft()
			p.boards[p.hcursor], cmd = p.boards[p.hcursor].Update(msg)
			return p, cmd
		case "enter":
			if p.inputFlag == move {
				p.handleMove()
				p.boards[p.hcursor].SetDelegate(DescDelegate)
				return p, nil
			}
			c := p.getCard()
			if c == nil {
				return p, nil
			}
			return p, func() tea.Msg { return card }
		case "n":
			p.inputFlag = new
			p.setInput()
			return p, p.Input.field.Focus()
		case "a":
			if p.project.Boards.Length() == 0 {
				return p, nil
			}
			p.inputFlag = add
			p.setInput()
			return p, p.Input.field.Focus()
		case "d":
			p.inputFlag = delete
			return p, nil
		case "m":
			p.inputFlag = move
			p.moveFrom = []int{p.hcursor, p.vcursor}
			p.boards[p.hcursor].SetDelegate(TopWhiteDelegate)
			return p, nil
		case "l":
			return p, func() tea.Msg { return label }
		case "esc":
			if p.inputFlag == move {
				p.boards[p.hcursor].SetDelegate(DescDelegate)
				p.inputFlag = none
				return p, nil
			}
			return p, func() tea.Msg { return menu }
		}
	}
	if p.project.Boards.Length() > 0 {
		p.boards[p.hcursor], cmd = p.boards[p.hcursor].Update(msg)
		p.vcursor = p.boards[p.hcursor].Cursor()
	}
	return p, cmd
}

func (p Project) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	titleStyled := ""
	emptyTxtStyled := ""
	bottomLines := ""
	inputStyled := ""
	boardStyled := ""
	boardsStyled := ""
	output := ""
	if p.project.Boards.Length() == 0 {
		emptyTxt := "No Boards.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		emptyTxtStyled = EmptyStyle.Render(emptyTxt)
		if p.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2-1; i++ {
				bottomLines += "\n"
			}
			inputStyled = InputFieldStyle.Render(p.Input.field.View())
		}
		output = lipgloss.Place(
			ws.width,
			ws.height,
			lipgloss.Center,
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Center,
				titleStyled,
				emptyTxtStyled,
				bottomLines,
				inputStyled,
			))
		return output
	}
	titleStyled = ProjectTitle.Render(p.project.Title)
	for i := range p.boards {
		if i == p.hcursor {
			boardStyled = ProjectListSelectedStyle.Render(p.boards[i].View())
		} else {
			boardStyled = ProjectListStyle.Render(p.boards[i].View())
		}
		boardsStyled = lipgloss.JoinHorizontal(lipgloss.Top, boardsStyled, boardStyled)
	}
	if p.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(p.Input.field.View())
	}
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Left,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyled,
			boardsStyled,
			inputStyled,
		))
	return output
}

// movement
func (p *Project) handleMoveLeft() {
	if p.project.Boards.Length() == 0 {
		return
	}
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(DescDelegate)
	}
	if p.hcursor == 0 {
		p.hcursor = p.project.Boards.Length() - 1
	} else {
		p.hcursor--
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(TopWhiteDelegate)
	}
}

func (p *Project) handleMoveRight() {
	if p.project.Boards.Length() == 0 {
		return
	}
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(DescDelegate)
	}
	if p.hcursor == p.project.Boards.Length()-1 {
		p.hcursor = 0
	} else {
		p.hcursor++
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(TopWhiteDelegate)
	}
}

// action
func (p *Project) getCard() *kanban.Card {
	if p.project.Boards.Length() == 0 {
		return nil
	}
	if len(p.boards[p.hcursor].Items()) == 0 {
		return nil
	}
	node, err := p.project.Boards.WalkTo(p.hcursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	board := node.Val().(*kanban.Board)
	node, err = board.Cards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	return node.Val().(*kanban.Card)
}

func (p *Project) deleteBoard() {
	if p.project.Boards.Length() == 0 {
		return
	}
	node, err := p.project.Boards.WalkTo(p.hcursor)
	if err != nil {
		log.Println(err)
		return
	}
	b := node.Val().(*kanban.Board)
	err = p.project.RemoveBoard(b)
	if err != nil {
		log.Println(err)
	}
	p.hcursor = 0
	p.setupBoards()
}

func (p *Project) deleteCard() {
	node, err := p.project.Boards.WalkTo(p.hcursor)
	if err != nil {
		log.Println(err)
		return
	}
	board := node.Val().(*kanban.Board)
	node, err = board.Cards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
		return
	}
	c := node.Val().(*kanban.Card)
	err = board.RemoveCard(c)
	if err != nil {
		log.Println(err)
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	p.setupBoards()
}

func (p *Project) handleDelete(key string) {
	switch key {
	case "b":
		p.deleteBoard()
		p.inputFlag = none
	case "c":
		p.deleteCard()
		p.inputFlag = none
	}
}

func (p *Project) handleMove() {
	var (
		node *dll.Node
		err  error
	)
	node, err = p.project.Boards.WalkTo(p.moveFrom[0])
	if err != nil {
		log.Println(err)
		return
	}
	boardFrom := node.Val().(*kanban.Board)
	node, err = p.project.Boards.WalkTo(p.hcursor)
	if err != nil {
		log.Println(err)
		return
	}
	boardTo := node.Val().(*kanban.Board)
	node, err = boardFrom.Cards.WalkTo(p.moveFrom[1])
	if err != nil {
		log.Println(err)
		return
	}
	card := node.Val().(*kanban.Card)
	cardVal := *card
	boardFrom.Cards.RemoveAt(p.moveFrom[1])
	boardTo.Cards.InsertAt(p.vcursor, &cardVal)
	p.setupBoards()
}
