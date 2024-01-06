package ui

import "github.com/charmbracelet/lipgloss"

var DefaultBorderColor = lipgloss.Color("#fc5603")

func DefaultStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Width(50).
		Padding(1).
		Bold(true)
}

func InputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		AlignHorizontal(lipgloss.Left).
		MarginLeft(2).
		Width(100).
		Bold(true)
}

func EmptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Width(50).
		Padding(1).
		Bold(true)
}
