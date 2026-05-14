package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Database   DatabaseConfig   `mapstructure:",squash"`
	Clover     CloverConfig     `mapstructure:",squash"`
	Server     ServerConfig     `mapstructure:",squash"`
	Security   SecurityConfig   `mapstructure:",squash"`
	AppEnv     string           `mapstructure:"APP_ENV"`
}

// DatabaseConfig holds PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host         string `mapstructure:"DB_HOST"`
	Port         string `mapstructure:"DB_PORT"`
	User         string `mapstructure:"DB_USER"`
	Password     string `mapstructure:"DB_PASSWORD"`
	Database     string `mapstructure:"DB_NAME"`
	SSLMode      string `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
}

// CloverConfig holds Clover API credentials.
type CloverConfig struct {
	ClientID     string `mapstructure:"CLOVER_CLIENT_ID"`
	ClientSecret string `mapstructure:"CLOVER_CLIENT_SECRET"`
	Env          string `mapstructure:"CLOVER_ENV"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port        string `mapstructure:"API_PORT"`
	FrontendURL string `mapstructure:"FRONTEND_URL"`
}

// SecurityConfig holds encryption and JWT secrets.
type SecurityConfig struct {
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
}

// Load reads configuration from .env file first, then environment variables.
// In development, .env is the primary source. Environment variables can override.
func Load() (*Config, error) {
	v := viper.New()

	// Read from .env file (primary source for local development)
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err != nil {
		// Only fail if .env exists but is unreadable; if missing, rely on env vars
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read .env: %w", err)
		}
	}

	// Environment variables override .env values
	v.SetEnvPrefix("") // no prefix
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
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
