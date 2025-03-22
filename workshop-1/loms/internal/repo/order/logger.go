package order

import (
	"context"
	"log/slog"
	"route256/cart/pkg/http/middleware"
)

func getLogger(ctx context.Context, op string) *slog.Logger {
	return middleware.GetLoggerFromContextOrDefault(ctx).With("op", "repo/order#"+op)
}
