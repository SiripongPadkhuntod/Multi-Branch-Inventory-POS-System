package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv       string
	Port         string
	DatabaseURL  string
	JWTSecret    string
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	CookieSecure bool
}

func Load() Config {
	accessMinutes := intEnv("ACCESS_TOKEN_TTL_MINUTES", 15)
	refreshHours := intEnv("REFRESH_TOKEN_TTL_HOURS", 168)
	appEnv := env("APP_ENV", "development")

	return Config{
		AppEnv:       appEnv,
		Port:         env("PORT", "8080"),
		DatabaseURL:  databaseURL(),
		JWTSecret:    env("JWT_SECRET", "change-me"),
		AccessTTL:    time.Duration(accessMinutes) * time.Minute,
		RefreshTTL:   time.Duration(refreshHours) * time.Hour,
		CookieSecure: appEnv == "production",
	}
}

func databaseURL() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://" + env("DB_USER", "pos") + ":" + env("DB_PASSWORD", "pos_password") +
		"@" + env("DB_HOST", "localhost") + ":" + env("DB_PORT", "5432") + "/" + env("DB_NAME", "pos") +
		"?sslmode=disable"
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func intEnv(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}
