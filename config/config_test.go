package config

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetAppConfig() {
	// Reset the singleton instance for testing purposes
	appConfig = bestiaryConfig{}
	onceAppConfig = sync.Once{}
}

func TestGetAppConfig(t *testing.T) {
	// Save current environment variables
	originalDatabaseURL := os.Getenv("BESTIARY_DATABASE_URL")
	originalPort := os.Getenv("BESTIARY_PORT")

	// Unset environment variables for the test
	os.Unsetenv("BESTIARY_DATABASE_URL")
	os.Unsetenv("BESTIARY_PORT")

	// Restore environment variables after the test
	defer func() {
		if originalDatabaseURL != "" {
			os.Setenv("BESTIARY_DATABASE_URL", originalDatabaseURL)
		}
		if originalPort != "" {
			os.Setenv("BESTIARY_PORT", originalPort)
		}
	}()

	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yml"

	// Sample configuration data
	configData := `db_url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
port: "8080"`

	// Write the sample config data to the temporary config file
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test without environment variables
	resetAppConfig() // Reset the config
	config := GetAppConfig(configPath)
	expectedConfig := &bestiaryConfig{
		DatabaseUrl: "postgres://user:password@localhost:5432/dbname?sslmode=disable",
		Port:        "8080",
	}
	assert.Equal(t, expectedConfig, config)

	// Test with environment variables
	os.Setenv("BESTIARY_DATABASE_URL", "postgres://envuser:envpassword@localhost:5432/envdbname?sslmode=disable")
	os.Setenv("BESTIARY_PORT", "9090")

	resetAppConfig() // Reset the config
	config = GetAppConfig(configPath)
	expectedConfig = &bestiaryConfig{
		DatabaseUrl: "postgres://envuser:envpassword@localhost:5432/envdbname?sslmode=disable",
		Port:        "9090",
	}
	assert.Equal(t, expectedConfig, config)
}

func TestGetAppConfig_InvalidPath(t *testing.T) {
	// Save current environment variables
	originalDatabaseURL := os.Getenv("BESTIARY_DATABASE_URL")
	originalPort := os.Getenv("BESTIARY_PORT")

	// Unset environment variables for the test
	os.Unsetenv("BESTIARY_DATABASE_URL")
	os.Unsetenv("BESTIARY_PORT")

	// Restore environment variables after the test
	defer func() {
		if originalDatabaseURL != "" {
			os.Setenv("BESTIARY_DATABASE_URL", originalDatabaseURL)
		}
		if originalPort != "" {
			os.Setenv("BESTIARY_PORT", originalPort)
		}
	}()

	// Test: Load configuration from an invalid path
	resetAppConfig() // Reset the config
	assert.Panics(t, func() {
		GetAppConfig("/invalid/path/config.yml")
	}, "The code did not panic on invalid file path")
}

func TestGetAppConfig_InvalidYaml(t *testing.T) {
	// Save current environment variables
	originalDatabaseURL := os.Getenv("BESTIARY_DATABASE_URL")
	originalPort := os.Getenv("BESTIARY_PORT")

	// Unset environment variables for the test
	os.Unsetenv("BESTIARY_DATABASE_URL")
	os.Unsetenv("BESTIARY_PORT")

	// Restore environment variables after the test
	defer func() {
		if originalDatabaseURL != "" {
			os.Setenv("BESTIARY_DATABASE_URL", originalDatabaseURL)
		}
		if originalPort != "" {
			os.Setenv("BESTIARY_PORT", originalPort)
		}
	}()

	// Setup: Create a temporary config file with invalid YAML
	tempDir := t.TempDir()
	configPath := tempDir + "/config_invalid.yml"

	// Invalid YAML configuration data
	configData := `
db_url: "postgres://user:password@localhost:5432/dbname?sslmode=disable
`

	// Write the invalid config data to the temporary config file
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test: Load the configuration and expect a panic
	resetAppConfig() // Reset the config
	assert.Panics(t, func() {
		GetAppConfig(configPath)
	}, "The code did not panic on invalid YAML")
}
