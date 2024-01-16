package ui

import "github.com/charmbracelet/bubbles/list"

var (
	menuItems []list.Item
)

type menuItem struct {
	title       string
	description string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

func (m *Menu) SetupList() {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), m.witdh/3, m.height-9)
	l.SetShowHelp(false)
	l.Title = "Projects"
	l.InfiniteScrolling = true
	m.list = l
}
