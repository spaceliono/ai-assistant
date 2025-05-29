package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`
	AI struct {
		Moonshot struct {
			BaseURL string `yaml:"baseURL"`
			APIKey  string `yaml:"api_key"`
			Model   string `yaml:"model"`
		} `yaml:"moonshot"`
		Deepseek struct {
			BaseURL       string `yaml:"baseURL"`
			APIKey        string `yaml:"api_key"`
			Model         string `yaml:"model"`
			SupportModels string `yaml:"supportModels"`
		} `yaml:"deepseek"`
	} `yaml:"ai"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
