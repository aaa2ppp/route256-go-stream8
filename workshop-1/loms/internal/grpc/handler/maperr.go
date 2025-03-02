package handler

import (
	"context"
	"errors"
	"route256/cart/pkg/http/middleware"
	"route256/loms/internal/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrPreconditionFailed):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, model.ErrInternalError):
		return status.Error(codes.Internal, err.Error())
	}
	log := middleware.GetLoggerFromContextOrDefault(ctx)
	log.Warn("unhandled error", "error", err)
	return status.Error(codes.Internal, model.ErrInternalError.Error())
}
