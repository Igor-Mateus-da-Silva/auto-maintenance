package database

import (
	"database/sql"
	"errors"
)

// DB holds the database connection pool.
// SQLite configuration will be fully implemented in Task 2.
var DB *sql.DB

var ErrNotImplemented = errors.New("database driver configuration scheduled for task 2")

func Init(dbPath string) error {
	// Placeholder ready for SQLite connection setup in Task 2
	return nil
}
