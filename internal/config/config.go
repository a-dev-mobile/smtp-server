package config

import (
	"errors"
	"fmt"

	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config defines the configuration structure.
type Config struct {
	Environment   string               `yaml:"environment"`
	Logging       LoggingConfig        `yaml:"logging"`
	SMTPProviders []SMTPProviderConfig `yaml:"smtpProviders"`
	GRPCServer    GRPCServerConfig     `yaml:"grpcServer"`
}

// SMTPProviderConfig defines settings for each SMTP provider.
type SMTPProviderConfig struct {
	Name      string `yaml:"name"`
	SMTPHost  string `yaml:"smtpHost"`
	SMTPPort  string `yaml:"smtpPort"`
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
	FromEmail string `yaml:"fromEmail"`
}

// LoggingConfig defines logging settings.
type LoggingConfig struct {
	Level      string           `yaml:"level"`
	FileOutput FileOutputConfig `yaml:"fileOutput"`
}

// FileOutputConfig defines settings for outputting logs to a file.
type FileOutputConfig struct {
	FilePath string `yaml:"filePath"`
}

// GRPCServerConfig defines the gRPC server settings.
type GRPCServerConfig struct {
	Port                 string `yaml:"port"`
	MaxConcurrentStreams int    `yaml:"maxConcurrentStreams"`
}

// GetConfig loads configuration from file
func GetConfig() (*Config, error) {
	configPath := "../config"

	configFile := filepath.Join(configPath, "config.yaml")

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("config file does not exist: %s", configFile)
	}

	return loadConfig(configFile)
}

// loadConfig reads and decodes the YAML configuration file.
func loadConfig(file string) (*Config, error) {
	configFile, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	return &config, nil
}
