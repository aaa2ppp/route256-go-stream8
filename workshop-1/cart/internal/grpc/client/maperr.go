package client

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/pkg/http/middleware"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(ctx context.Context, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		log := middleware.GetLoggerFromContextOrDefault(ctx)
		log.Warn("non-gRPC error detected", "error", err)
		return model.ErrInternalError
	}
	switch st.Code() {
	case codes.NotFound:
		return model.ErrNotFound
	case codes.FailedPrecondition:
		return model.ErrPreconditionFailed
	case codes.Internal:
		return model.ErrInternalError
	}
	log := middleware.GetLoggerFromContextOrDefault(ctx)
	log.Warn("unexpected status code", "code", st.Code())
	return model.ErrInternalError
}
