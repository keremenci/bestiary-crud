package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAppConfig(t *testing.T) {
	// Setup: Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.yml"

	// Sample configuration data
	configData := `
db_url: "postgres://user:password@localhost:5432/dbname?sslmode=disable"
`

	// Write the sample config data to the temporary config file
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test: Load the configuration
	config := GetAppConfig(configPath)

	// Verify: Check if the configuration was loaded correctly
	expectedConfig := &bestiaryConfig{
		DatabaseUrl: "postgres://user:password@localhost:5432/dbname?sslmode=disable",
	}
	assert.Equal(t, expectedConfig, config)
}

func TestGetAppConfig_InvalidPath(t *testing.T) {
	// Test: Load configuration from an invalid path
	assert.Panics(t, func() {
		GetAppConfig("/invalid/path/config.yml")
	}, "The code did not panic on invalid file path")
}

func TestGetAppConfig_InvalidYaml(t *testing.T) {
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
	assert.Panics(t, func() {
		GetAppConfig(configPath)
	}, "The code did not panic on invalid YAML")
}
