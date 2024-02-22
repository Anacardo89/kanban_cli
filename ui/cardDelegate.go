package ui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type CardStyles struct {
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style
}

func NewCardStyles() (s CardStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(WHITE)
	s.NormalDesc = lipgloss.NewStyle()
	s.SelectedTitle = lipgloss.NewStyle().
		Foreground(SelectedListItemColor)
	s.SelectedDesc = lipgloss.NewStyle()
	return s
}

type CardDelegate struct {
	ShowDescription bool
	Styles          CardStyles
	UpdateFunc      func(tea.Msg, *list.Model) tea.Cmd
	height          int
	spacing         int
}

func NewCardDelegate() CardDelegate {
	return CardDelegate{
		ShowDescription: true,
		Styles:          NewCardStyles(),
		height:          2,
		spacing:         0,
	}
}

func (d CardDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

func (d CardDelegate) Spacing() int {
	return d.spacing
}

func (d CardDelegate) Update(msg tea.Msg, l *list.Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, l)
}

func (d CardDelegate) Render(w io.Writer, l list.Model, index int, item list.Item) {
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
	if isSelected {
		title = s.SelectedTitle.Render(title)
	} else {
		title = s.NormalTitle.Render(title)
	}

	desc = formatDescription(desc)

	fmt.Fprintf(w, "%s\n%s", title, desc)
}

func formatDescription(desc string) string {
	descSlice := strings.Split(desc, " ")
	if len(descSlice[2]) == 0 {
		return desc
	}

	labels := strings.Split(descSlice[2], ",")
	var labelID []string
	var labelColor []string
	for _, label := range labels {
		splitLabel := strings.Split(label, "#")
		labelID = append(labelID, splitLabel[0])
		labelColor = append(labelColor, splitLabel[1])
	}
	if len(labelID) != len(labelColor) {
		log.Println("Labels and colors don't match")
	}

	labelout := ""
	for i := 0; i < len(labelID); i++ {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(labelColor[i]))
		styled := style.Render(labelID[i])
		labelout += styled
	}

	return descSlice[0] + " " + descSlice[1] + " " + labelout

}
