package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// bestiaryConfig holds the configuration details
// Each field can have a corresponding environment variable.
type bestiaryConfig struct {
	DatabaseUrl string `yaml:"db_url"`
	Port        string `yaml:"port"`
}

var appConfig bestiaryConfig
var onceAppConfig sync.Once

// GetEnv returns the value of the environment variable if set, otherwise returns the default value.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetAppConfig(filepath string) *bestiaryConfig {
	// Singleton
	onceAppConfig.Do(func() {
		var err error
		rawConfig, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(rawConfig, &appConfig)
		if err != nil {
			panic(err)
		}

		// Override config entries with environment variables if they are set
		appConfig.DatabaseUrl = GetEnv("BESTIARY_DATABASE_URL", appConfig.DatabaseUrl)
		appConfig.Port = GetEnv("BESTIARY_PORT", appConfig.Port)
	})

	return &appConfig
}
