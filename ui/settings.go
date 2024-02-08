package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Window
type WindowSize struct {
	width  int
	height int
}

var ws WindowSize

func updateWindowSize(msg tea.WindowSizeMsg) {
	ws.width = msg.Width
	ws.height = msg.Height
	updateStyles()
}

// Colors
var (
	BLACK                 = lipgloss.Color("#000000")
	WHITE                 = lipgloss.Color("#ffffff")
	DefaultBorderColor    = lipgloss.Color("#fc5603")
	ListItemColor         = lipgloss.Color("#42daf5")
	SelectedListItemColor = lipgloss.Color("#e9f542")
)

// Styles
var (
	DefaultStyle             lipgloss.Style
	EmptyStyle               lipgloss.Style
	InputFieldStyle          lipgloss.Style
	ListStyle                lipgloss.Style
	ProjectListStyle         lipgloss.Style
	ProjectListSelectedStyle lipgloss.Style
	ProjectTitle             lipgloss.Style
)

func updateStyles() {
	DefaultStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Width(ws.width / 2).
		Padding(1).
		Bold(true)

	EmptyStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		MarginTop(1).
		Padding(1, 3).
		Bold(true)

	InputFieldStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		AlignHorizontal(lipgloss.Left).
		Width(ws.width - 2).
		Bold(true)

	ListStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Margin(1, 1, 0).
		Padding(1).
		Bold(true)

	ProjectListStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Margin(0, 1, 1).
		Padding(1).
		Bold(true)

	ProjectListSelectedStyle = ProjectListStyle.Copy().
		BorderForeground(WHITE)

	ProjectTitle = lipgloss.NewStyle().
		Align(lipgloss.Center).
		MarginLeft(2).
		Bold(true)
}
