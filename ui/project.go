package ui

import (
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Project struct {
	styles  []lipgloss.Style
	project *kanban.Project
	boards  []list.Model
	labels  list.Model
	cursor  int
	Input   InputField
}

func OpenProject(p *kanban.Project) Project {
	pj := Project{
		styles:  make([]lipgloss.Style, 5),
		cursor:  0,
		Input:   InputField{field: textinput.New()},
		project: p,
	}
	pj.SetStyles()
	return pj
}

func (p Project) Init() tea.Cmd {
	return nil
}

func (p Project) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ws.width = msg.Width
		ws.height = msg.Height
		p.SetStyles()
		return p, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return p, tea.Quit
		case "n":

		}
	}
	return p, nil
}

func (p Project) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	emptyTxtStyled := ""
	bottomLines := ""
	inputStyled := ""
	output := ""
	if p.project.Boards.Length() == 0 {
		emptyTxt := "No Boards.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		emptyTxtStyled = p.styles[empty].Render(emptyTxt)
		if p.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2; i++ {
				bottomLines += "\n"
			}
			inputStyled = p.styles[input].Render(p.Input.field.View())
		}
		output = lipgloss.Place(ws.width, ws.height, lipgloss.Center, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Center, emptyTxtStyled, bottomLines, inputStyled))
		return output
	}
	return output
}
