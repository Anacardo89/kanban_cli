package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type LabelItemStyles struct {
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
}

func NewLabelItemStyles() (s LabelItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(BLACK)
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(YELLOW)
	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(BLACK)
	return s
}

type LabelListDelegate struct {
	ShowDescription bool
	Styles          LabelItemStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	height          int
	spacing         int
}

func NewLabelListDelegate() LabelListDelegate {
	return LabelListDelegate{
		ShowDescription: true,
		Styles:          NewLabelItemStyles(),
		height:          2,
		spacing:         0,
	}
}

func (d LabelListDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

func (d LabelListDelegate) Spacing() int {
	return d.spacing
}

func (d LabelListDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, l)
}

func (d LabelListDelegate) Render(w io.Writer, l list.Model, index int, item list.Item) {
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
	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= d.height-1 {
			break
		}
		lines = append(lines, truncate.StringWithTail(line, textwidth, "..."))
	}
	desc = strings.Join(lines, "\n")

	isSelected := index == l.Index()

	if isSelected {
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Background(lipgloss.Color(desc)).Render(desc)
	} else {
		title = s.NormalTitle.Render(title)
	}
	desc = s.NormalDesc.Background(lipgloss.Color(desc)).Render(desc)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
