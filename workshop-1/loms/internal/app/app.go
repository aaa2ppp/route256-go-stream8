package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256/cart/pkg/http/middleware"
	"route256/loms/internal/config"
	grpcHandler "route256/loms/internal/grpc/handler"
	httpHandler "route256/loms/internal/http/handler"
	"route256/loms/internal/memstor"
	"route256/loms/internal/service"
	"route256/loms/pkg/api/order/v1"
	"route256/loms/pkg/api/stock/v1"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

func Run() int {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("can't load config", "error", err)
		return 1
	}

	setupDefaultLogger(cfg.Logger)

	lomsService := createLOMSService()

	// Создаем менеджер и добавляем серверы
	manager := &serverManager{}
	manager.AddServer(newHTTPServer(cfg.HTTPServer, lomsService))
	manager.AddServer(newGRPCServer(cfg.GRPCServer, lomsService))

	// Запускаем серверы
	manager.StartAll()

	// Ожидаем сигнала или падения серверов
	manager.WaitForShutdown()

	// Завершаем работу серверов
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	manager.StopAll(ctx)

	if manager.failed {
		return 1
	}

	return 0
}

func setupDefaultLogger(cfg *config.Logger) {
	var handler slog.Handler
	if cfg.PlainText {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level})
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: cfg.Level})
	}
	slog.SetDefault(slog.New(handler).With("app", "loms"))
}

func createLOMSService() *service.LOMS {
	// TODO: here using memstore for test only
	return service.NewLOMS(
		memstor.NewOrder(),
		memstor.NewRandomStock(productsForTest),
	)
}

// server представляет интерфейс для управления сервером.
type server interface {
	Start(fail chan<- struct{}) // Запуск сервера
	Stop(ctx context.Context)   // Остановка сервера
}

// serverManager управляет несколькими серверами.
type serverManager struct {
	servers []server
	fail    chan struct{}
	failed  bool
	once    sync.Once
}

// AddServer добавляет сервер в менеджер.
func (m *serverManager) AddServer(server server) {
	m.servers = append(m.servers, server)
}

// StartAll запускает все серверы.
func (m *serverManager) StartAll() {
	m.once.Do(func() {
		m.fail = make(chan struct{}, len(m.servers))
		for _, svr := range m.servers {
			go svr.Start(m.fail)
		}
	})
}

// StopAll останавливает все серверы.
func (m *serverManager) StopAll(ctx context.Context) {
	var wg sync.WaitGroup
	for _, svr := range m.servers {
		wg.Add(1)
		go func(s server) {
			defer wg.Done()
			s.Stop(ctx)
		}(svr)
	}
	wg.Wait()
}

// WaitForShutdown ожидает сигнала завершения или падения серверов.
func (m *serverManager) WaitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		slog.Info("shutdown by signal")
	case <-m.fail:
		slog.Info("shutdown by server failure")
		m.failed = true
	}
}

type httpServer struct {
	server *http.Server
}

func newHTTPServer(cfg *config.HTTPServer, lomsService *service.LOMS) *httpServer {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pong)
	mux.Handle("POST /order/create", httpHandler.CreateOrder(lomsService.CreateOrder))
	mux.Handle("POST /order/info", httpHandler.GetOrderInfo(lomsService.GetOrderInfo))
	mux.Handle("POST /order/pay", httpHandler.PayOrder(lomsService.PayOrder))
	mux.Handle("POST /order/cancel", httpHandler.CancelOrder(lomsService.CancelOrder))
	mux.Handle("POST /stock/info", httpHandler.GetStockInfo(lomsService.GetStockInfo))

	return &httpServer{
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      middleware.Logging(slog.Default(), mux),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func pong(w http.ResponseWriter, h *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		slog.Error("ping: can't write pong", "error", err)
	}
}

func (s *httpServer) Start(fail chan<- struct{}) {
	slog.Info("http server startup", "addr", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("http server failed", "error", err)
		fail <- struct{}{} // Отправляем сигнал о падении
	}
}

func (s *httpServer) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error("failed to stop http server", "error", err)
	}
}

type grpcServer struct {
	server *grpc.Server
	cfg    *config.GRPCServer
}

func newGRPCServer(cfg *config.GRPCServer, lomsService *service.LOMS) *grpcServer {
	server := grpc.NewServer(grpc.UnaryInterceptor(middleware.GRPCLogging(slog.Default())))
	order.RegisterOrderServer(server, grpcHandler.NewOrder(lomsService))
	stock.RegisterStockServer(server, grpcHandler.NewStock(lomsService))

	return &grpcServer{
		server: server,
		cfg:    cfg,
	}
}

func (s *grpcServer) Start(fail chan<- struct{}) {
	slog.Info("grcp server startup", "addr", s.cfg.Addr)
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		slog.Error("grpc server failed", "error", err)
		fail <- struct{}{} // Отправляем сигнал о падении
		return
	}
	if err := s.server.Serve(lis); err != nil {
		slog.Error("grpc server failed", "error", err)
		fail <- struct{}{} // Отправляем сигнал о падении
	}
}

func (s *grpcServer) Stop(ctx context.Context) {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
	case <-stopped:
	}
}
