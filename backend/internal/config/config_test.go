package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"car-maintenance/internal/config"
)

func TestLoad_MissingRequiredSecrets(t *testing.T) {
	// Run inside an empty temp dir so local .env does not interfere
	t.Chdir(t.TempDir())
	t.Setenv("APP_ACCESS_KEY", "")
	t.Setenv("JWT_SECRET", "")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected error when APP_ACCESS_KEY is missing")
	}

	t.Setenv("APP_ACCESS_KEY", "my-strong-key")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	}

	t.Setenv("JWT_SECRET", "too-short")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error when JWT_SECRET is shorter than 32 chars")
	}
}

func TestLoad_ValidEnvVars(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("APP_ACCESS_KEY", "strong-access-key")
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AppAccessKey != "strong-access-key" {
		t.Errorf("got AppAccessKey=%q", cfg.AppAccessKey)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default Port 8080, got %q", cfg.Port)
	}
}

func TestLoad_FromDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("APP_ACCESS_KEY", "")
	t.Setenv("JWT_SECRET", "")

	dotEnvContent := "APP_ACCESS_KEY=from-dotenv-key\nJWT_SECRET=abcdef0123456789abcdef0123456789\nPORT=9090\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(dotEnvContent), 0600); err != nil {
		t.Fatalf("failed to write temp .env: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading from .env: %v", err)
	}
	if cfg.AppAccessKey != "from-dotenv-key" || cfg.Port != "9090" {
		t.Errorf("unexpected config from .env: %+v", cfg)
	}
}
