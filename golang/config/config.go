package config

import "errors"

type Config struct {
	TasksQueueLimit int
}

func (c *Config) Validate() error {
	if c.TasksQueueLimit < 0 {
		return errors.New("invalid config")
	}

	return nil
}
