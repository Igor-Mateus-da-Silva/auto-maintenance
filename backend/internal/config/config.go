package config

import (
	"os"
)

type Config struct {
	Port               string
	Env                string
	DBPath             string
	CORSAllowedOrigins string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./car_maintenance.db"
	}

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:4200"
	}

	return &Config{
		Port:               port,
		Env:                env,
		DBPath:             dbPath,
		CORSAllowedOrigins: corsOrigins,
	}
}
