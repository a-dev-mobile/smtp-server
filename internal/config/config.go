package config

import (

	"errors"
	"fmt"

	"os"
	"path/filepath"
	"github.com/a-dev-mobile/smtp-server/internal/environment"

	"gopkg.in/yaml.v3"
)

// Config определяет структуру конфигурации.
type Config struct {
	Logging       LoggingConfig        `yaml:"logging"`
	SMTPProviders []SMTPProviderConfig `yaml:"smtpProviders"`
	GRPCServer    GRPCServerConfig     `yaml:"grpcServer"`
}
// SMTPProviderConfig определяет настройки для каждого поставщика SMTP.
type SMTPProviderConfig struct {
	Name      string `yaml:"name"`
	SMTPHost  string `yaml:"smtpHost"`
	SMTPPort  string `yaml:"smtpPort"`
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
	FromEmail string `yaml:"fromEmail"`
}

// LoggingConfig определяет настройки логирования.
type LoggingConfig struct {
	Level      string           `yaml:"level"`
	FileOutput FileOutputConfig `yaml:"fileOutput"`
}

// FileOutputConfig определяет настройки вывода логов в файл.
type FileOutputConfig struct {
	FilePath string `yaml:"filePath"`
}

// GRPCServerConfig определяет настройки gRPC сервера.
type GRPCServerConfig struct {
	Port                 string `yaml:"port"`
	MaxConcurrentStreams int    `yaml:"maxConcurrentStreams"`
}


// GetConfig загружает конфигурацию из файла в зависимости от окружения.
func GetConfig(appEnv environment.Environment) (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	var configFile string
	switch appEnv {
	case environment.Prod:
		configFile = filepath.Join(configPath, "config.prod.yaml")
	case environment.Dev:
		configFile = filepath.Join(configPath, "config.dev.yaml")
	default:
		return nil, fmt.Errorf("invalid environment: expected %s or %s, got %v", environment.Prod, environment.Dev, appEnv)
	}

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("config file does not exist: %s", configFile)
	}

	return loadConfig(configFile)
}

// loadConfig читает и декодирует YAML файл конфигурации.
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