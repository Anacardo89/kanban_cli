package ui

import (
	"fmt"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/Anacardo89/kanban_cli/logger"
	"github.com/Anacardo89/kanban_cli/storage"
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

// Implements tea.Model
type Card struct {
	card      *kanban.Card
	checklist list.Model
	labels    list.Model
	cursor    cursorPos
	textinput textinput.Model
	textarea  textarea.Model
	flag      actionFlag
}

func (c Card) Init() tea.Cmd {
	return nil
}

func (c Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowSize(msg)
		c.setLists()
		return c, nil
	case tea.KeyMsg:
		cmd = c.keyPress(msg)
		return c, cmd
	}
	if c.cursor == checkPos {
		c.checklist, cmd = c.checklist.Update(msg)
	} else if c.cursor == labelPos {
		c.labels, cmd = c.labels.Update(msg)
	}
	return c, cmd
}

func (c Card) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	if c.flag != none {
		return c.viewCertify()
	}
	return c.cardView()

}

// called by model_selector
func OpenCard(kc *kanban.Card) Card {
	c := Card{
		card:      kc,
		textinput: textinput.New(),
		textarea:  textarea.New(),
		cursor:    0,
		flag:      none,
	}
	c.setInput()
	c.setTxtArea()
	c.setLists()
	return c
}

func (c *Card) UpdateCard() {
	c.setInput()
	c.setTxtArea()
	c.setLists()
}

func (c *Card) getCheckItem() *kanban.CheckItem {
	if c.card.CheckList.Length() == 0 {
		return nil
	}
	ci, err := c.card.CheckList.GetAt(c.checklist.Cursor())
	if err != nil {
		logger.Error.Fatal(err)
	}
	return ci.(*kanban.CheckItem)
}

func (c *Card) getCardLabel() *kanban.Label {
	if c.card.CardLabels.Length() == 0 {
		return nil
	}
	l, err := c.card.CardLabels.GetAt(c.labels.Cursor())
	if err != nil {
		logger.Error.Fatal(err)
	}
	return l.(*kanban.Label)
}

// Update
func (c *Card) keyPress(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	if c.textinput.Focused() {
		cmd = c.inputFocused(msg)
		return cmd
	}
	if c.textarea.Focused() {
		switch msg.String() {
		case "esc":
			storage.UpdateCardDesc(c.card.Id, c.textarea.Value())
			c.textarea.Blur()
		}
		c.textarea, cmd = c.textarea.Update(msg)
		return cmd
	}
	if c.flag != none {
		cmd = c.checkFlag(msg)
		return cmd
	}
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "h", "left":
		c.handleMoveLeft()
	case "l", "right":
		c.handleMoveRight()
	case "esc":
		return func() tea.Msg { return upProject }
	case "enter":
		cmd = c.enterKeyPress()
		return cmd
	case "n":
		if c.cursor == checkPos {
			c.textinput.Placeholder = "CheckItem Title"
			return c.textinput.Focus()
		}
		return nil
	case "r":
		switch c.cursor {
		case titlePos:
			c.textinput.Placeholder = "New Card Title"
			return c.textinput.Focus()
		case checkPos:
			c.flag = rename
			c.textinput.Placeholder = "New CheckItem Title"
			return c.textinput.Focus()
		}
		return nil
	case "d":
		switch c.cursor {
		case checkPos:
			c.flag = dCheck
			return nil
		case labelPos:
			c.flag = dLabel
			return nil
		}
	}
	if c.cursor == checkPos {
		c.checklist, cmd = c.checklist.Update(msg)
	} else if c.cursor == labelPos {
		c.labels, cmd = c.labels.Update(msg)
	}
	return cmd
}

func (c *Card) inputFocused(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		c.textinput.SetValue("")
		c.textinput.Blur()
		c.flag = none
		return nil
	case "enter":
		c.txtInputEnter()
	}
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

func (c *Card) txtInputEnter() {
	if c.textinput.Value() == "" {
		return
	}
	switch c.cursor {
	case titlePos:
		storage.UpdateCardTitle(c.card.Id, c.textinput.Value())
		c.card.RenameCard(c.textinput.Value())
	case checkPos:
		if c.flag == rename {
			ci := c.getCheckItem()
			storage.UpdateCheckItemTitle(ci.Id, c.textinput.Value())
			ci.RenameCheckItem(c.textinput.Value())
		} else {
			res := storage.CreateCheckItem(c.textinput.Value(), 0, c.card.Id)
			id64, err := (res.LastInsertId())
			if err != nil {
				logger.Error.Fatal(err)
			}
			id := int(id64)
			c.card.AddCheckItem(id, c.textinput.Value(), false)
		}
	}
	c.setLists()
	c.textinput.SetValue("")
	c.textinput.Blur()
	c.flag = none
}

func (c *Card) enterKeyPress() tea.Cmd {
	switch c.cursor {
	case titlePos:
		c.textinput.Placeholder = "New Card Title"
		c.textinput.Focus()
	case descPos:
		c.textarea.Focus()
	case checkPos:
		if c.card.CheckList.Length() == 0 {
			return nil
		}
		ci := c.getCheckItem()
		done := 0
		if ci.Check {
			done = 1
		}
		storage.UpdateCheckItemDone(ci.Id, done)
		ci.CheckCheckItem()
		c.setLists()
	case labelPos:
		return func() tea.Msg { return labelState }
	}
	return nil
}

func (c *Card) checkFlag(msg tea.KeyMsg) tea.Cmd {
	switch c.flag {
	case dCheck:
		switch msg.String() {
		case "n", "enter", "esc":
			c.flag = none
		case "y":
			ci := c.getCheckItem()
			storage.DeleteCheckItem(ci.Id)
			c.card.RemoveCheckItem(ci)
			c.setLists()
			c.flag = none
		}
		return nil
	case dLabel:
		switch msg.String() {
		case "n", "enter", "esc":
			c.flag = none
		case "y":
			cl := c.getCardLabel()
			storage.DeleteCardLabel(cl.Id)
			c.card.RemoveLabel(cl)
			c.setLists()
			c.flag = none
		}
		return nil
	}
	return nil
}

// actions
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

// View
func (c *Card) cardView() string {
	var inputStyled = ""
	if c.textinput.Focused() {
		inputStyled = InputFieldStyle.Render(c.textinput.View())
	} else {
		switch c.cursor {
		case titlePos:
			inputStyled = InputNoFieldStyle.Render(
				"[hl] [left/right] movement * [ESC] project * [ENTER][R]ename ",
			)
		case descPos:
			inputStyled = InputNoFieldStyle.Render(
				"[hl] [left/right] movement * [ESC] project * [ENTER]/[ESC] focus/unfocus",
			)
		case checkPos:
			inputStyled = InputNoFieldStyle.Render(
				"[hl] [left/right] movement * [ESC] project * [N]ew [D]elete [R]ename * [ENTER] check/uncheck",
			)
		case labelPos:
			inputStyled = InputNoFieldStyle.Render(
				"[hl] [left/right] movement * [ESC] project * [D]elete * [ENTER] select label",
			)
		}
	}
	cardStyled := c.renderCard()
	return lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			cardStyled,
			inputStyled,
		),
	)
}

func (c *Card) viewCertify() string {
	var (
		toDelete interface{}
		areUsure string
	)
	if c.flag == dCheck {
		toDelete = c.getCheckItem()
		areUsure = fmt.Sprintf(
			"Are you sure you wish to delete item\n\n%s\n\nfrom the checklist?\n\ny/N",
			toDelete.(*kanban.CheckItem).Title,
		)
	} else {
		toDelete = c.getCardLabel()
		areUsure = fmt.Sprintf(
			"Are you sure you wish to delete label\n\n%s\n\nfrom the card?\n\ny/N",
			toDelete.(*kanban.Label).Title,
		)
	}

	areUsureStyled := EmptyStyle.Render(areUsure)
	return lipgloss.Place(
		ws.width, ws.height,
		lipgloss.Center, lipgloss.Center,
		areUsureStyled,
	)
}

func (c *Card) renderCard() string {
	emptyLine := ""
	titleStyled := "Title"
	descriptionStyled := "Description"
	txtareaStyled := TextAreaStyle.Render(c.textarea.View())
	checklistStyled := ListStyle.Render(c.checklist.View())
	cardlabelsStyled := ListStyle.Render(c.labels.View())
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
	listsStyled := lipgloss.JoinHorizontal(
		lipgloss.Top,
		checklistStyled,
		cardlabelsStyled,
	)
	return CardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyled,
		c.card.Title,
		emptyLine,
		descriptionStyled,
		txtareaStyled,
		listsStyled,
	))
}

// bubbles
// textinput
func (c *Card) setInput() {
	c.textinput.Prompt = ": "
	c.textinput.CharLimit = 120
	c.textinput.Cursor.Blink = true
}

// textarea
func (c *Card) setTxtArea() {
	c.textarea.Prompt = ""
	c.textarea.Placeholder = "Card Description"
	c.textarea.ShowLineNumbers = true
	c.textarea.Cursor.Blink = true
	c.textarea.SetValue(c.card.Description)
}

// list
var checklistDelegate = NewCheckListDelegate()

func (c *Card) setLists() {
	c.setCheckList()
	c.setCardLabels()
}

func (c *Card) setCheckList() {
	var checklistItems []list.Item
	cl := list.New([]list.Item{}, checklistDelegate, ws.width/2, ws.height/3+1)
	cl.SetShowHelp(false)
	cl.Title = "Checklist"
	cl.InfiniteScrolling = true
	for i := 0; i < c.card.CheckList.Length(); i++ {
		ci, _ := c.card.CheckList.GetAt(i)
		item := Item{
			title: ci.(*kanban.CheckItem).Title,
		}
		if ci.(*kanban.CheckItem).Check {
			item.description = "1"
		} else {
			item.description = "0"
		}
		checklistItems = append(checklistItems, item)
	}
	cl.SetItems(checklistItems)
	c.checklist = cl
}

func (c *Card) setCardLabels() {
	var labelItems []list.Item
	ll := list.New([]list.Item{}, NewLabelListDelegate(), ws.width/2, ws.height/3+1)
	ll.SetShowHelp(false)
	ll.Title = "Card Labels"
	ll.InfiniteScrolling = true
	for i := 0; i < c.card.CardLabels.Length(); i++ {
		l, _ := c.card.CardLabels.GetAt(i)
		item := Item{
			title: l.(*kanban.Label).Title,
		}
		labelItems = append(labelItems, item)
	}
	ll.SetItems(labelItems)
	c.labels = ll
}
