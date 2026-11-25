package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"car-maintenance/internal/database"
)

func TestDatabase_ConnectAndMigrations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "subfolder", "test.db")

	db, err := database.Connect(dbPath)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify file was created in subdirectory
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("Database file was not created at expected path: %s", dbPath)
	}

	// Verify foreign_keys pragma is enabled
	var foreignKeysEnabled int
	if err := db.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeysEnabled); err != nil {
		t.Fatalf("Failed to query PRAGMA foreign_keys: %v", err)
	}
	if foreignKeysEnabled != 1 {
		t.Errorf("Expected foreign_keys = 1, got %d", foreignKeysEnabled)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Check that tables exist
	tables := []string{"vehicles", "maintenance_types", "vehicle_maintenances"}
	for _, table := range tables {
		var name string
		query := "SELECT name FROM sqlite_master WHERE type='table' AND name=?;"
		if err := db.QueryRow(query, table).Scan(&name); err != nil {
			t.Errorf("Table %s was not found after migrations: %v", table, err)
		}
	}
}
