package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host      string   `yaml:"host"`
	Port      int      `yaml:"port"`
	StateFile string   `yaml:"state"`
	CertFile  string   `yaml:"cert"`
	KeyFile   string   `yaml:"key"`
	APIKeys   []string `yaml:"keys"`
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
		Host:      "",
		Port:      0,
		StateFile: "",
		CertFile:  "",
		KeyFile:   "",
		APIKeys:   []string{},
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}

	if c.StateFile == "" {
		c.StateFile = strings.TrimSuffix(strings.TrimSuffix(file, ".yaml"), ".yml") + ".state.yaml"
	}
	return c, nil
}
