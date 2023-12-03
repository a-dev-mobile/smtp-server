package config

import (
    "errors"
    // Errors package for handling errors.
    "fmt"
    // Fmt package for formatting strings.
    "os"
    // Os package for interacting with the operating system, like file handling.
    "path/filepath"
    // Filepath package for manipulating filename paths.
    "gopkg.in/yaml.v3"
    // Yaml.v3 package for YAML processing.
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






// GetConfig loads configuration from file.
// This function is generic and works with any type T.
// It takes a configuration path and name, then returns a pointer to the config struct or an error.
func GetConfig[T any](configPath string, configName string) (*T, error) {
    configFile := filepath.Join(configPath, configName)
    // Joins the path and filename to create the full path to the configuration file.

    if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
        // Checks if the file exists. If it does not, returns an error.
        return nil, fmt.Errorf("config file does not exist: %s", configFile)
    }

    // If the file exists, calls loadConfig to read and parse the file.
    return loadConfig[T](configFile)
}

// loadConfig reads and decodes the YAML configuration file.
// It is a private function, indicated by the lowercase first letter.
// Takes the file path as input and returns a pointer to the config struct or an error.
func loadConfig[T any](file string) (*T, error) {
    configFile, err := os.ReadFile(file)
    // Reads the file. If there is an error reading, it returns an error.
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config T
    // Declares a variable of type T to hold the configuration data.

    if err := yaml.Unmarshal(configFile, &config); err != nil {
        // Tries to unmarshal the YAML file into the config variable. If it fails, returns an error.
        return nil, fmt.Errorf("error decoding config file: %w", err)
    }

    // Returns a pointer to the config struct if successful.
    return &config, nil
}
