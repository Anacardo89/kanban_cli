package fsops

import (
	"log"
	"os"
)

var (
	LoggerPath string
	DBPath     string
	YamlPath   string
)

func Home() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot get HOME:", err)
	}
	return home
}

func CreateDir() {
	err := os.Mkdir(Home()+"/kanboards", 0755)
	if err != nil {
		log.Fatal("Cannot create kanboards directory:", err)
	}
}

func SetPaths() {
	LoggerPath = Home() + "/kanboards/kb.log"
	DBPath = Home() + "/kanboards/kb.db"
	YamlPath = Home() + "/kanboards/kb.yaml"
}
