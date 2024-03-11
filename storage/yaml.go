package storage

import (
	"github.com/Anacardo89/kanboards/logger"
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
	Title string `yaml:"title"`
	Color string `yaml:"color"`
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
		logger.Error.Println("Cannot export to YAML", err)
	}
	datastr := string(data)
	return datastr
}

func ToFile(string) {
	// TODO
}
