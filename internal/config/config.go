// Package config loads configuration from environment variables and .env file.
package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	APIName          string `env:"MB_API_NAME"`
	APIVersion       string `env:"MB_API_VERSION"`
	ServerPort       string `env:"MB_API_SERVER_PORT"`
	ServerEnv        string `env:"MB_API_SERVER_ENV,optional"`
	ServerLogLevel   string `env:"MB_API_SERVER_LOG_LEVEL"`
	PostgresSchema   string `env:"MB_API_POSTGRES_SCHEMA"`
	PostgresDsn      string `env:"MB_API_POSTGRES_DSN"`
	PostgresLogLevel string `env:"MB_API_POSTGRES_LOG_LEVEL"`
	TelegramBotToken string `env:"MB_API_TELEGRAM_BOT_TOKEN"`
	TelegramChatID   string `env:"MB_API_TELEGRAM_CHAT_ID"`
}

var (
	SingleLine string = "--------------------------------------------------"
	DoubleLine string = "=================================================="
)

var (
	instance *Config
	once     sync.Once
	err      error
)

// Get returns the application configuration
func Get() (*Config, error) {
	once.Do(func() {
		instance, err = loadConfig()
	})
	return instance, err
}

// loadConfig loads configuration from .env file and environment variables
func loadConfig() (*Config, error) {
	// Set environment based on hostname before loading any env files
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %w", err)
	}

	// Set default environment based on hostname
	if hostname == "moneybots-app" {
		os.Setenv("MB_API_SERVER_ENV", "production")
	} else {
		os.Setenv("MB_API_SERVER_ENV", "development")
	}

	// Load .env file if it exists
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}
	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// loadEnvFile attempts to load the .env file
func loadEnvFile() error {
	// Look for .env file in the current directory and parent directories
	envFiles := []string{
		".env",
		"../.env",
		"../../.env",
	}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			return godotenv.Load(envFile)
		}
	}

	// If no .env file is found, return without error as environment variables
	// might be set through other means (e.g., docker, kubernetes)
	return nil
}

// loadFromEnv loads configuration from environment variables
func (c *Config) loadFromEnv() error {
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(c).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			return fmt.Errorf("missing env tag for field %s", field.Name)
		}

		// Parse the env tag to get the key and options
		tagParts := strings.Split(envTag, ",")
		envKey := tagParts[0]
		isOptional := len(tagParts) > 1 && tagParts[1] == "optional"

		value := os.Getenv(envKey)
		if value == "" && !isOptional {
			return fmt.Errorf("env variable %s is required but not set", envKey)
		}

		v.Field(i).SetString(value)
	}

	return nil
}

// String returns the configuration as a string
func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString(SingleLine + "\n")
	sb.WriteString("Config: \n")
	sb.WriteString(SingleLine + "\n")

	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).String()

		// Mask sensitive fields
		value = maskSensitiveField(field.Name, value)
		sb.WriteString(fmt.Sprintf("  %s:  %s\n", field.Name, value))
	}

	sb.WriteString(SingleLine + "\n")

	return sb.String()
}

func maskSensitiveField(fieldName, value string) string {
	sensitiveFields := []string{"token", "dsn", "secret", "password", "url"}

	fieldNameLower := strings.ToLower(fieldName)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldNameLower, sensitive) {
			return maskValue(value)
		}
	}

	return value
}

func maskValue(value string) string {
	if len(value) <= 3 {
		return strings.Repeat("*", 7)
	}
	return value[:3] + strings.Repeat("*", 7)
}
