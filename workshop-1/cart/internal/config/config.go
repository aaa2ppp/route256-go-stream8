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

type HTTPProductClient struct {
	BaseURL        string
	GetEndpoint    string
	RequestTimeout time.Duration
}

type HTTPLOMSClient struct {
	BaseURL              string
	CreateOrderEndpoint  string
	GetStockInfoEndpoint string
	RequestTimeout       time.Duration
}

type GRPCLOMSClient struct {
	Addr string
}

type Config struct {
	Logger            *Logger
	HTTPServer        *HTTPServer
	HTTPLOMSClient    *HTTPLOMSClient
	HTTPProductClient *HTTPProductClient
	GRPCLOMSClient    *GRPCLOMSClient
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
		HTTPLOMSClient: &HTTPLOMSClient{
			BaseURL:              "http://loms:8080",
			CreateOrderEndpoint:  "/order/create",
			GetStockInfoEndpoint: "/stock/info",
			RequestTimeout:       10 * time.Second,
		},
		HTTPProductClient: &HTTPProductClient{
			BaseURL:        "http://route256.pavl.uk:8080",
			GetEndpoint:    "/get_product",
			RequestTimeout: 10 * time.Second,
		},
		GRPCLOMSClient: &GRPCLOMSClient{
			Addr: "loms:50051",
		},
	}, nil
}
