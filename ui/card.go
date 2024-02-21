package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Card struct {
	card      *kanban.Card
	checklist list.Model
	labels    list.Model
	cursor    int
	icursor   int
	Input     InputField
	TxtArea   TextArea
	inputFlag inputFlag
}

func OpenCard(kc *kanban.Card) Card {
	c := Card{
		card:    kc,
		Input:   InputField{field: textinput.New()},
		TxtArea: TextArea{field: textarea.New()},
		cursor:  0,
	}
	c.setupLists()
	return c
}

func (c Card) Init() tea.Cmd {
	return nil
}

func (c Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	c.setupLists()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		return c, nil
	case tea.KeyMsg:
		if c.Input.field.Focused() {
			c.handleInput(msg.String())
			c.Input.field, cmd = c.Input.field.Update(msg)
			return c, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "esc":
			return c, func() tea.Msg { return project }
		case "right":
			c.handleMoveRight()
		case "left":
			c.handleMoveLeft()
		case "enter":
			c.setInput()
			c.handleEnter()
		}
	}
	if c.cursor == 1 {
		c.checklist, cmd = c.checklist.Update(msg)
	} else if c.cursor == 2 {
		c.labels, cmd = c.labels.Update(msg)
	}
	return c, cmd
}

func (c Card) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	var (
		cardStyled        = ""
		titleStyled       = "Title"
		descriptionStyled = "Description"
		txtareaStyled     = ""
		emptyLine         = ""
		bottomLines       = ""
		inputStyled       = ""
		output            = ""
	)
	checklistStyled := ListStyle.Render(c.checklist.View())
	cardlabelsStyled := ListStyle.Render(c.labels.View())
	switch c.cursor {
	case 0:
		titleStyled = SelectedTxtStyle.Render(titleStyled)
	case 1:
		descriptionStyled = SelectedTxtStyle.Render(descriptionStyled)
	case 2:
		checklistStyled = SelectedListStyle.Render(c.checklist.View())
	case 3:
		cardlabelsStyled = SelectedListStyle.Render(c.labels.View())
	}
	listsStyled := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistStyled,
		cardlabelsStyled,
	)
	if c.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(c.Input.field.View())
	}
	if c.TxtArea.field.Focused() {
		txtareaStyled = InputFieldStyle.Render(c.TxtArea.field.View())
	}
	cardStyled = CardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyled,
		c.card.Title,
		emptyLine,
		descriptionStyled,
		txtareaStyled,
		c.card.Description,
		emptyLine,
		listsStyled,
	))
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			cardStyled,
			bottomLines,
			inputStyled,
		))
	return output
}

// movement
func (c *Card) handleMoveRight() {
	if c.cursor == 3 {
		c.cursor = 0
	} else {
		c.cursor++
	}
}

func (c *Card) handleMoveLeft() {
	if c.cursor == 0 {
		c.cursor = 3
	} else {
		c.cursor--
	}
}

// action
func (c *Card) handleEnter() {
	switch c.cursor {
	case 0:
		c.Input.field.Focus()
	case 1:
		c.TxtArea.field.Focus()
	case 2:

	}
}
