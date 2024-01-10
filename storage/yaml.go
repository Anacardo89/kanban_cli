package storage

import (
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Title  string  `yaml:"title"`
	Lists  []List  `yaml:"lists"`
	Labels []Label `yaml:"labels"`
}

type List struct {
	Title string `yaml:"title"`
	Cards []Card `yaml:"cards"`
}

type Label struct {
	Title string         `yaml:"title"`
	Color lipgloss.Color `yaml:"color"`
}

type Card struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	CheckList   []CheckItem `yaml:"checklist"`
	CardLabels  []Label     `yaml:"card_labels"`
}

type CheckItem struct {
	Title string `yaml:"title"`
	Check bool   `yaml:"check"`
}
