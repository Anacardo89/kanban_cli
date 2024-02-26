package storage

import (
	"database/sql"

	"github.com/Anacardo89/kanban_cli/logger"
)

// projects
func CreateProject(title string) sql.Result {
	stmt, err := DB.Prepare(CreateProjectSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateProject(id int, title string) sql.Result {
	stmt, err := DB.Prepare(UpdateProjectSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteProject(id int) {
	stmt, err := DB.Prepare(DeleteProjectSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// boards
func CreateBoard(title string, projectId int) sql.Result {
	stmt, err := DB.Prepare(CreateBoardSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(title, projectId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateBoard(id int, title string) sql.Result {
	stmt, err := DB.Prepare(UpdateBoardSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteBoard(id int) {
	stmt, err := DB.Prepare(DeleteBoardSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// labels
func CreateLabel(title string, color string, projectId int) sql.Result {
	stmt, err := DB.Prepare(CreateLabelSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(title, color, projectId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateLabelTitle(id int, title string) sql.Result {
	stmt, err := DB.Prepare(UpdateLabelTitleSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateLabelColor(id int, color string) sql.Result {
	stmt, err := DB.Prepare(UpdateLabelColorSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, color)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteLabel(id int) {
	stmt, err := DB.Prepare(DeleteLabelSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// cards
func CreateCard(title string, boardId int) sql.Result {
	stmt, err := DB.Prepare(CreateCardSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(title, boardId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateCardTitle(id int, title string) sql.Result {
	stmt, err := DB.Prepare(UpdateCardTitleSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateCardDesc(id int, desc string) sql.Result {
	stmt, err := DB.Prepare(UpdateCardDescSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, desc)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateCardParent(id int, boardId int) sql.Result {
	stmt, err := DB.Prepare(UpdateCardParentSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, boardId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteCard(id int) {
	stmt, err := DB.Prepare(DeleteCardSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// card_labels
func CreateCardLabel(cardId int, labelId int) sql.Result {
	stmt, err := DB.Prepare(CreateCardLabelSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(cardId, labelId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteCardLabel(id int) {
	stmt, err := DB.Prepare(DeleteCardLabelSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// check_items
func CreateCheckItem(title string, done int, cardId int) sql.Result {
	stmt, err := DB.Prepare(CreateCheckItemSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(title, done, cardId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateCheckItemTitle(id int, title string) sql.Result {
	stmt, err := DB.Prepare(UpdateCheckItemTitleSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, title)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func UpdateCheckItemDone(id int, done int) sql.Result {
	stmt, err := DB.Prepare(UpdateCheckItemDoneSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	res, err := stmt.Exec(id, done)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
	return res
}

func DeleteCheckItem(id int) {
	stmt, err := DB.Prepare(DeleteCheckItemSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}
