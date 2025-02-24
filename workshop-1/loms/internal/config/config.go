package config

import (
	"log/slog"
	"time"
)

type Logger struct {
	Level     slog.Level
	PlainText bool
}

type HTTPServer struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type Config struct {
	Logger     *Logger
	HTTPServer *HTTPServer
}

func Load() (Config, error) {
	return Config{
		Logger: &Logger{
			Level:     slog.LevelDebug,
			PlainText: true,
		},
		HTTPServer: &HTTPServer{
			Addr:            ":8080",
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
	}, nil
}
