package storage

import (
	"os"

	"github.com/Anacardo89/kanboards/logger"
	"gopkg.in/yaml.v2"
)

type Menu struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Id     int64   `yaml:"id"`
	Title  string  `yaml:"title"`
	Boards []Board `yaml:"boards"`
	Labels []Label `yaml:"labels"`
}

type Board struct {
	Id    int64  `yaml:"id"`
	Pos   int    `yaml:"position"`
	Title string `yaml:"title"`
	Cards []Card `yaml:"cards"`
}

type Label struct {
	Id    int64  `yaml:"id"`
	Title string `yaml:"title"`
	Color string `yaml:"color"`
}

type Card struct {
	Id          int64       `yaml:"id"`
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	CheckList   []CheckItem `yaml:"checklist"`
	CardLabels  []Label     `yaml:"card_labels"`
}

type CheckItem struct {
	Id    int64  `yaml:"id"`
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

func FromYAML(data []byte) *Menu {
	m := Menu{}
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		logger.Error.Println("Cannot import from YAML", err)
	}
	return &m
}

func ToFile(data string) {
	yamlPath := "./kb.yaml"
	f, err := os.Create(yamlPath)
	if err != nil {
		logger.Error.Println(err)
	}
	defer f.Close()
	f.WriteString(data)
}

func FromFile() []byte {
	yamlPath := "./kb.yaml"
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		logger.Error.Println(err)
	}
	return data
}
