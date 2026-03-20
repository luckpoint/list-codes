package tui

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Include []string       `yaml:"include,omitempty"`
	Exclude []string       `yaml:"exclude,omitempty"`
	Options *ConfigOptions `yaml:"options,omitempty"`
}

type ConfigOptions struct {
	IncludeTests bool   `yaml:"include-tests,omitempty"`
	MaxFileSize  string `yaml:"max-file-size,omitempty"`
	MaxDepth     int    `yaml:"max-depth,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
