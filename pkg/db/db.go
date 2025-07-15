package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

var DB *sql.DB

const schema = `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR NOT NULL DEFAULT '',
			comment TEXT NOT NULL DEFAULT '',
			repeat VARCHAR NOT NULL DEFAULT ''
		);
		CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
	`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
