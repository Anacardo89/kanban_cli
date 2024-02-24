package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type CheckItemStyles struct {
	NormalTitle       lipgloss.Style
	NormalDoneTitle   lipgloss.Style
	SelectedTitle     lipgloss.Style
	SelectedDoneTitle lipgloss.Style
}

func NewCheckItemStyles() (s CheckItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.NormalDoneTitle = lipgloss.NewStyle().
		Foreground(DoneItemColor).
		Strikethrough(true)
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(YELLOW)
	s.SelectedDoneTitle = lipgloss.NewStyle().
		Foreground(YELLOW).
		Strikethrough(true)
	return s
}

type CheckListDelegate struct {
	ShowDescription bool
	Styles          CheckItemStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	height          int
	spacing         int
}

func NewCheckListDelegate() CheckListDelegate {
	return CheckListDelegate{
		ShowDescription: true,
		Styles:          NewCheckItemStyles(),
		height:          1,
		spacing:         0,
	}
}

func (d CheckListDelegate) Height() int {
	return d.height
}

func (d CheckListDelegate) Spacing() int {
	return d.spacing
}

func (d CheckListDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, l)
}

func (d CheckListDelegate) Render(w io.Writer, l list.Model, index int, item list.Item) {
	var (
		title, desc string
		s           = &d.Styles
	)

	if i, ok := item.(Item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if l.Width() <= 0 {
		return
	}

	textwidth := uint(l.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, "...")

	isSelected := index == l.Index()
	done := desc == "1"

	if isSelected {
		if done {
			title = s.SelectedDoneTitle.Render(title)
		}
		title = s.SelectedTitle.Render(title)
	} else {
		if done {
			title = s.NormalDoneTitle.Render(title)
		}
		title = s.NormalTitle.Render(title)
	}

	fmt.Fprintf(w, "%s", title)
}
