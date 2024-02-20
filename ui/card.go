package ui

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	card      *kanban.Card
	checklist list.Model
	labels    list.Model
	hcursor   int
	vcursor   int
	sc        *dll.Node
	Input     InputField
	inputFlag inputFlag
}

func OpenCard(kc *kanban.Card) Card {
	c := Card{
		card:  kc,
		Input: InputField{field: textinput.New()},
	}
	c.setupLists()
	return c
}

func (c Card) Init() tea.Cmd {
	return nil
}

func (c Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		c.setupLists()
		return c, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "esc":
			return c, func() tea.Msg { return project }
		}
	}
	return c, nil
}

func (c Card) Vitew() string {
	if ws.width == 0 {
		return "loading..."
	}
	output := ""
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Center,
		"card",
	)
	return output
}

func (c Card) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	var (
		emptyLine   = "\n"
		inputStyled = ""
		output      = ""
	)
	checklistStyled := ListStyle.Render(c.checklist.View())
	cardlabelsStyled := ListStyle.Render(c.labels.View())
	listsStyled := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistStyled,
		cardlabelsStyled,
	)
	if c.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(c.Input.field.View())
	}
	output = CardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		"Title:",
		c.card.Title,
		emptyLine,
		"Description:",
		c.card.Description,
		emptyLine,
		listsStyled,
		inputStyled,
	))
	return output
}
