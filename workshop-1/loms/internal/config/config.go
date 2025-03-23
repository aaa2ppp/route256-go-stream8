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

type DB struct {
	Addr          string
	Name          string
	User          string
	Password      string
	SSLMode       string
	WaitUpTimeout time.Duration
}

type Config struct {
	ShutdownTimeout time.Duration
	Logger          *Logger
	HTTPServer      *HTTPServer
	GRPCServer      *GRPCServer
	DB              *DB
}

func Load() (Config, error) {
	const required = true
	var ge getenv

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
		DB: &DB{
			Addr:          "db",
			Name:          ge.String("DB_NAME", !required, "loms"),
			User:          ge.String("DB_USER", !required, "loms"),
			Password:      ge.String("DB_PASSWORD", required, ""),
			SSLMode:       "disable",
			WaitUpTimeout: 30 * time.Second,
		},
	}, nil
}
