package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPClient
}

type HTTPClient struct {
	Timeout time.Duration `yaml:"timeout"`
}

func LoadConfig(filename string, cfg interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	return nil
}
