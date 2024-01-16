package ui

import "github.com/charmbracelet/lipgloss"

type styledef int

const (
	defaultStyle styledef = iota
	empty
	input
	listStyle
)

var (
	DefaultBorderColor = lipgloss.Color("#fc5603")
)

func (m Menu) SetStyles() {
	m.styles[defaultStyle] = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Width(m.witdh / 2).
		Padding(1).
		Bold(true)
	m.styles[empty] = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		MarginTop(1).
		Padding(1, 3).
		Bold(true)
	m.styles[input] = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		AlignHorizontal(lipgloss.Left).
		Width(m.witdh - 2).
		Bold(true)
	m.styles[listStyle] = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Margin(1, 1, 0).
		Padding(1).
		Bold(true)
}
