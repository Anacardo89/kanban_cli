package ui

import "github.com/charmbracelet/lipgloss"

func EmptyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		BorderForeground(lipgloss.Color("150")).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(1).
		Width(50)

}
