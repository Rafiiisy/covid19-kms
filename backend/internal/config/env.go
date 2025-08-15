package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnvFile loads environment variables from a .env file
func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) {
			value = value[1 : len(value)-1]
		}

		// Set environment variable
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// LoadEnvFileIfExists loads .env file if it exists, otherwise continues silently
func LoadEnvFileIfExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, that's okay
		return nil
	}
	return LoadEnvFile(filename)
}

// LoadDefaultEnv loads environment from common .env file locations
func LoadDefaultEnv() error {
	// Try to load from current directory
	if err := LoadEnvFileIfExists(".env"); err != nil {
		return err
	}

	// Try to load from config directory
	if err := LoadEnvFileIfExists("config/.env"); err != nil {
		return err
	}

	// Try to load from parent directory
	if err := LoadEnvFileIfExists("../.env"); err != nil {
		return err
	}

	return nil
}

// GetEnvWithDefault gets an environment variable with a default value
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetRequiredEnv gets a required environment variable, returns error if not set
func GetRequiredEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("required environment variable %s is not set", key)
}

// ValidateRequiredEnvs validates that all required environment variables are set
func ValidateRequiredEnvs(requiredKeys []string) error {
	var missing []string

	for _, key := range requiredKeys {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

// GetRequiredEnvsForProduction returns the list of environment variables required for production
func GetRequiredEnvsForProduction() []string {
	return []string{
		"YOUTUBE_API_KEY",
		"GOOGLE_NEWS_API_KEY",
		"INSTAGRAM_API_KEY",
		"INDONESIA_NEWS_API_KEY",
	}
}

// GetRequiredEnvsForDevelopment returns the list of environment variables required for development
func GetRequiredEnvsForDevelopment() []string {
	return []string{
		// For development, API keys are optional (can use mock data)
		"ENV",
	}
}

// IsProduction checks if the environment is set to production
func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

// IsDevelopment checks if the environment is set to development
func IsDevelopment() bool {
	return os.Getenv("ENV") == "development" || os.Getenv("ENV") == ""
}

// IsTest checks if the environment is set to test
func IsTest() bool {
	return os.Getenv("ENV") == "test"
}
