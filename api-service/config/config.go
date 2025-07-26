package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(stage string) (*Config, error) {
	viper.AddConfigPath(".env")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Faled to read config file: %s", err.Error())
		return nil, err
	}

	var cfg FConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Errorf("Faled to read config data: %s", err.Error())
		return nil, err
	}

	if stage == "test" {
		return &cfg.Test, nil
	}
	return &cfg.Prod, nil
}
