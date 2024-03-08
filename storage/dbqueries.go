package storage

import (
	"database/sql"

	"github.com/Anacardo89/kanban_cli/logger"
)

// projects
type ProjectSql struct {
	Id    int64
	Title string
}

func GetAllProjects() []ProjectSql {
	var items []ProjectSql
	rows, err := DB.Query(SelectAllProjectsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	for rows.Next() {
		var i ProjectSql
		if err := rows.Scan(
			&i.Id,
			&i.Title,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	rows.Close()
	return items
}

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

func UpdateProject(id int64, title string) sql.Result {
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

func DeleteProject(id int64) {
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
type BoardSql struct {
	Id        int64
	Title     string
	ProjectId int64
}

func GetAllBoards() []BoardSql {
	var items []BoardSql
	rows, err := DB.Query(SelectAllBoardsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	defer rows.Close()
	for rows.Next() {
		var i BoardSql
		if err := rows.Scan(
			&i.Id,
			&i.Title,
			&i.ProjectId,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	return items
}

func CreateBoard(title string, projectId int64) sql.Result {
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

func UpdateBoard(id int64, title string) sql.Result {
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

func DeleteBoard(id int64) {
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
type LabelSql struct {
	Id        int64
	Title     string
	Color     string
	ProjectId int64
}

func GetAllLabels() []LabelSql {
	var items []LabelSql
	rows, err := DB.Query(SelectAllLabelsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	defer rows.Close()
	for rows.Next() {
		var i LabelSql
		if err := rows.Scan(
			&i.Id,
			&i.Title,
			&i.Color,
			&i.ProjectId,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	return items
}

func CreateLabel(title string, color string, projectId int64) sql.Result {
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

func UpdateLabelTitle(id int64, title string) sql.Result {
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

func UpdateLabelColor(id int64, color string) sql.Result {
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

func DeleteLabel(id int64) {
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
type CardSql struct {
	Id      int64
	Title   string
	Desc    sql.NullString
	BoardId int64
}

func GetAllCards() []CardSql {
	var items []CardSql
	rows, err := DB.Query(SelectAllCardsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	defer rows.Close()
	for rows.Next() {
		var i CardSql
		if err := rows.Scan(
			&i.Id,
			&i.Title,
			&i.Desc,
			&i.BoardId,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	return items
}

func CreateCard(title string, boardId int64) sql.Result {
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

func UpdateCardTitle(id int64, title string) sql.Result {
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

func UpdateCardDesc(id int64, desc string) sql.Result {
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

func UpdateCardParent(id int64, boardId int64) sql.Result {
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

func DeleteCard(id int64) {
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
type CardLabelSql struct {
	CardId  int64
	LabelId int64
}

func GetAllCardLabels() []CardLabelSql {
	var items []CardLabelSql
	rows, err := DB.Query(SelectAllCardLabelsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	defer rows.Close()
	for rows.Next() {
		var i CardLabelSql
		if err := rows.Scan(
			&i.CardId,
			&i.LabelId,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	return items
}

func CreateCardLabel(cardId int64, labelId int64) sql.Result {
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

func DeleteCardLabel(labelId int64) {
	stmt, err := DB.Prepare(DeleteCardLabelSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(labelId)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}

// check_items
type CheckItemSql struct {
	Id     int64
	Title  string
	Done   int
	CardId int64
}

func GetAllCheckItems() []CheckItemSql {
	var items []CheckItemSql
	rows, err := DB.Query(SelectAllCheckItemsSql)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt)
	}
	defer rows.Close()
	for rows.Next() {
		var i CheckItemSql
		if err := rows.Scan(
			&i.Id,
			&i.Title,
			&i.Done,
		); err != nil {
			logger.Error.Println(ErrSQLrowScan, err)
		}
		items = append(items, i)
	}
	return items
}

func CreateCheckItem(title string, done int, cardId int64) sql.Result {
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

func UpdateCheckItemTitle(id int64, title string) sql.Result {
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

func UpdateCheckItemDone(id int64, done int) sql.Result {
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

func DeleteCheckItem(id int64) {
	stmt, err := DB.Prepare(DeleteCheckItemSql)
	if err != nil {
		logger.Error.Fatal(ErrCreatSQLstmt, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatal(ErrExecSQLstmt, err)
	}
}
