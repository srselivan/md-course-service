package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Level       string
	AppName     string
	LogFilePath string
}

func New(cfg Config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse level: %w", err)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:          os.Stdout,
		NoColor:      true,
		TimeFormat:   time.RFC3339Nano,
		TimeLocation: time.UTC,
	}

	logFile, err := openOrCreateLogFile(cfg.LogFilePath, cfg.AppName)
	if err != nil {
		return nil, fmt.Errorf("open or create log file: %w", err)
	}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)

	logger := zerolog.New(multiWriter).With().Caller().Timestamp().Logger()
	logger.Level(level)

	return &logger, nil
}

func openOrCreateLogFile(filepath string, filename string) (*os.File, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err = os.Mkdir(filepath, 0777); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}
	}
	logsFilePath := fmt.Sprintf("%s/%s.log", filepath, filename)
	file, err := os.OpenFile(logsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("openfile: %w", err)
	}
	return file, nil
}
