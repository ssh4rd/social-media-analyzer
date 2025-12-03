package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	VK       VKConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type VKConfig struct {
	AccessToken string
	APIVersion  string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "3000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "database"),
		},
		VK: VKConfig{
			AccessToken: getEnv("VK_ACCESS_TOKEN", ""),
			APIVersion:  getEnv("VK_API_VERSION", "5.131"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}
