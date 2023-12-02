package models



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
