package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	JwtToken      string `yaml:"jwt_token"`
	OutServiceURL string `yaml:"out_service_url"`
}

func ReadConfig() *Config {
	var config *Config
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		logrus.Fatalf("Err read config.yaml file - %s", err)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		logrus.Fatalf("Err unmarshal config.yaml to struct - %s", err)
	}

	return config
}
