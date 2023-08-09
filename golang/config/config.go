package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Task     TaskConfig    `mapstructure:"task"`
	WorkTime time.Duration `mapstructure:"workTime"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
