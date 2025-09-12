package db

import (
	"testing"
)

func TestInsertQuery(t *testing.T) {

	// newDb, err := NewDbService()
	// if err != nil {
	// 	t.Errorf("[error] %v", err)
	// }
	// defer newDb.Close()

	// newDb.db.Exec(`TRUNCATE TABLE items`)

	// lastInsert, err := newDb.GetDocId()
	// if err != nil && !strings.HasPrefix(err.Error(), NoEntryError) {
	// 	t.Errorf("[error-01] %v", err)
	// } else if strings.HasPrefix(err.Error(), NoEntryError) {
	// 	if err := newDb.Insert("word"); err != nil {
	// 		t.Errorf("Insert(%s) : %v want <nil>", "word", err)
	// 	}
	// 	lastInsert, err = newDb.GetDocId()
	// }
	// if err != nil {
	// 	t.Errorf("[error-02] %v", err)
	// }
	// if lastInsert != 1 {
	// 	t.Errorf("GetDocId() = %v want 1", lastInsert)
	// }

	// word, err := newDb.Query(1)
	// if err != nil {
	// 	t.Errorf("[error-02] %v", err)
	// }

	// if word != "word" {
	// 	t.Errorf("Query(%d) = %v want %s", lastInsert, word, "word")
	// }
}
