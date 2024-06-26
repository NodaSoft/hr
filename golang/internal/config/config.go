package config

import (
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TaskResultReportingPeriod time.Duration `yaml:"task_result_reporting_period"`
	TaskExpirationDuration    time.Duration `yaml:"task_expiration_duration"`
	TaskGenerationDuration    time.Duration `yaml:"task_generation_duration"`
	TaskExecutorsLimit        int           `yaml:"executors_limit"`
}

var (
	config Config
	once   sync.Once
)

func MustNew() Config {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("could not get pwd for config path: %v", err)
			}
			configPath = path.Join(pwd, "/config/config.yaml")
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}

		if err := cleanenv.ReadConfig(configPath, &config); err != nil {
			log.Fatalf("cannot read config: %s", err)
		}
	})

	return config
}
