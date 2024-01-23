package ui

import (
	"log"

	"github.com/Anacardo89/ds/lists/dll"
	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	styles   []lipgloss.Style
	menu     *kanban.Menu
	cursor   int
	selected *dll.Node
	list     list.Model
	Input    InputField
}

func NewMenu() Menu {
	return Menu{
		styles: make([]lipgloss.Style, 5),
		cursor: 0,
		menu:   kanban.StartMenu(),
		Input:  InputField{field: textinput.New()},
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ws.width = msg.Width
		ws.height = msg.Height
		m.SetStyles()
		m.SetupList()
		return m, nil
	case tea.KeyMsg:
		if m.Input.field.Focused() {
			m.handleInput(msg.String())
			m.Input.field, cmd = m.Input.field.Update(msg)
			return m, cmd
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.handleMoveUp()
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		case "down":
			m.handleMoveDown()
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		case "enter":
			p := m.getProject()
			if p == nil {
				return m, nil
			}
			return m, func() tea.Msg { return project }
		case "n":
			m.setInput()
			return m, m.Input.field.Focus()
		}

	}
	return m, nil
}

func (m Menu) View() string {
	if ws.width == 0 {
		return "loading..."
	}
	var (
		bottomLines    = ""
		emptyTxtStyled = ""
		inputStyled    = ""
		menustyled     = ""
		output         = ""
	)

	if m.menu.Projects.Length() == 0 {
		emptyTxt := "No projects.\n\nPress 'n' to create a new Project Board\nor 'q' to quit"
		emptyTxtStyled = m.styles[empty].Render(emptyTxt)
		if m.Input.field.Focused() {
			_, h := lipgloss.Size(emptyTxtStyled)
			for i := 0; i < ws.height-h-h/2; i++ {
				bottomLines += "\n"
			}
			inputStyled = m.styles[input].Render(m.Input.field.View())
		}
		output = lipgloss.Place(ws.width, ws.height, lipgloss.Center, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Center, emptyTxtStyled, bottomLines, inputStyled))
		return output
	}

	menustyled = m.styles[listStyle].Render(m.list.View())
	if m.Input.field.Focused() {
		inputStyled = m.styles[input].Render(m.Input.field.View())
	}
	output = lipgloss.Place(ws.width, ws.height, lipgloss.Left, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, menustyled, bottomLines, inputStyled))
	return output
}

func (m *Menu) handleMoveUp() {
	if m.menu.Projects.Length() == 0 {
		return
	}
	var err error
	var node *dll.Node
	if m.cursor == 0 {
		m.cursor = m.menu.Projects.Length() - 1
		node, err = m.menu.Projects.TailNode()
		if err != nil {
			log.Println(err)
		}
		m.selected = node
		return
	}
	m.cursor--
	m.selected, err = m.selected.Prev()
	if err != nil {
		log.Println(err)
	}
}

func (m *Menu) handleMoveDown() {
	if m.menu.Projects.Length() == 0 {
		return
	}
	var err error
	var node *dll.Node
	if m.cursor == m.menu.Projects.Length()-1 {
		m.cursor = 0
		node, err = m.menu.Projects.HeadNode()
		if err != nil {
			log.Println(err)
		}
		m.selected = node
		return
	}
	m.cursor++
	m.selected, err = m.selected.Prev()
	if err != nil {
		log.Println(err)
	}
}

func (m *Menu) getProject() *kanban.Project {
	if m.menu.Projects.Length() == 0 {
		return nil
	}
	node, err := m.menu.Projects.WalkTo(m.cursor)
	if err != nil {
		log.Println(err)
	}
	return node.Val().(*kanban.Project)
}

func (m *Menu) setInput() {
	m.Input.field.Prompt = ": "
	m.Input.field.CharLimit = 120
	m.Input.field.Placeholder = "Project Title"
}

func (m *Menu) handleInput(key string) {
	switch key {
	case "esc":
		m.Input.field.SetValue("")
		m.Input.data = ""
		m.Input.field.Blur()
		return
	case "enter":
		log.Println(m.Input.field.Value())
		m.Input.data = m.Input.field.Value()
		menuItem := menuItem{
			title: m.Input.data,
		}
		menuItems = append(menuItems, menuItem)
		m.list.SetItems(menuItems)
		m.menu.AddProject(m.Input.data)
		m.Input.data = ""
		m.Input.field.Blur()
		m.cursor = 0
		node, err := m.menu.Projects.HeadNode()
		if err != nil {
			log.Println(err)
		}
		m.selected = node
		return
	}
}
