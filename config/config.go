package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type bestiaryConfig struct {
	DatabaseUrl string `yaml:"db_url"`
}

var appConfig bestiaryConfig
var onceAppConfig sync.Once

func GetAppConfig() *bestiaryConfig {
	onceAppConfig.Do(func() {
		var err error
		rawConfig, err := os.ReadFile("config/config.yml")
		if err != nil {
			fmt.Println("error:", err)
		}

		err = yaml.Unmarshal(rawConfig, &appConfig)
		if err != nil {
			fmt.Println("error:", err)
		}
		if err != nil {
			panic(err)
		}
	})

	return &appConfig
}
