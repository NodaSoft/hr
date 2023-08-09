package config

import "time"

type TaskConfig struct {
	TaskRelevanceTimDuration time.Duration `mapstructure:"taskRelevanceTime"`
}
