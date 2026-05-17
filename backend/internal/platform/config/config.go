package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Database DatabaseConfig
	Clover   CloverConfig
	Server   ServerConfig
	Security SecurityConfig
	AppEnv   string
}

// DatabaseConfig holds PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

// CloverConfig holds Clover API credentials.
type CloverConfig struct {
	ClientID     string
	ClientSecret string
	Env          string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port        string
	FrontendURL string
}

// SecurityConfig holds encryption and JWT secrets.
type SecurityConfig struct {
	EncryptionKey string
	JWTSecret     string
}

// findEnvFile searches for .env starting from the current working directory
// and walking up the directory tree until found or reaching the root.
func findEnvFile() (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		candidate := filepath.Join(wd, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break // reached root
		}
		wd = parent
	}
	return "", false
}

// Load reads configuration from .env file first, then environment variables.
// In development, .env is the primary source. Environment variables can override.
func Load() (*Config, error) {
	// Load .env file into environment variables (development default)
	if envPath, found := findEnvFile(); found {
		if err := godotenv.Load(envPath); err != nil {
			return nil, fmt.Errorf("read .env at %s: %w", envPath, err)
		}
	}

	v := viper.New()

	// Environment variables override .env values
	v.SetEnvPrefix("") // no prefix
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	setDefaults(v)

	// Manual mapping — more reliable than Unmarshal with squash
	cfg := Config{
		Database: DatabaseConfig{
			Host:         v.GetString("DB_HOST"),
			Port:         v.GetString("DB_PORT"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			Database:     v.GetString("DB_NAME"),
			SSLMode:      v.GetString("DB_SSL_MODE"),
			MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
		},
		Clover: CloverConfig{
			ClientID:     v.GetString("CLOVER_CLIENT_ID"),
			ClientSecret: v.GetString("CLOVER_CLIENT_SECRET"),
			Env:          v.GetString("CLOVER_ENV"),
		},
		Server: ServerConfig{
			Port:        v.GetString("API_PORT"),
			FrontendURL: v.GetString("FRONTEND_URL"),
		},
		Security: SecurityConfig{
			EncryptionKey: v.GetString("ENCRYPTION_KEY"),
			JWTSecret:     v.GetString("JWT_SECRET"),
		},
		AppEnv: v.GetString("APP_ENV"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	// Database
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "sales_insight")
	v.SetDefault("DB_PASSWORD", "sales_insight_secret")
	v.SetDefault("DB_NAME", "sales_insight")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)

	// Clover
	v.SetDefault("CLOVER_ENV", "sandbox")

	// Server
	v.SetDefault("API_PORT", "8080")
	v.SetDefault("FRONTEND_URL", "http://localhost:5173")

	// App
	v.SetDefault("APP_ENV", "development")
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.Security.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

// DatabaseDSN builds the PostgreSQL connection string from config.
func (c *Config) DatabaseDSN() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + c.Database.Port +
		"/" + c.Database.Database + "?sslmode=" + c.Database.SSLMode
}
