package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// NewLogger 创建一个 zerolog.Logger 实例，支持控制台友好格式
func NewLogger(logFilePath string) zerolog.Logger {

	// 确保目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}
	multi := io.MultiWriter(consoleWriter, file)
	logger := zerolog.New(multi).
		With().
		Timestamp().
		Caller().
		Logger()
	logger.Info().Msg("日志初始化成功")
	return logger
}
