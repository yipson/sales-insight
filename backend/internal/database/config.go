package database

import "os"

// Config holds database connection parameters.
type Config struct {
	Host         string `mapstructure:"DB_HOST"`
	Port         string `mapstructure:"DB_PORT"`
	User         string `mapstructure:"DB_USER"`
	Password     string `mapstructure:"DB_PASSWORD"`
	Database     string `mapstructure:"DB_NAME"`
	SSLMode      string `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
}

// DefaultConfig returns development defaults.
// In production, load these from environment variables via Viper.
func DefaultConfig() Config {
	return Config{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnv("DB_PORT", "5432"),
		User:         getEnv("DB_USER", "sales_insight"),
		Password:     getEnv("DB_PASSWORD", "sales_insight_secret"),
		Database:     getEnv("DB_NAME", "sales_insight"),
		SSLMode:      getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns: 25,
		MaxIdleConns: 10,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

// DSN builds the PostgreSQL connection string.
func (c Config) DSN() string {
	return "postgres://" + c.User + ":" + c.Password +
		"@" + c.Host + ":" + c.Port +
		"/" + c.Database + "?sslmode=" + c.SSLMode
}
