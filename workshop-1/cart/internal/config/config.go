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

type ProductClient struct {
	BaseURL        string
	GetEndpoint    string
	RequestTimeout time.Duration
}

type OrderClient struct {
	BaseURL              string
	CreateOrderEndpoint  string
	GetStockInfoEndpoint string
	RequestTimeout       time.Duration
}

type Config struct {
	Logger        *Logger
	HTTPServer    *HTTPServer
	OrderClient   *OrderClient
	ProductClient *ProductClient
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
		OrderClient: &OrderClient{
			BaseURL:              "http://loms:8080",
			CreateOrderEndpoint:  "/order/create",
			GetStockInfoEndpoint: "/stock/info",
			RequestTimeout:       10 * time.Second,
		},
		ProductClient: &ProductClient{
			BaseURL:        "http://route256.pavl.uk:8080",
			GetEndpoint:    "/get_product",
			RequestTimeout: 10 * time.Second,
		},
	}, nil
}
