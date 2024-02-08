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
	styles    []lipgloss.Style
	project   *kanban.Project
	boards    []list.Model
	hcursor   int
	vcursor   int
	mcursor   int
	sb        *dll.Node
	sc        *dll.Node
	mb        *dll.Node
	mc        *dll.Node
	Input     InputField
	inputFlag inputFlag
}

func OpenProject(kp *kanban.Project) Project {
	var err error
	p := Project{
		styles:  make([]lipgloss.Style, 10),
		Input:   InputField{field: textinput.New()},
		project: kp,
	}
	setBoardItemDelegate()
	setMoveDelegate()
	p.sb, err = p.project.Boards.WalkTo(p.hcursor)
	if err != nil {
		log.Println(err)
	}
	if p.sb != (*dll.Node)(nil) {
		board := p.sb.Val().(*kanban.Board)
		p.sc, err = board.Cards.WalkTo(p.vcursor)
		if err != nil {
			log.Println(err)
		}
	}
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
		case "up":
			p.handleMoveUp()
			p.boards[p.hcursor], cmd = p.boards[p.hcursor].Update(msg)
			return p, cmd
		case "down":
			p.handleMoveDown()
			p.boards[p.hcursor], cmd = p.boards[p.hcursor].Update(msg)
			return p, cmd
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
				p.boards[p.hcursor].SetDelegate(boardDelegate)
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
			p.boards[p.hcursor].SetDelegate(moveDelegate)
			p.mb = p.sb
			p.mc = p.sc
			p.mcursor = p.vcursor
			return p, nil
		case "l":
			return p, func() tea.Msg { return label }
		case "esc":
			if p.inputFlag != none {
				p.boards[p.hcursor].SetDelegate(boardDelegate)
				p.inputFlag = none
				return p, nil
			}
			return p, func() tea.Msg { return menu }
		}
	}
	return p, nil
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
func (p *Project) handleMoveUp() {
	board := p.sb.Val().(*kanban.Board)
	if board.Cards.Length() == 0 {
		return
	}
	var err error
	var node *dll.Node
	if p.vcursor == 0 {
		p.vcursor = board.Cards.Length() - 1
		node, err = board.Cards.TailNode()
		if err != nil {
			log.Println(err)
		}
		p.sc = node
		return
	}
	p.vcursor--
	p.sc, err = p.sc.Prev()
	if err != nil {
		log.Println(err)
	}
}

func (p *Project) handleMoveDown() {
	board := p.sb.Val().(*kanban.Board)
	if board.Cards.Length() == 0 {
		return
	}
	var err error
	var node *dll.Node
	if p.vcursor == board.Cards.Length()-1 {
		p.vcursor = 0
		node, err = board.Cards.HeadNode()
		if err != nil {
			log.Println(err)
		}
		p.sc = node
		return
	}
	p.vcursor++
	p.sc, err = p.sc.Next()
	if err != nil {
		log.Println(err)
	}
}

func (p *Project) handleMoveLeft() {
	if p.project.Boards.Length() == 0 {
		return
	}
	var (
		err  error
		node *dll.Node
	)
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(boardDelegate)
	}
	if p.hcursor == 0 {
		p.hcursor = p.project.Boards.Length() - 1
		if p.inputFlag == move {
			p.boards[p.hcursor].SetDelegate(moveDelegate)
		}
		node, err = p.project.Boards.TailNode()
		if err != nil {
			log.Println(err)
		}
		p.sb = node
		return
	}
	p.hcursor--
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(moveDelegate)
	}
	p.sb, err = p.sb.Prev()
	if err != nil {
		log.Println(err)
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	board := p.sb.Val().(*kanban.Board)
	p.sc, err = board.Cards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
	}
}

func (p *Project) handleMoveRight() {
	if p.project.Boards.Length() == 0 {
		return
	}
	var (
		err  error
		node *dll.Node
	)
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(boardDelegate)
	}
	if p.hcursor == p.project.Boards.Length()-1 {
		p.hcursor = 0
		if p.inputFlag == move {
			p.boards[p.hcursor].SetDelegate(moveDelegate)
		}
		node, err = p.project.Boards.HeadNode()
		if err != nil {
			log.Println(err)
		}
		p.sb = node
		return
	}
	p.hcursor++
	if p.inputFlag == move {
		p.boards[p.hcursor].SetDelegate(moveDelegate)
	}
	p.sb, err = p.sb.Next()
	if err != nil {
		log.Println(err)
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	board := p.sb.Val().(*kanban.Board)
	p.sc, err = board.Cards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
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
	}
	board := node.Val().(*kanban.Board)
	node, err = board.Cards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
	}
	p.sc = node
	return node.Val().(*kanban.Card)
}

func (p *Project) deleteBoard() {
	var err error
	if p.project.Boards.Length() == 0 {
		return
	}
	b := p.sb.Val().(*kanban.Board)
	err = p.project.RemoveBoard(b)
	if err != nil {
		log.Println(err)
	}
	p.setupBoards()
}

func (p *Project) deleteCard() {
	var err error
	board := p.sb.Val().(*kanban.Board)
	if board.Cards.Length() == 0 {
		return
	}
	c := p.sc.Val().(*kanban.Card)
	err = board.RemoveCard(c)
	if err != nil {
		log.Println(err)
	}
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
	b1 := p.mb.Val().(*kanban.Board)
	b2 := p.sb.Val().(*kanban.Board)
	card := p.mc.Val().(*kanban.Card)
	p.project.MoveCard(b1, b2, card, p.vcursor, p.mcursor)
	p.setupBoards()
}
