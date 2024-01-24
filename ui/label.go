package ui

import (
	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Label struct {
	styles    []lipgloss.Style
	project   *kanban.Project
	labels    list.Model
	cursor    int
	sl        *dll.Node
	Input     InputField
	inputFlag inputFlag
}

func OpenLabel(kp *kanban.Project) Label {
	l := Label{
		styles:  make([]lipgloss.Style, 10),
		project: kp,
		Input:   InputField{field: textinput.New()},
	}
	return l
}

func (l Label) Init() tea.Cmd {
	return nil
}

func (l Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ws.width = msg.Width
		ws.height = msg.Height
		return l, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return l, tea.Quit
		case "b", "esc":
			return l, func() tea.Msg { return project }
		}
	}
	return l, nil
}

func (l Label) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	output := ""
	output = lipgloss.Place(
		ws.width,
		ws.height,
		lipgloss.Center,
		lipgloss.Center,
		"label",
	)
	return output
}
