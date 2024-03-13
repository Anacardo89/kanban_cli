package logger

import (
	"log"
	"os"

	"github.com/Anacardo89/kanboards/fsops"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func CreateLogger() (*os.File, error) {
	if _, err := os.Open(fsops.LoggerPath); err != nil {
		fsops.CreateDir()
	}
	file, err := os.OpenFile(fsops.LoggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access log file:", err)
	}
	Info = log.New(file, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	return file, nil
}
