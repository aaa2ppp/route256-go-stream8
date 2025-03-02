package middleware

import (
	"cmp"
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"runtime/debug"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loggerContextKey struct{}

// ContextWithLogger добавляет логгер в контекст.
func ContextWithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, log)
}

// GetLoggerFromContext возвращает логгер из контекста или nil, если логгер не найден.
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if v := ctx.Value(loggerContextKey{}); v != nil {
		return v.(*slog.Logger)
	}
	return nil
}

// GetLoggerFromContextOrDefault возвращает логгер из контекста или логгер по умолчанию,
// если логгер не найден.
func GetLoggerFromContextOrDefault(ctx context.Context) *slog.Logger {
	return cmp.Or(GetLoggerFromContext(ctx), slog.Default())
}

// Logging создает middleware для логирования HTTP-запросов.
func Logging(log *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.With("httpReqID", rand.Uint64())

		url := r.URL.String()
		log.Debug("http request begin", "fromAddr", r.RemoteAddr, "method", r.Method, "url", url)

		w = newWriteHeaderHook(w, func(statusCode int) {
			log.Debug("http request end", "statusCode", statusCode, "url", url)
		})

		ctx := ContextWithLogger(r.Context(), log)
		r = r.WithContext(ctx)

		defer func() {
			if p := recover(); p != nil {
				log.Error("*** panic recovered ***", "panic", p, "stack", debug.Stack())
				http.Error(w, "internal error", 500)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

type writeHeaderHook struct {
	http.ResponseWriter
	hook func(statusCode int)
	flag atomic.Bool
}

func newWriteHeaderHook(w http.ResponseWriter, hook func(statusCode int)) *writeHeaderHook {
	return &writeHeaderHook{
		ResponseWriter: w,
		hook:           hook,
	}
}

func (hk *writeHeaderHook) WriteHeader(statusCode int) {
	if !hk.flag.Swap(true) {
		hk.hook(statusCode)
		hk.ResponseWriter.WriteHeader(statusCode)
	}
}

func (hk *writeHeaderHook) Write(b []byte) (int, error) {
	hk.WriteHeader(http.StatusOK)
	return hk.ResponseWriter.Write(b)
}

func GRPCLogging(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		log := log.With("grpcReqID", rand.Uint64())
		log.Debug("grpc request begin", "fromAddr", "???", "method", info.FullMethod, "req", req)

		defer func() {
			if p := recover(); p != nil {
				log.Error("*** panic recovered ***", "panic", p, "stack", debug.Stack())
				err = status.Error(codes.Internal, "internal error")
			}
		}()

		resp, err := handler(ContextWithLogger(ctx, log), req)
		log.Debug("grpc request end", "status", status.Code(err), "resp", resp)

		return resp, err
	}
}
