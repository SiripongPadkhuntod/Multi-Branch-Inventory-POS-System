package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTTL       time.Duration
	RefreshTTL      time.Duration
	CookieSecure    bool
	RunMigrations   bool
	MigrationsPath  string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioUseSSL     bool
	MinioBucketName string
	MinioPublicURL  string
}

func Load() Config {
	accessMinutes := intEnv("ACCESS_TOKEN_TTL_MINUTES", 15)
	refreshHours := intEnv("REFRESH_TOKEN_TTL_HOURS", 168)
	appEnv := env("APP_ENV", "development")

	return Config{
		AppEnv:          appEnv,
		Port:            env("PORT", "8080"),
		DatabaseURL:     databaseURL(),
		JWTSecret:       env("JWT_SECRET", "change-me"),
		AccessTTL:       time.Duration(accessMinutes) * time.Minute,
		RefreshTTL:      time.Duration(refreshHours) * time.Hour,
		CookieSecure:    appEnv == "production",
		RunMigrations:   boolEnv("RUN_MIGRATIONS", true),
		MigrationsPath:  env("MIGRATIONS_PATH", "migrations"),
		MinioEndpoint:   env("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:  env("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:  env("MINIO_SECRET_KEY", "minioadmin"),
		MinioUseSSL:     boolEnv("MINIO_USE_SSL", false),
		MinioBucketName: env("MINIO_BUCKET_NAME", "pos-products"),
		MinioPublicURL:  env("MINIO_PUBLIC_URL", ""),
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

func boolEnv(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "TRUE", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "no", "NO", "off", "OFF":
		return false
	default:
		return fallback
	}
}
