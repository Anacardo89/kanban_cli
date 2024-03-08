package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type UnfocusedCardStyles struct {
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
}

func NewUnfocusedCardStyles() (s UnfocusedCardStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.NormalDesc = lipgloss.NewStyle()
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.SelectedDesc = lipgloss.NewStyle()
	return s
}

type UnfocusedCardDelegate struct {
	ShowDescription bool
	Styles          UnfocusedCardStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	height          int
	spacing         int
}

func NewUnfocusedCardDelegate() UnfocusedCardDelegate {
	return UnfocusedCardDelegate{
		ShowDescription: true,
		Styles:          NewUnfocusedCardStyles(),
		height:          2,
		spacing:         0,
	}
}

func (d UnfocusedCardDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

func (d UnfocusedCardDelegate) Spacing() int {
	return d.spacing
}

func (d UnfocusedCardDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, l)
}

func (d UnfocusedCardDelegate) Render(w io.Writer, l list.Model, index int, item list.Item) {
	var (
		title, desc string
		meta        []Meta
		s           = &d.Styles
	)

	if i, ok := item.(Item); ok {
		title = i.Title()
		desc = i.Description()
		meta = i.meta
	} else {
		return
	}

	if l.Width() <= 0 {
		return
	}

	textwidth := uint(l.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, "...")
	isSelected := index == l.Index()
	if isSelected {
		title = s.SelectedTitle.Render(title)
	} else {
		title = s.NormalTitle.Render(title)
	}

	desc = formatDescription(desc, meta)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
