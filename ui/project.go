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

type Project struct {
	styles    []lipgloss.Style
	project   *kanban.Project
	boards    []list.Model
	hcursor   int
	vcursor   int
	sb        *dll.Node
	sc        *dll.Node
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
	p.SetStyles()
	p.SetupBoards()
	return p
}

func (p Project) Init() tea.Cmd {
	return nil
}

func (p Project) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ws.width = msg.Width
		ws.height = msg.Height
		p.SetStyles()
		p.SetupBoards()
		return p, nil
	case tea.KeyMsg:
		if p.Input.field.Focused() {
			p.handleInput(msg.String())
			p.Input.field, cmd = p.Input.field.Update(msg)
			return p, cmd
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
		case "l":
			return p, func() tea.Msg { return label }
		case "esc":
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
		emptyTxtStyled = p.styles[empty].Render(emptyTxt)
		if p.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2-1; i++ {
				bottomLines += "\n"
			}
			inputStyled = p.styles[input].Render(p.Input.field.View())
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
	titleStyled = p.styles[title].Render(p.project.Title)
	for i := range p.boards {
		if i == p.hcursor {
			boardStyled = p.styles[selected].Render(p.boards[i].View())
		} else {
			boardStyled = p.styles[listStyle].Render(p.boards[i].View())
		}
		boardsStyled = lipgloss.JoinHorizontal(lipgloss.Top, boardsStyled, boardStyled)
	}
	if p.Input.field.Focused() {
		inputStyled = p.styles[input].Render(p.Input.field.View())
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
	var err error
	var node *dll.Node
	if p.hcursor == 0 {
		p.hcursor = p.project.Boards.Length() - 1
		node, err = p.project.Boards.TailNode()
		if err != nil {
			log.Println(err)
		}
		p.sb = node
		return
	}
	p.hcursor--
	p.sb, err = p.sb.Prev()
	if err != nil {
		log.Println(err)
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	p.sc, err = p.project.Boards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
	}
}

func (p *Project) handleMoveRight() {
	if p.project.Boards.Length() == 0 {
		return
	}
	var err error
	var node *dll.Node
	if p.hcursor == p.project.Boards.Length()-1 {
		p.hcursor = 0
		node, err = p.project.Boards.HeadNode()
		if err != nil {
			log.Println(err)
		}
		p.sb = node
		return
	}
	p.hcursor++
	p.sb, err = p.sb.Next()
	if err != nil {
		log.Println(err)
	}
	p.vcursor = p.boards[p.hcursor].Cursor()
	p.sc, err = p.project.Boards.WalkTo(p.vcursor)
	if err != nil {
		log.Println(err)
	}
}

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
	return node.Val().(*kanban.Card)
}

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
		switch p.inputFlag {
		case new:
			p.project.AddBoard(p.Input.data)
			p.SetupBoards()
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
