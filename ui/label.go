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
	labelName string
)

type Label struct {
	project   *kanban.Project
	list      list.Model
	cursor    int
	Input     InputField
	inputFlag inputFlag
}

func OpenLabels(kp *kanban.Project) Label {
	l := Label{
		project: kp,
		Input:   InputField{field: textinput.New()},
		cursor:  0,
	}
	setLabelItemDelegate()
	l.setupList()
	return l
}

func (l Label) Init() tea.Cmd {
	return nil
}

func (l Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		l.setInput()
		l.setupList()
		return l, nil
	case tea.KeyMsg:
		if l.Input.field.Focused() {
			l.handleInput(msg.String())
			l.Input.field, cmd = l.Input.field.Update(msg)
			return l, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit
		case "n":
			l.inputFlag = title
			return l, l.Input.field.Focus()
		case "d":
			label := l.getLabel()
			if label == nil {
				return l, nil
			}
			l.deleteLabel()
			return l, nil
		case "b", "esc":
			return l, func() tea.Msg { return project }
		}
	}
	l.list, cmd = l.list.Update(msg)
	l.cursor = l.list.Cursor()
	return l, cmd
}

func (l Label) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	var (
		bottomLines    = ""
		emptyTxtStyled = ""
		inputStyled    = ""
		labelStyled    = ""
		output         = ""
	)

	if l.project.Labels.Length() == 0 {
		emptyTxt := "No labels.\n\nPress 'n' to create a new Label\nor 'esc' to go back"
		emptyTxtStyled = EmptyStyle.Render(emptyTxt)
		if l.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2; i++ {
				bottomLines += "\n"
			}
			inputStyled = InputFieldStyle.Render(l.Input.field.View())
		}
		output = lipgloss.Place(
			ws.width,
			ws.height,
			lipgloss.Center,
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Center,
				emptyTxtStyled,
				bottomLines,
				inputStyled,
			))
		return output
	}

	labelStyled = ListStyle.Render(l.list.View())
	if l.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(l.Input.field.View())
	}
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			labelStyled,
			bottomLines,
			inputStyled,
		))
	return output
}

// action
func (l *Label) getLabel() *kanban.Label {
	if l.project.Labels.Length() == 0 {
		return nil
	}
	node, err := l.project.Labels.WalkTo(l.cursor)
	if err != nil {
		log.Println(err)
	}
	return node.Val().(*kanban.Label)
}

func (l *Label) deleteLabel() {
	var err error
	if l.project.Labels.Length() == 0 {
		return
	}
	node, err := l.project.Labels.WalkTo(l.cursor)
	if err != nil {
		log.Println(err)
		return
	}
	label := node.Val().(*kanban.Label)
	err = l.project.RemoveLabel(label)
	if err != nil {
		log.Println(err)
	}
	l.setupList()
	l.cursor = l.list.Cursor()
}
