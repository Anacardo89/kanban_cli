package ui

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tmpTitle   = ""
	tmpColor   = ""
	recolorClr = ""
)

// Implements tea.Model
type Label struct {
	project   *kanban.Project
	list      list.Model
	cursor    int
	textinput textinput.Model
	flag      actionFlag
	empty     bool
}

func (l Label) Init() tea.Cmd {
	return nil
}

func (l Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		l.setList()
		return l, nil
	case tea.KeyMsg:
		cmd = l.keyPress(msg)
		return l, cmd
	}
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l Label) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	if l.empty {
		return l.viewEmpty()
	}
	if l.flag == delete {
		return l.viewCertify()
	}
	return l.viewLabels()
}

// called by model_selector
func OpenLabels(kp *kanban.Project) Label {
	l := Label{
		project:   kp,
		textinput: textinput.New(),
		cursor:    0,
		flag:      none,
	}
	if l.project.Labels.Length() == 0 {
		l.empty = true
	}
	l.setTxtInput()
	l.setList()
	return l
}

func (l *Label) getLabel() *kanban.Label {
	if l.empty {
		return nil
	}
	label, err := l.project.Labels.GetAt(l.list.Cursor())
	if err != nil {
		log.Println(err)
	}
	return label.(*kanban.Label)
}

// Update
func (l *Label) keyPress(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if l.textinput.Focused() {
		cmd = l.inputFocused(msg)
		return cmd
	}
	if l.flag == delete {
		switch msg.String() {
		case "n", "enter", "esc":
			l.flag = none
		case "y":
			label := l.getLabel()
			l.deleteLabelFromCards(label)
			l.deleteLabel(label)
			l.flag = none
		}
		return nil
	}
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "enter":
		if l.empty {
			return nil
		} else {
			return func() tea.Msg { return upLabel }
		}
	case "b", "esc":
		return func() tea.Msg { return upProject }
	case "n":
		l.flag = title
		l.textinput.Placeholder = "Label Title"
		return l.textinput.Focus()
	case "r":
		l.textinput.Placeholder = "New Label Title"
		l.flag = rename
		label := l.getLabel()
		tmpTitle = label.Title
		tmpColor = label.Color
		return l.textinput.Focus()
	case "d":
		l.flag = delete
		return nil
	case "c":
		l.textinput.Placeholder = "New Label Hex Color (without #)"
		l.flag = recolor
		label := l.getLabel()
		tmpTitle = label.Title
		tmpColor = label.Color
		return l.textinput.Focus()
	}
	l.list, cmd = l.list.Update(msg)
	return cmd
}

func (l *Label) inputFocused(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		l.textinput.SetValue("")
		l.textinput.Blur()
		tmpTitle = ""
		tmpColor = ""
		l.flag = none
	case "enter":
		l.txtInputEnter()
	}
	l.textinput, cmd = l.textinput.Update(msg)
	return cmd
}

func (l *Label) txtInputEnter() {
	switch l.flag {
	case title:
		if l.textinput.Value() == "" {
			l.flag = none
			return
		}
		tmpTitle = l.textinput.Value()
		l.textinput.SetValue("")
		l.textinput.Placeholder = "Label Hex Color (without #)"
		l.flag = color
		return
	case color:
		tmpColor = l.textinput.Value()
		_, err := hex.DecodeString(tmpColor)
		if err != nil {
			return
		}
		tmpColor = string('#') + tmpColor
		l.project.AddLabel(tmpTitle, tmpColor)
		l.empty = false
		l.setList()
		l.textinput.SetValue("")
		l.textinput.Blur()
		l.flag = none
	case rename:
		if l.textinput.Value() == "" {
			l.flag = none
			return
		}
		label := l.getLabel()
		label.Title = l.textinput.Value()
		l.setList()
		l.renameCardLabels()
		l.textinput.SetValue("")
		l.textinput.Blur()
		l.flag = none
		return
	case recolor:
		if l.textinput.Value() == "" {
			l.flag = none
			return
		}
		recolorClr = l.textinput.Value()
		_, err := hex.DecodeString(recolorClr)
		if err != nil {
			return
		}
		recolorClr = string('#') + recolorClr
		label := l.getLabel()
		label.Color = recolorClr
		l.setList()
		l.recolorCardLabels()
		l.textinput.SetValue("")
		l.textinput.Blur()
		l.flag = none
		return
	}
}

// actions
func (l *Label) renameCardLabels() {
	cards := l.getAllCards()
	for _, card := range cards {
		for i := 0; i < card.CardLabels.Length(); i++ {
			cl, err := card.CardLabels.GetAt(i)
			if err != nil {
				log.Println(err)
			}
			if cl.(*kanban.Label).Title == tmpTitle {
				cl.(*kanban.Label).Title = l.textinput.Value()
			}
		}
	}
}

func (l *Label) recolorCardLabels() {
	cards := l.getAllCards()
	for _, card := range cards {
		for i := 0; i < card.CardLabels.Length(); i++ {
			cl, err := card.CardLabels.GetAt(i)
			if err != nil {
				log.Println(err)
			}
			if cl.(*kanban.Label).Color == tmpTitle {
				cl.(*kanban.Label).Color = l.textinput.Value()
			}
		}
	}
}

func (l *Label) deleteLabel(label *kanban.Label) {
	var err error
	if l.empty {
		return
	}
	err = l.project.RemoveLabel(label)
	if err != nil {
		log.Println(err)
	}
	if l.project.Labels.Length() == 0 {
		l.empty = true
	}
	l.setList()
}

func (l *Label) deleteLabelFromCards(label *kanban.Label) {
	if l.empty {
		return
	}
	cards := l.getAllCards()
	for _, card := range cards {
		card.RemoveLabel(label)
	}
}

func (l *Label) getAllCards() []*kanban.Card {
	var cards []*kanban.Card
	for i := 0; i < l.project.Boards.Length(); i++ {
		b, err := l.project.Boards.GetAt(i)
		if err != nil {
			log.Println(err)
			return nil
		}
		for j := 0; j < b.(*kanban.Board).Cards.Length(); j++ {
			c, err := b.(*kanban.Board).Cards.GetAt(j)
			if err != nil {
				log.Println(err)
				return nil
			}
			cards = append(cards, c.(*kanban.Card))
		}
	}
	return cards
}

// View
func (l *Label) viewEmpty() string {
	var (
		bottomLines string
		inputStyled string
	)
	emptyTxtStyled := EmptyStyle.Render(
		"No labels.\n\nPress 'n' to create a new Label\nor 'esc' to go back",
	)
	if l.textinput.Focused() {
		_, h := lipgloss.Size(emptyTxtStyled)
		for i := 0; i < ws.height-h-h/2; i++ {
			bottomLines += "\n"
		}
		inputStyled = InputFieldStyle.Render(l.textinput.View())
	}
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Center,
			emptyTxtStyled,
			bottomLines,
			inputStyled,
		),
	)
}

func (l *Label) viewCertify() string {
	toDelete := l.getLabel()
	areUsure := fmt.Sprintf(
		"Are you sure you wish to delete label\n\n%s\n\nThe label will also be deleted from the cards\nThis operation cannot be reverted\n\ny/N",
		toDelete.Title,
	)
	areUsureStyled := EmptyStyle.Render(areUsure)
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Center,
		areUsureStyled,
	)
}

func (l *Label) viewLabels() string {
	var (
		bottomLines string
		inputStyled string
	)
	labelStyled := ListStyle.Render(l.list.View())
	if l.textinput.Focused() {
		inputStyled = InputFieldStyle.Render(l.textinput.View())
	} else {
		inputStyled = InputNoFieldStyle.Render(
			"[kj] [down/up] movement * [ESC] menu [B]oard * [N]ew [D]elete [R]ename [C] recolor * [ENTER] add to Card",
		)
	}
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Center,
			labelStyled,
			bottomLines,
			inputStyled,
		),
	)
}

// bubbles
// text input
func (l *Label) setTxtInput() {
	l.textinput.Prompt = ": "
	l.textinput.CharLimit = 120
	l.textinput.Cursor.Blink = true
}

// list
var labelDelegate = NewLabelListDelegate()

func (l *Label) setList() {
	var labelItems []list.Item
	lst := list.New([]list.Item{}, labelDelegate, ws.width/3, ws.height-9)
	lst.SetShowHelp(false)
	lst.Title = "Labels"
	lst.InfiniteScrolling = true
	for i := 0; i < l.project.Labels.Length(); i++ {
		label, _ := l.project.Labels.GetAt(i)
		item := Item{
			title:       label.(*kanban.Label).Title,
			description: label.(*kanban.Label).Color,
		}
		labelItems = append(labelItems, item)
	}
	lst.SetItems(labelItems)
	l.list = lst
}
