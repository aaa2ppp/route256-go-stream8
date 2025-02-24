package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/pkg/http/middleware"
	"route256/loms/internal/config"
	"route256/loms/internal/http/handler"
	"route256/loms/internal/memstor"
	"route256/loms/internal/service"
	"syscall"
)

func Run() int {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	setupDefaultLogger(cfg.Logger)

	lomsService := service.NewLOMS(
		memstor.NewOrder(),
		memstor.NewRandomStock(productsForTest), // xxx
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pong)
	mux.Handle("POST /order/create", handler.CreateOrder(lomsService.CreateOrder))
	mux.Handle("POST /order/info", handler.GetOrderInfo(lomsService.GetOrderInfo))
	mux.Handle("POST /order/pay", handler.PayOrder(lomsService.PayOrder))
	mux.Handle("POST /order/cancel", handler.CancelOrder(lomsService.CancelOrder))
	mux.Handle("POST /stock/info", handler.GetStockInfo(lomsService.GetStockInfo))

	httpServer := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      middleware.Logging(slog.Default(), mux),
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	done := make(chan int)
	go func() {
		defer close(done)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		signal := <-c

		slog.Info("shutdown by singnal", "singnal", signal)
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Error("can't shutdown http server", "error", err)
			done <- 1
		}
	}()

	slog.Info("http server startup", "addr", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("http server fail", "error", err)
		return 1
	}

	return <-done
}

func setupDefaultLogger(cfg *config.Logger) {
	var log *slog.Logger
	if cfg.PlainText {
		log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level}))
	}
	slog.SetDefault(log.With("app", "loms"))
}

func pong(w http.ResponseWriter, h *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		slog.Error("ping: can't write pong")
	}
}
