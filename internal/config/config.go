package config

import (
	"encoding/json"
	"os"
)

// Config is the configuration structure
type Config struct {
	// Imap is the imap configuration
	Imap struct {
		// Host is the imap host
		Host string `yaml:"host"`
		// Port is the imap port
		Port int `yaml:"port"`
		// Username is the imap username
		Username string `yaml:"username"`
		// Password is the imap password$
		Password string `yaml:"password"`
	} `yaml:"imap"`
}

// FromFile reads the configuration from a file
func FromFile(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromData(content)
}

// FromData reads the configuration from a byte array
func FromData(content []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
