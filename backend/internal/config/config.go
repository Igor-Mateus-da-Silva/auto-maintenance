package config

import (
	"os"
)

type Config struct {
	Port               string
	Env                string
	DBPath             string
	CORSAllowedOrigins string
	AppAccessKey       string
	JWTSecret          string
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
		dbPath = "./data/car_maintenance.db"
	}

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:4200"
	}

	appAccessKey := os.Getenv("APP_ACCESS_KEY")
	if appAccessKey == "" {
		appAccessKey = "admin123"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-jwt-key-auto-maintenance"
	}

	return &Config{
		Port:               port,
		Env:                env,
		DBPath:             dbPath,
		CORSAllowedOrigins: corsOrigins,
		AppAccessKey:       appAccessKey,
		JWTSecret:          jwtSecret,
	}
}
