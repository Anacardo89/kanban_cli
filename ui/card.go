package ui

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cursorPos int

const (
	titlePos cursorPos = iota
	descPos
	checkPos
	labelPos
)

type Card struct {
	card      *kanban.Card
	checklist list.Model
	labels    list.Model
	cursor    cursorPos
	icursor   int
	Input     InputField
	textarea  textarea.Model
	inputFlag inputFlag
}

func OpenCard(kc *kanban.Card) Card {
	c := Card{
		card:     kc,
		Input:    InputField{field: textinput.New()},
		textarea: textarea.New(),
		cursor:   0,
	}
	c.setupLists()
	c.setTxtArea()
	return c
}

func (c *Card) UpdateCard() {
	c.setupLists()
	c.setTxtArea()
}

func (c Card) Init() tea.Cmd {
	return nil
}

func (c Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		c.setupLists()
		c.setTxtArea()
		return c, nil
	case tea.KeyMsg:
		if c.Input.field.Focused() {
			c.handleInput(msg.String())
			c.Input.field, cmd = c.Input.field.Update(msg)
			return c, cmd
		}
		if c.textarea.Focused() {
			c.handleTextArea(msg.String())
			c.textarea, cmd = c.textarea.Update(msg)
			return c, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return c, tea.Quit
		case "esc":
			return c, func() tea.Msg { return upProject }
		case "right":
			c.handleMoveRight()
		case "left":
			c.handleMoveLeft()
		case "n":
			if c.cursor == checkPos {
				c.inputFlag = new
				c.setInput()
				c.Input.field.Focus()
			}
			return c, nil
		case "d":
			if c.cursor == checkPos || c.cursor == labelPos {
				c.handleDelete()
			}
		case "enter":
			c.setInput()
			switch c.cursor {
			case titlePos:
				c.Input.field.Focus()
			case descPos:
				c.textarea.Focus()
			case checkPos:
				checkitem := c.getCheckItem()
				checkitem.CheckCheckItem()
				c.setupLists()
			case labelPos:
				return c, func() tea.Msg { return label }
			}
		}
	}

	if c.cursor == checkPos {
		c.checklist, cmd = c.checklist.Update(msg)
		c.icursor = c.checklist.Cursor()
	} else if c.cursor == labelPos {
		c.labels, cmd = c.labels.Update(msg)
		c.icursor = c.labels.Cursor()
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
		inputStyled       = ""
		output            = ""
	)

	if c.Input.field.Focused() {
		inputStyled = InputFieldStyle.Render(c.Input.field.View())
	}
	txtareaStyled = TextAreaStyle.Render(c.textarea.View())
	checklistStyled := ListStyle.Render(c.checklist.View())
	cardlabelsStyled := ListStyle.Render(c.labels.View())

	// highlight selected
	switch c.cursor {
	case titlePos:
		titleStyled = SelectedTxtStyle.Render(titleStyled)
	case descPos:
		descriptionStyled = SelectedTxtStyle.Render(descriptionStyled)
	case checkPos:
		checklistStyled = SelectedListStyle.Render(c.checklist.View())
	case labelPos:
		cardlabelsStyled = SelectedListStyle.Render(c.labels.View())
	}

	// build output
	listsStyled := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistStyled,
		cardlabelsStyled,
	)

	cardStyled = CardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyled,
		c.card.Title,
		emptyLine,
		descriptionStyled,
		txtareaStyled,
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
			inputStyled,
		))
	return output
}

// movement
func (c *Card) handleMoveRight() {
	if c.cursor == labelPos {
		c.cursor = titlePos
	} else {
		c.cursor++
	}
}

func (c *Card) handleMoveLeft() {
	if c.cursor == titlePos {
		c.cursor = labelPos
	} else {
		c.cursor--
	}
}

// action
func (c *Card) getCheckItem() *kanban.CheckItem {
	if c.card.CheckList.Length() == 0 {
		return nil
	}
	c.icursor = c.checklist.Cursor()
	node, err := c.card.CheckList.WalkTo(c.icursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	return node.Val().(*kanban.CheckItem)
}

func (c *Card) getCardLabel() *kanban.Label {
	if c.card.CardLabels.Length() == 0 {
		return nil
	}
	c.icursor = c.labels.Cursor()
	node, err := c.card.CardLabels.WalkTo(c.icursor)
	if err != nil {
		log.Println(err)
		return nil
	}
	return node.Val().(*kanban.Label)
}

func (c *Card) handleDelete() {
	switch c.cursor {
	case checkPos:
		node, err := c.card.CheckList.WalkTo(c.checklist.Cursor())
		if err != nil {
			log.Println(err)
		}
		checkitem := node.Val().(*kanban.CheckItem)
		c.card.RemoveCheckItem(checkitem)
	case labelPos:
		node, err := c.card.CardLabels.WalkTo(c.labels.Cursor())
		if err != nil {
			log.Println(err)
		}
		cardlabel := node.Val().(*kanban.Label)
		c.card.RemoveLabel(cardlabel)
	}
	c.setupLists()
}
