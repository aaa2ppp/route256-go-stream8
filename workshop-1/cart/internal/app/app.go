package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"route256/cart/internal/config"
	grpcClient "route256/cart/internal/grpc/client"
	httpClient "route256/cart/internal/http/client"
	"route256/cart/internal/http/handler"
	"route256/cart/internal/memstor"
	"route256/cart/internal/service"
	"route256/cart/pkg/http/middleware"
)

func Run() int {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	setupDefaultLogger(cfg.Logger)

	// lomsClient := httpClient.NewOrder(cfg.HTTPLOMSClient)
	lomsClient, err := grpcClient.NewOrder(cfg.GRPCLOMSClient)
	if err != nil {
		log.Fatal(err)
	}

	cartService := service.NewCart(
		memstor.NewCart(),
		lomsClient,
		httpClient.NewProduct(cfg.HTTPProductClient),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pong)
	mux.Handle("POST /cart/item/add", handler.CartAddItem(cartService.Add))
	mux.Handle("POST /cart/item/delete", handler.CartDeleteItem(cartService.Delete))
	mux.Handle("POST /cart/list", handler.CartList(cartService.List))
	mux.Handle("POST /cart/clear", handler.CartClear(cartService.Clear))
	mux.Handle("POST /cart/checkout", handler.CartCheckout(cartService.Checkout))

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
		sig := <-c

		slog.Info("shutdown by signal", "signal", sig)
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

	slog.Info("http server stopped")
	return <-done
}

func setupDefaultLogger(cfg *config.Logger) {
	var handler slog.Handler
	if cfg.PlainText {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level})
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level})
	}
	slog.SetDefault(slog.New(handler).With("app", "cart"))
}

func pong(w http.ResponseWriter, h *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		slog.Error("ping: can't write pong")
	}
}
