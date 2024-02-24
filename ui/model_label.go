package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tmpTitle = ""
	tmpColor = ""
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
	l.cursor = l.list.Cursor()
	return l, cmd
}

func (l Label) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	if l.project.Labels.Length() == 0 {
		return l.viewEmpty()
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
	l.setTxtInput()
	l.setList()
	return l
}

func (l *Label) getLabel() *kanban.Label {
	if l.project.Labels.Length() == 0 {
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
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "enter":
		return func() tea.Msg { return upLabel }
	case "b", "esc":
		return func() tea.Msg { return projectState }
	case "n":
		l.flag = title
		l.textinput.Placeholder = "Label Title"
		return l.textinput.Focus()
	case "d":
		l.deleteLabelFromCards()
		l.deleteLabel()
		return nil
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
	if l.flag == title {
		if l.textinput.Value() == "" {
			return
		}
		tmpTitle = l.textinput.Value()
		l.textinput.SetValue("")
		l.textinput.Placeholder = "Label Hex Color"
		l.flag = color
		return
	} else if l.flag == color {
		tmpColor = l.textinput.Value()
		if tmpColor[0] != '#' {
			tmpColor = string('#') + tmpColor
		}
		if len(tmpColor) != 7 {
			return
		}
		l.project.AddLabel(tmpTitle, tmpColor)
		l.setList()
		l.textinput.SetValue("")
		l.textinput.Blur()
		l.flag = none
	}
}

// actions
func (l *Label) deleteLabel() {
	var err error
	if l.empty {
		return
	}
	label, err := l.project.Labels.GetAt(l.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	err = l.project.RemoveLabel(label.(*kanban.Label))
	if err != nil {
		log.Println(err)
	}
	l.setList()
	l.cursor = l.list.Cursor()
}

func (l *Label) deleteLabelFromCards() {
	var err error
	if l.empty {
		return
	}
	label, err := l.project.Labels.GetAt(l.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	cards := l.getAllCards()
	for _, card := range cards {
		card.RemoveLabel(label.(*kanban.Label))
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

func (l *Label) viewLabels() string {
	var (
		bottomLines string
		inputStyled string
	)
	labelStyled := ListStyle.Render(l.list.View())
	if l.textinput.Focused() {
		inputStyled = InputFieldStyle.Render(l.textinput.View())
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
