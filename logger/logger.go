package logger

import (
	"log"
	"os"
)

var (
	Info       *log.Logger
	Error      *log.Logger
	loggerPath = "./log.log"
)

func CreateLogger() (*os.File, error) {
	file, err := os.OpenFile(loggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("Cannot access log file:", err)
	}
	Info = log.New(file, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
	return file, nil
}
