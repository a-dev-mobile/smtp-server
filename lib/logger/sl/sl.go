package sl

import (
	"github.com/a-dev-mobile/smtp-server/internal/environment"

	"golang.org/x/exp/slog"

	"log"
	"os"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)
const LogFilePath = "/app/logs/smtp-server.log"

func SetupLogger(env environment.Environment, logLevel string, logFilePath string) *slog.Logger {
	var logger *slog.Logger
	level := parseLogLevel(logLevel)

	// Проверка, задан ли путь к файлу для логов
	if logFilePath != "" {
		// Попытка открыть файл журнала
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level}))
	} else {
		// Использование стандартного вывода, если путь к файлу не задан
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
		if env != environment.Dev {
			logger.Warn("Unknown environment: %s, using default logger settings", env)
		}
	}

	return logger
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func parseLogLevel(level string) slog.Level {
	if !isValidLogLevel(level) {
		log.Fatalf("Invalid logging level: %s", level)
	}

	switch level {
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	case LevelDebug:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func isValidLogLevel(level string) bool {
	switch level {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return true
	default:
		return false
	}
}
