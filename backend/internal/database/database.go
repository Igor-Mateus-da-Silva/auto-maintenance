package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Connect opens a SQLite connection, ensures directory exists, and sets PRAGMAs.
func Connect(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	// PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;
	if _, err := db.Exec(`PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to configure sqlite pragmas: %w", err)
	}

	return db, nil
}

// RunMigrations creates the required tables if they do not exist.
func RunMigrations(db *sql.DB) error {
	queries := `
	CREATE TABLE IF NOT EXISTS vehicles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		plate TEXT,
		year INTEGER,
		current_mileage INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS maintenance_types (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		expiration_type TEXT NOT NULL CHECK(expiration_type IN ('time', 'km', 'both')),
		interval_days INTEGER,
		interval_km INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vehicle_maintenances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vehicle_id INTEGER NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
		maintenance_type_id INTEGER NOT NULL REFERENCES maintenance_types(id) ON DELETE CASCADE,
		last_performed_date TEXT NOT NULL,
		last_performed_mileage INTEGER NOT NULL,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(queries); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
