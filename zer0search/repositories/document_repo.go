package repositories

import (
	"database/sql"
	"fmt"
	"log/slog"
	"searchengine/db"
)

type DocumentRepo struct {
	db *sql.DB
}

func NewDocumentRepo(db *sql.DB) *DocumentRepo {
	return &DocumentRepo{
		db: db,
	}
}

func (d *DocumentRepo) Insert(document string) (int64, error) {
	res, err1 := d.db.Exec(db.InsertStmt, document)
	docId, err2 := res.LastInsertId()
	if err1 == nil && err2 == nil {
		return docId, nil
	}
	slog.Info("[document_repo.go] [Insert()] document insertion error : ", err1, err2)
	return 0, fmt.Errorf("%w, %w", err1, err2)
}

func (d *DocumentRepo) Query(id int) (string, error) {
	var document string
	if err := d.db.QueryRow(db.QueryStmt, id).Scan(&document); err != nil {
		slog.Error("[document_repo.go] [Query()] document retriving error : ", err)
		return "", err
	}
	return document, nil
}

func (d *DocumentRepo) DeleteAt(docId int) {

}
