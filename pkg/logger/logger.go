package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level     string
	Filename  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
	IsDev     string
}

func NewLogger(config LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var write io.Writer

	if config.IsDev == "development" {
		write = PrettyJSONWriter{Writer: os.Stdout}
	} else {
		write = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize, // megabytes
			MaxBackups: config.MaxBackup,
			MaxAge:     config.MaxAge, //days
			Compress:   config.Compress,
		}
	}

	logger := zerolog.New(write).With().Timestamp().Logger()

	return &logger
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, p, "", "   ")

	if err != nil {
		return w.Writer.Write(p)
	}

	return w.Writer.Write(prettyJSON.Bytes())
}
