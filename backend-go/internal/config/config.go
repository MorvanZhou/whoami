package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds runtime settings for the Go practice backend.
type Config struct {
	Env            string
	HTTPAddr       string
	DatabaseURL    string
	ProfileDataDir string
	JWTSecret      string
	JWTExpireDays  int
	CookieSecure   bool
	FrontendURL    string
	APIKeyPrefix   string

	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURI  string
}

// Load reads .env (or .env.dev when WHOAMI_ENV=dev) then environment variables.
func Load() (*Config, error) {
	envFile := ".env"
	if os.Getenv("WHOAMI_ENV") == "dev" {
		envFile = ".env.dev"
	}
	_ = godotenv.Load(envFile)
	_ = godotenv.Load(".env") // fallback / overlay for local defaults

	cfg := &Config{
		Env:            getEnv("ENV", "dev"),
		HTTPAddr:       getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "./data/whoami.db"),
		ProfileDataDir: getEnv("PROFILE_DATA_DIR", "./data/profiles"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-to-a-random-secret"),
		JWTExpireDays:  getEnvInt("JWT_EXPIRE_DAYS", 7),
		CookieSecure:   getEnvBool("COOKIE_SECURE", false),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3000"),
		APIKeyPrefix:   getEnv("API_KEY_PREFIX", "wai_"),

		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURI:  getEnv("GITHUB_REDIRECT_URI", "http://localhost:8080/api/auth/github/callback"),
	}

	if cfg.JWTSecret == "" || cfg.JWTSecret == "change-me-to-a-random-secret" {
		if cfg.Env == "prod" {
			return nil, fmt.Errorf("JWT_SECRET must be set in prod")
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var n int
	if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "off", "OFF":
		return false
	default:
		return fallback
	}
}
