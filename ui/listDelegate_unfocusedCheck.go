package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type UnfocusedCheckItemStyles struct {
	NormalTitle       lipgloss.Style
	NormalDoneTitle   lipgloss.Style
	SelectedTitle     lipgloss.Style
	SelectedDoneTitle lipgloss.Style
}

func NewUnfocusedCheckItemStyles() (s UnfocusedCheckItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.NormalDoneTitle = lipgloss.NewStyle().
		Foreground(DoneItemColor).
		Strikethrough(true)
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.SelectedDoneTitle = lipgloss.NewStyle().
		Foreground(DoneItemColor).
		Strikethrough(true)
	return s
}

type UnfocusedCheckListDelegate struct {
	ShowDescription bool
	Styles          UnfocusedCheckItemStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	height          int
	spacing         int
}

func NewUnfocusedCheckListDelegate() UnfocusedCheckListDelegate {
	return UnfocusedCheckListDelegate{
		ShowDescription: true,
		Styles:          NewUnfocusedCheckItemStyles(),
		height:          1,
		spacing:         0,
	}
}

func (d UnfocusedCheckListDelegate) Height() int {
	return d.height
}

func (d UnfocusedCheckListDelegate) Spacing() int {
	return d.spacing
}

func (d UnfocusedCheckListDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, l)
}

func (d UnfocusedCheckListDelegate) Render(w io.Writer, l list.Model, index int, item list.Item) {
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
