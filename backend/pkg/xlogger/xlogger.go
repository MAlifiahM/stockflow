package xlogger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/MAlifiahM/stockflow/backend/internal/config"
	"github.com/rs/zerolog"
)

var (
	Logger *zerolog.Logger
)

func NewFileLogger(logFilePath string, isDevelopment bool) *zerolog.Logger {
	var output io.Writer

	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fallBackLogger := zerolog.New(os.Stderr).
			With().Timestamp().Logger().With().Str("component", "xlogger").Logger()
		fallBackLogger.Error().Err(err).Str("logDir", logDir).Msg("Failed to create log directory, falling back to stderr")
		output = os.Stderr
	} else {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fallBackLogger := zerolog.New(os.Stderr).
				With().
				Timestamp().
				Logger().
				With().
				Str("component", "xlogger").
				Logger()
			fallBackLogger.Error().
				Err(err).
				Str("logFilePath", logFilePath).
				Msg("Failed to open log file, falling back to stderr")
			output = os.Stderr
		} else {
			output = file
		}
	}

	var l zerolog.Logger
	if isDevelopment {
		multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}, output)
		l = zerolog.New(multi).With().Timestamp().Logger()
		l.Level(zerolog.DebugLevel)
	} else {
		l = zerolog.New(output).With().Timestamp().Logger()
		l.Level(zerolog.InfoLevel)
	}
	return &l
}

func Setup(cfg config.Config) {
	if cfg.IsDevelopment {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
		l.Level(zerolog.DebugLevel)
		Logger = &l
		return
	}
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	Logger = &l
}
