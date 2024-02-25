package ui

import (
	"fmt"
	"log"
	"strconv"

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
	if p.flag == dBoard || p.flag == dCard {
		return p.viewCertify()
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

func (p *Project) getBoard() *kanban.Board {
	if p.empty {
		return nil
	}
	board, err := p.project.Boards.GetAt(p.cursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	return board.(*kanban.Board)
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

func (p *Project) viewCertify() string {
	var (
		toDelete interface{}
		areUsure string
		err      error
	)
	if p.flag == dBoard {
		toDelete = p.getBoard()
		if err != nil {
			log.Println(err)
		}
		areUsure = fmt.Sprintf(
			"Are you sure you wish to delete the board\n\n%s\n\nThis will also delete all cards in the board\nThis operation cannot be reverted\n\ny/N",
			toDelete.(*kanban.Board).Title,
		)
	} else {
		toDelete = p.getCard()
		if err != nil {
			log.Println(err)
		}
		areUsure = fmt.Sprintf(
			"Are you sure you wish to delete the card\n\n%s\n\nThis operation cannot be reverted\n\ny/N",
			toDelete.(*kanban.Card).Title,
		)
	}

	areUsureStyled := EmptyStyle.Render(areUsure)
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Center,
		areUsureStyled,
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
		return InputNoFieldStyle.Render("[Enter] append Card to Board or [Esc] to go back")
	}
	switch p.flag {
	case new:
		return InputNoFieldStyle.Render("Create new: [B]oard or [C]ard")
	case rename:
		return InputNoFieldStyle.Render("Rename: [P]roject or [B]oard")
	case delete:
		return InputNoFieldStyle.Render("Delete: [B]oard or [C]ard")
	case mvBoard:
		return InputNoFieldStyle.Render("Use [h/l] or [left/right] to move board. [Enter] to confirm position")
	}
	return "[hjkl][arrows] movement * [ESC] menu [i] label [ENTER] selected card * [N]ew [M]ove [R]ename [D]elete"
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
