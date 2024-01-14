package storage

import (
	"log"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

type Menu struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Title  string  `yaml:"title"`
	Boards []Board `yaml:"lists"`
	Labels []Label `yaml:"labels"`
}

type Board struct {
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

func (m *Menu) ToYAML() string {
	data, err := yaml.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	datastr := string(data)
	return datastr
}

func ToFile(string) {

}
