package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// NewLogger 创建一个 zerolog.Logger 实例，支持控制台友好格式
func NewLogger() zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}

	logger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}
