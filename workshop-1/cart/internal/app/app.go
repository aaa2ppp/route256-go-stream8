package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/internal/config"
	grpcClient "route256/cart/internal/grpc/client"
	httpClient "route256/cart/internal/http/client"
	"route256/cart/internal/http/handler"
	cartRepo "route256/cart/internal/repo/cart"
	"route256/cart/internal/service"
	"route256/cart/pkg/http/middleware"
	"syscall"
)

// чтобы именованые пакеты не "убегали"
var (
	_ = (*grpcClient.Product)(nil)
	_ = (*httpClient.Product)(nil)
	_ = (*cartRepo.Adapter)(nil)
)

func Run() int {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	setupDefaultLogger(cfg.Logger)

	dbpool, err := openDB(context.Background(), cfg.DB)
	if err != nil {
		slog.Error(err.Error())
		return 1
	}
	defer dbpool.Close()

	// cartStor := memstor.NewCart()
	cartStor := cartRepo.New(dbpool)

	// lomsClient := httpClient.NewOrder(cfg.HTTPLOMSClient)
	lomsClient, err := grpcClient.NewOrder(cfg.GRPCLOMSClient)
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	// productClient := httpClient.NewProduct(cfg.HTTPProductClient)
	productClient, err := grpcClient.NewProduct(cfg.GRPCProductClient)
	if err != nil {
		slog.Error(err.Error())
		return 1
	}

	cartService := service.NewCart(
		cartStor,
		lomsClient,
		productClient,
	)

	cartMux := http.NewServeMux()
	cartMux.Handle("POST /cart/item/add", handler.CartAddItem(cartService.Add))
	cartMux.Handle("POST /cart/item/delete", handler.CartDeleteItem(cartService.Delete))
	cartMux.Handle("POST /cart/list", handler.CartList(cartService.List))
	cartMux.Handle("POST /cart/clear", handler.CartClear(cartService.Clear))
	cartMux.Handle("POST /cart/checkout", handler.CartCheckout(cartService.Checkout))

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pong)
	mux.Handle("/cart/", middleware.Auth(cartMux))

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
