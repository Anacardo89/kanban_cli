package storage

import (
	"database/sql"
	"os"

	"github.com/Anacardo89/kanban_cli/logger"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB     *sql.DB
	dbPath string
)

func DBExists() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error.Println("Cannot extract HOME:", err)
	}
	dbPath = home + "/.kanboards/db.db"
	_, err = os.Open(dbPath)
	if err != nil {
		return false
	}
	return true
}

func SetDB() {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error.Println(err)
		err = nil
	}
}

func CreateDB() {
	file, err := os.OpenFile(dbPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.Error.Fatal("Cannot create DB file:", err)
	}
	defer file.Close()
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error.Fatal("Cannot establish DB connection:", err)
	}
	createDB, err := DB.Prepare(CreateDBsql)
	if err != nil {
		logger.Error.Fatal("Error creating SQL statement:", err)
	}
	createDB.Exec()
}

func CreateDBfile() {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error.Println("Cannot extract HOME:", err)
	}
	os.Mkdir(home+"/.kanboards", 0755)
}
