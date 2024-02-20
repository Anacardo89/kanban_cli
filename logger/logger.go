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

func CreateLogger() {
	file, err := os.OpenFile(loggerPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Info = log.New(file, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
