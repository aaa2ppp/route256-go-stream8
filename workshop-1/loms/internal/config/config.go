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
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GRPCServer struct {
	Addr string
}

type Config struct {
	ShutdownTimeout time.Duration
	Logger          *Logger
	HTTPServer      *HTTPServer
	GRPCServer      *GRPCServer
}

func Load() (Config, error) {
	return Config{
		ShutdownTimeout: 30 * time.Second,
		Logger: &Logger{
			Level:     slog.LevelDebug,
			PlainText: true,
		},
		HTTPServer: &HTTPServer{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		GRPCServer: &GRPCServer{
			Addr: ":50051",
		},
	}, nil
}
