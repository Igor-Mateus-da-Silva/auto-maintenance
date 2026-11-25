package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const MinJWTSecretLength = 32

type Config struct {
	Port               string
	Env                string
	DBPath             string
	CORSAllowedOrigins string
	AppAccessKey       string
	JWTSecret          string
}

// Load reads configuration from environment variables (and an optional .env file
// in the working directory for variables not already set in the environment).
// It fails fast if required security secrets are missing or too weak.
func Load() (*Config, error) {
	fileEnv, err := readDotEnv(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	get := func(key, defaultVal string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		if val, ok := fileEnv[key]; ok && val != "" {
			return val
		}
		return defaultVal
	}

	appAccessKey := strings.TrimSpace(get("APP_ACCESS_KEY", ""))
	if appAccessKey == "" {
		return nil, errors.New("APP_ACCESS_KEY é obrigatória e não pode estar vazia")
	}

	jwtSecret := get("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET é obrigatória e não pode estar vazia")
	}
	if len(jwtSecret) < MinJWTSecretLength {
		return nil, fmt.Errorf("JWT_SECRET deve ter pelo menos %d caracteres", MinJWTSecretLength)
	}

	return &Config{
		Port:               get("PORT", "8080"),
		Env:                get("ENV", "development"),
		DBPath:             get("DB_PATH", "./data/car_maintenance.db"),
		CORSAllowedOrigins: get("CORS_ALLOWED_ORIGINS", "http://localhost:4200"),
		AppAccessKey:       appAccessKey,
		JWTSecret:          jwtSecret,
	}, nil
}

func readDotEnv(path string) (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return result, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"'`)
		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
