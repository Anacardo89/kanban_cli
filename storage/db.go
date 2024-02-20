package storage

import (
	"database/sql"

	"github.com/Anacardo89/kanban_cli/logger"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	stmt *sql.Stmt
)

func SetDB() {
	var err error
	db, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		logger.Error.Println(err)
		err = nil
	}
}

func FillDB() {
	var err error
	stmt, err = db.Prepare(CreateDB)
	if err != nil {
		logger.Error.Println(err)
		err = nil
	}
	stmt.Exec()
}
