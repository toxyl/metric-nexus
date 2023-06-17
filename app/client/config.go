package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	APIKey string `yaml:"key"`
}

func LoadConfig(file string) (*Config, error) {
	file, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	if !fileExists(file) {
		return nil, fmt.Errorf("file does not exist")
	}
	c := &Config{
		Host:   "",
		Port:   0,
		APIKey: "",
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
