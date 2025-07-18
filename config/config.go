package config

import (
	"os"
	"time"

	"EffectiveMobile/pkg/logster"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HTTPClient `yaml:"httpClient"`
	Log        logster.Config `yaml:"log"`
}

type HTTPClient struct {
	Timeout     time.Duration `yaml:"timeout"`
	PrivateAddr string        `yaml:"privateAddr"`
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
