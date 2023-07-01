package main

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AuthToken   string   `yaml:"authToken"`
	MonitGuilds []string `yaml:"monitGuilds"`
	MyUID       string   `yaml:"myUID"`
}

func InitConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("main - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	if cfg.AuthToken == "" {
		return nil, errors.New("it is required to set the authentication token in config.yml")
	}
	if len(cfg.MonitGuilds) == 0 {
		return nil, errors.New("you need to specify at least 1 guild to track the change of the nickname in config.yml")
	}
	if cfg.MyUID == "" {
		return nil, errors.New("it is required to set the myUID in config.yml")
	}

	return cfg, nil
}
