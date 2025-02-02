package database

import (
	"database/sql"
	"time"

	"OpnLaaS.cyber.unh.edu/lib"
	_ "modernc.org/sqlite"
)

const (
	maxRetries = 5
	baseDelay  = 100 * time.Millisecond
)

var db *sql.DB

func open() (*sql.DB, error) {
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("sqlite", lib.Config.DBFile)
		if err == nil {
			break
		}

		time.Sleep(baseDelay * time.Duration(i))
	}

	return db, err
}

func Connect() bool {
	lib.Log.Basic("Connecting to database...")

	var err error

	db, err = open()

	if err != nil {
		lib.Log.Error("Could not connect to database: " + err.Error())
		return false
	}

	lib.Log.Basic("Checking tables...")

	if _, err = db.Exec(USERS_STATEMENT); err != nil {
		lib.Log.Error("Could not create users table: " + err.Error())
		return false
	}

	lib.Log.Success("Database is ready")

	return true
}
