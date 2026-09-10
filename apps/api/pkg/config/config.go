package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
	OIDC     OIDCConfig
	R2       R2Config
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CORSConfig struct {
	AllowOrigins string
}

type OIDCConfig struct {
	Issuer       string
	ServerURL    string
	FrontendURL  string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string
	PublicBase      string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: envOrDefault("SERVER_PORT", "3711"),
			Mode: envOrDefault("SERVER_MODE", "dev"),
		},
		Database: DatabaseConfig{
			URL:             os.Getenv("DATABASE_URL"),
			MaxOpenConns:    envOrDefaultInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    envOrDefaultInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: envOrDefaultInt("DB_CONN_MAX_LIFETIME", 300),
		},
		Redis: RedisConfig{
			Host:     envOrDefault("REDIS_HOST", "127.0.0.1"),
			Port:     envOrDefault("REDIS_PORT", "6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       envOrDefaultInt("REDIS_DB", 0),
		},
		CORS: CORSConfig{
			AllowOrigins: envOrDefault("CORS_ALLOW_ORIGINS", "http://127.0.0.1:3710"),
		},
		OIDC: OIDCConfig{
			Issuer:       envOrDefault("OIDC_ISSUER", "http://127.0.0.1:9277"),
			ServerURL:    envOrDefault("OIDC_SERVER_URL", "http://127.0.0.1:9277/api/v1"),
			FrontendURL:  envOrDefault("OIDC_FRONTEND_URL", "http://127.0.0.1:9420"),
			ClientID:     os.Getenv("OIDC_CLIENT_ID"),
			ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
			RedirectURI:  envOrDefault("OIDC_REDIRECT_URI", "http://127.0.0.1:3710/auth/callback"),
		},
		R2: R2Config{
			AccountID:       os.Getenv("R2_ACCOUNT_ID"),
			AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
			Bucket:          envOrDefault("R2_BUCKET", "vnfm"),
			Endpoint:        os.Getenv("R2_ENDPOINT"),
			PublicBase:      os.Getenv("R2_PUBLIC_BASE"),
		},
	}, nil
}

func RequireDatabase(cfg *Config) error {
	if cfg.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}
