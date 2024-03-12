package storage

import (
	"database/sql"
	"os"

	"github.com/Anacardo89/kanboards/logger"
	_ "github.com/mattn/go-sqlite3"
)

const (
	ErrCreatSQLstmt string = "Error creating SQL statement:"
	ErrExecSQLstmt  string = "Error executing SQL statement:"
	ErrSQLrowScan   string = "Error scanning rows:"
)

var (
	home     string
	dbPath   string
	yamlPath string
	DB       *sql.DB
)

func SetPaths() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		logger.Error.Fatal("Cannot get HOME:", err)
	}
	dbPath = home + "/kanboards/db.db"
	yamlPath = home + "/kanboards/kb.yaml"
}

func DBExists() bool {
	_, err := os.Open(dbPath)
	return err == nil
}

func OpenDB() {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Error.Fatal("Cannot establish DB connection:", err)
		err = nil
	}
}

func CreateDBfile() {
	err := os.Mkdir(home+"/kanboards", 0755)
	if err != nil {
		logger.Error.Println("Cannot create working directory:", err)
	}
}

func CreateDBTables() {
	file, err := os.OpenFile(dbPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.Error.Fatal("Cannot open DB file:", err)
	}
	defer file.Close()

	OpenDB()

	createProjects, err := DB.Prepare(CreateTableProjects)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createProjects.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createBoards, err := DB.Prepare(CreateTableBoards)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createBoards.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createLabels, err := DB.Prepare(CreateTableLabels)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createLabels.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCards, err := DB.Prepare(CreateTableCards)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCards.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCardLabels, err := DB.Prepare(CreateTableCardLabels)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCardLabels.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}

	createCheckItems, err := DB.Prepare(CreateTableCheckItems)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = createCheckItems.Exec()
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}
