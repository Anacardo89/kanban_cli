package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// implements tea.Model
type Project struct {
	project    *kanban.Project
	boards     []list.Model
	empty      bool
	emptyBoard []bool
	cursor     int
	moveFrom   []int
	textinput  textinput.Model
	flag       actionFlag
}

func (p Project) Init() tea.Cmd {
	return nil
}

func (p Project) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		p.setLists()
		return p, nil
	case tea.KeyMsg:
		if p.textinput.Focused() {
			cmd = p.inputFocused(msg)
			return p, cmd
		}
		if p.flag != none {
			cmd = p.checkFlag(msg)
			return p, cmd
		}
		cmd = p.keyPress(msg)
		return p, cmd
	}
	if !p.empty {
		p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	}
	return p, cmd
}

func (p Project) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	if p.empty {
		return p.viewEmpty()
	}
	return p.viewBoards()
}

// called by model_selector
func OpenProject(kp *kanban.Project) Project {
	p := Project{
		textinput: textinput.New(),
		project:   kp,
		cursor:    0,
		moveFrom:  []int{-1, 0},
		flag:      none,
	}
	p.setEmptyBoard()
	p.setInput()
	setMoveDelegate()
	p.setLists()
	return p
}

func (p *Project) setEmptyBoard() {
	var emptyBoard []bool
	if p.project.Boards.Length() == 0 {
		p.empty = true
	}
	for i := 0; i < p.project.Boards.Length(); i++ {
		p.empty = false
		b, err := p.project.Boards.GetAt(i)
		if err != nil {
			log.Println(err)
			return
		}
		if b.(*kanban.Board).Cards.Length() == 0 {
			emptyBoard = append(emptyBoard, true)
		} else {
			emptyBoard = append(emptyBoard, false)
		}
	}
	p.emptyBoard = emptyBoard
}

func (p *Project) UpdateProject() {
	p.setEmptyBoard()
	p.setLists()
}

func (p *Project) getCard() *kanban.Card {
	if p.empty {
		return nil
	}
	if p.emptyBoard[p.cursor] {
		return nil
	}
	board, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	card, err := board.(*kanban.Board).Cards.GetAt(p.boards[p.cursor].Cursor())
	if err != nil {
		log.Println(err)
		return nil
	}
	return card.(*kanban.Card)
}

// Update
// textinput
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

func (p *Project) txtInputEnter() {
	if p.textinput.Value() == "" {
		return
	}
	switch p.flag {
	case board:
		p.project.AddBoard(p.textinput.Value())
		if p.empty {
			p.empty = false
		}
		p.emptyBoard = append(p.emptyBoard, true)
		p.cursor = 0
	case card:
		board, err := p.project.Boards.GetAt(p.cursor)
		if err != nil {
			log.Println(err)
			err = nil
			return
		}
		board.(*kanban.Board).AddCard(p.textinput.Value())
		p.emptyBoard[p.cursor] = false
	case rename:
		if strings.Contains(p.textinput.Placeholder, "Project") {
			p.project.RenameProject(p.textinput.Value())
		} else {
			b, err := p.project.Boards.GetAt(p.cursor)
			if err != nil {
				log.Println(err)
				return
			}
			b.(*kanban.Board).RenameBoard(p.textinput.Value())
		}
	}
	p.flag = none
	p.setLists()
	p.textinput.SetValue("")
	p.textinput.Blur()
}

// actionFlag
func (p *Project) checkFlag(msg tea.KeyMsg) tea.Cmd {
	switch p.flag {
	case new:
		if p.empty {
			p.textinput.Placeholder = BOARD_TITLE
			return p.textinput.Focus()
		} else {
			switch msg.String() {
			case "b":
				p.flag = board
				p.textinput.Placeholder = BOARD_TITLE
				return p.textinput.Focus()
			case "c":
				p.flag = card
				p.textinput.Placeholder = CARD_TITLE
				return p.textinput.Focus()
			}
		}
	case move:
		p.moveLogic(msg)
	case rename:
		switch msg.String() {
		case "p":
			p.textinput.Placeholder = PROJECT_TITLE
			return p.textinput.Focus()
		case "b":
			p.textinput.Placeholder = BOARD_TITLE
			return p.textinput.Focus()
		}

	case delete:
		switch msg.String() {
		case "b":
			p.deleteBoard()
		case "c":
			p.deleteCard()
		}
		p.flag = none
		return nil
	case board:
		switch msg.String() {
		case "left":
			p.moveBoardLeft()
		case "right":
			p.moveBoardRight()
		case "enter", "esc":
			p.moveFrom = []int{-1, 0}
		}
		return nil
	case card:
		p.moveFrom = []int{p.cursor, p.boards[p.cursor].Cursor()}
		p.flag = move
	}
	return nil
}

func (p *Project) moveLogic(msg tea.KeyMsg) {
	if p.moveFrom[0] == -1 {
		switch msg.String() {
		case "b":
			p.flag = board
		case "c":
			p.flag = card
		}
		return
	} else {
		switch msg.String() {
		case "esc":
			p.flag = none
			return
		case "enter":
			if p.flag == card {
				p.moveCard()
				p.moveFrom = []int{-1, 0}
			}
			p.flag = none
			return
		}
	}
}

// key presses
func (p *Project) keyPress(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "left":
		p.moveLeft()
		p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
		return cmd
	case "right":
		p.moveRight()
		p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
		return cmd
	case "l":
		return func() tea.Msg { return labelState }
	case "enter":
		if p.empty {
			return nil
		}
		if p.emptyBoard[p.cursor] {
			return nil
		}
		return func() tea.Msg { return cardState }
	case "esc":
		return func() tea.Msg { return upMenu }
	case "n":
		if p.empty {
			p.flag = board
			return p.textinput.Focus()
		} else {
			p.flag = new
			return nil
		}
	case "m":
		p.flag = move
		return nil
	case "r":
		p.flag = rename
		return nil
	case "d":
		p.flag = delete
		return nil

	}
	p.boards[p.cursor], cmd = p.boards[p.cursor].Update(msg)
	return cmd
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

// move
func (p *Project) moveBoardLeft() {
	b, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	bVal := *b.(*kanban.Board)
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
	b, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	bVal := *b.(*kanban.Board)
	p.project.Boards.RemoveAt(p.cursor)
	if p.cursor == p.project.Boards.Length() {
		p.project.Boards.Prepend(&bVal)
		p.cursor = 0
	} else {
		p.project.Boards.InsertAt(p.cursor, &bVal)
		p.cursor++
	}
	p.setLists()
}

func (p *Project) moveCard() {
	bf, err := p.project.Boards.GetAt(p.moveFrom[0])
	if err != nil {
		log.Println(err)
		return
	}
	bt, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	c, err := bf.(*kanban.Board).Cards.GetAt(p.moveFrom[1])
	if err != nil {
		log.Println(err)
		return
	}
	cardVal := *c.(*kanban.Card)
	bf.(*kanban.Board).Cards.RemoveAt(p.moveFrom[1])
	bt.(*kanban.Board).Cards.Append(&cardVal)
	p.setLists()
}

// delete
func (p *Project) deleteBoard() {
	if p.empty {
		return
	}
	node, err := p.project.Boards.WalkTo(p.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	b := node.Val().(*kanban.Board)
	err = p.project.RemoveBoard(b)
	if err != nil {
		log.Println(err)
	}
	if p.project.Boards.Length() == 0 {
		p.empty = true
	}
	p.cursor = 0
	p.setLists()
}

func (p *Project) deleteCard() {
	b, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	board := b.(*kanban.Board)
	c, err := board.Cards.GetAt(p.boards[p.cursor].Cursor())
	if err != nil {
		log.Println(err)
		return
	}
	card := c.(*kanban.Card)
	err = board.RemoveCard(card)
	if err != nil {
		log.Println(err)
	}
	p.setLists()
}

// View
func (p *Project) viewEmpty() string {
	var (
		bottomLines string
		inputStyled string
	)
	titleStyled := ProjectTitle.Render(p.project.Title)
	emptyTxtStyled := EmptyStyle.Render(
		"No Boards.\n\nPress 'n' to create a new Project Board\nor 'q' to quit",
	)
	if p.textinput.Focused() {
		_, h := lipgloss.Size(emptyTxtStyled)
		for i := 0; i < ws.height-h-h/2-1; i++ {
			bottomLines += "\n"
		}
		inputStyled = InputFieldStyle.Render(p.textinput.View())
	}
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Center,
			titleStyled,
			emptyTxtStyled,
			bottomLines,
			inputStyled,
		),
	)
}

func (p *Project) viewBoards() string {
	var (
		boardStyled  string
		boardsStyled string
	)
	titleStyled := ProjectTitle.Render(p.project.Title)
	for i := range p.boards {
		if i == p.cursor {
			boardStyled = ProjectListSelectedStyle.Render(p.boards[i].View())
		} else {
			boardStyled = ProjectListStyle.Render(p.boards[i].View())
		}
		boardsStyled = lipgloss.JoinHorizontal(lipgloss.Top, boardsStyled, boardStyled)
	}
	inputStyled := p.inputStyled()
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Left, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyled,
			boardsStyled,
			inputStyled,
		),
	)

}

func (p *Project) inputStyled() string {
	if p.textinput.Focused() {
		return InputFieldStyle.Render(p.textinput.View())
	} else if p.flag == move && p.moveFrom[0] == -1 {
		return InputNoFieldStyle.Render("Move: [B]oard or [C]ard")
	} else if p.flag == move && p.moveFrom[0] != -1 {
		return InputNoFieldStyle.Render("Press [Enter] to append Card to Board or [Esc] to go back")
	}
	switch p.flag {
	case new:
		return InputNoFieldStyle.Render("Create new: [B]oard or [C]ard")
	case rename:
		return InputNoFieldStyle.Render("Rename: [P]roject or [B]oard")
	case delete:
		return InputNoFieldStyle.Render("Delete: [B]oard or [C]ard")
	case board:
		return InputNoFieldStyle.Render("Use [Left] and [Right] to move highlighted board. [Enter] to confirm position")
	}
	return ""
}

// bubbles
// text input
func (p *Project) setInput() {
	p.textinput.Prompt = ": "
	p.textinput.CharLimit = 120
	p.textinput.Cursor.Blink = true
}

// list
var (
	cardDelegate = NewCardDelegate()
	moveDelegate CardDelegate
)

func setMoveDelegate() {
	moveDelegate = NewCardDelegate()
	moveDelegate.Styles.NormalTitle.Foreground(WHITE)
	moveDelegate.Styles.SelectedTitle.
		Foreground(YELLOW).
		Padding(0, 0, 0, 2).
		BorderLeft(false).
		BorderTop(true).
		BorderForeground(GREEN)
}

func (p *Project) setLists() {
	node, _ := p.project.Boards.HeadNode()
	if node == nil {
		return
	}
	var boards []list.Model
	board := node.Val().(*kanban.Board)
	for i := 0; i < p.project.Boards.Length(); i++ {
		l, err := p.listFromBoard(board)
		if err != nil {
			log.Println(err)
			return
		}
		boards = append(boards, l)
		node, _ = node.Next()
		if node != nil {
			board = node.Val().(*kanban.Board)
		}
	}
	p.boards = boards
}

func (p *Project) listFromBoard(b *kanban.Board) (list.Model, error) {
	var items []list.Item
	l := list.New([]list.Item{}, cardDelegate, ws.width/3, ws.height-9)
	l.SetShowHelp(false)
	l.Title = b.Title
	l.InfiniteScrolling = true
	for i := 0; i < b.Cards.Length(); i++ {
		c, err := b.Cards.GetAt(i)
		if err != nil {
			log.Println(err)
			return list.Model{}, err
		}
		checkTotal, checkDone := p.getCheckListInfo(c.(*kanban.Card))
		labelLen, metaSlice := p.getLabelInfo(c.(*kanban.Card))
		item := Item{
			title: c.(*kanban.Card).Title,
			description: fmt.Sprintf("[✓]%s/%s %sL",
				strconv.Itoa(checkDone),
				strconv.Itoa(checkTotal),
				strconv.Itoa(labelLen),
			),
			meta: metaSlice,
		}
		items = append(items, item)
		l.SetItems(items)
	}
	return l, nil
}

func (p *Project) getCheckListInfo(c *kanban.Card) (int, int) {
	checkDone := 0
	checkTotal := c.CheckList.Length()
	for i := 0; i < checkTotal; i++ {
		ci, err := c.CheckList.GetAt(i)
		if err != nil {
			log.Println(err)
		}
		if ci.(*kanban.CheckItem).Check {
			checkDone++
		}
	}
	return checkTotal, checkDone
}

func (p *Project) getLabelInfo(c *kanban.Card) (int, []Meta) {
	var metaSlice []Meta
	labelLen := c.CardLabels.Length()
	for i := 0; i < labelLen; i++ {
		l, err := c.CardLabels.GetAt(i)
		if err != nil {
			log.Println(err)
		}
		meta := Meta{
			initial: string(l.(*kanban.Label).Title[0]),
			color:   l.(*kanban.Label).Color,
		}
		metaSlice = append(metaSlice, meta)
	}
	return labelLen, metaSlice
}
