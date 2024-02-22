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
	RED                   = lipgloss.Color("#ba3525")
	BLUE                  = lipgloss.Color("#77ccee")
	YELLOW                = lipgloss.Color("#fecc00")
	GREEN                 = lipgloss.Color("#0edd1e")
	DefaultBorderColor    = lipgloss.Color("#fc5603")
	SelectedListItemColor = lipgloss.Color("#e9f542")
	DoneItemColor         = lipgloss.Color("#0ff702")
)

// Styles
var (
	DefaultStyle             lipgloss.Style
	EmptyStyle               lipgloss.Style
	SelectedTxtStyle         lipgloss.Style
	InputFieldStyle          lipgloss.Style
	TextAreaStyle            lipgloss.Style
	ListStyle                lipgloss.Style
	SelectedListStyle        lipgloss.Style
	DoneItemStyle            lipgloss.Style
	ProjectListStyle         lipgloss.Style
	ProjectListSelectedStyle lipgloss.Style
	ProjectTitle             lipgloss.Style
	CardStyle                lipgloss.Style
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

	SelectedTxtStyle = lipgloss.NewStyle().
		Foreground(SelectedListItemColor)

	InputFieldStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		AlignHorizontal(lipgloss.Left).
		Width(ws.width - 2).
		Bold(true)

	TextAreaStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderLeftForeground(DefaultBorderColor).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(2).
		Width(ws.width - 4).
		Height(ws.height / 4).
		Bold(true)

	ListStyle = lipgloss.NewStyle().
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder()).
		Margin(0, 1).
		Padding(0).
		Bold(true)

	SelectedListStyle = lipgloss.NewStyle().
		BorderForeground(WHITE).
		BorderStyle(lipgloss.RoundedBorder()).
		Margin(0, 1).
		Padding(0).
		Bold(true)

	DoneItemStyle = lipgloss.NewStyle().
		Strikethrough(true)

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

	CardStyle = lipgloss.NewStyle().
		Width(ws.width - 2).
		BorderForeground(DefaultBorderColor).
		BorderStyle(lipgloss.RoundedBorder())
}
