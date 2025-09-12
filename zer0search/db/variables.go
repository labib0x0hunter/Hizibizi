package db

var (
	user        = "root"
	password    = ""
	ip          = "127.0.0.1"
	port        = "3306"
	dbName      = "testDb"
	tableName   = "items"
	InsertStmt  = "INSERT INTO items (content) VALUES (?)"
	QueryStmt   = "SELECT content FROM items WHERE id = ?"
	deleteTable = "DELETE FROM "
	resetTable  = "ALTER TABLE " + "items" + " AUTO_INCREMENT = 1"

	NoEntryError    = "sql: Scan error on column index 0"
	InsertTempError = "<nil>, <nil>"
)
