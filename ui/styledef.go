package ui

import "github.com/charmbracelet/lipgloss"

var (
	defaultStyle = lipgloss.NewStyle().
			BorderForeground(DefaultBorderColor).
			BorderStyle(lipgloss.RoundedBorder()).
			Align(lipgloss.Center).
			Width(50).
			Padding(1).
			Bold(true)
	emptyStyle = lipgloss.NewStyle().
			BorderForeground(DefaultBorderColor).
			BorderStyle(lipgloss.RoundedBorder()).
			Align(lipgloss.Center).
			Width(50).
			Padding(1).
			Bold(true)
	inputStyle = lipgloss.NewStyle().
			BorderForeground(DefaultBorderColor).
			BorderStyle(lipgloss.RoundedBorder()).
			AlignHorizontal(lipgloss.Left).
			MarginLeft(2).
			Width(100).
			Bold(true)
	menuListStyle = lipgloss.NewStyle().
			BorderForeground(DefaultBorderColor).
			BorderStyle(lipgloss.RoundedBorder()).
			Margin(1, 2).
			Padding(1).
			Bold(true)
	DefaultBorderColor = lipgloss.Color("#fc5603")
)
