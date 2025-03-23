package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"reflect"
	"route256/cart/internal/config"
	"slices"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getDBConnectString(cfg *config.DB) string {
	q := url.Values{}
	if cfg.SSLMode != "" {
		q.Set("sslmode", cfg.SSLMode)
	}
	u := url.URL{
		Scheme:   "postgres",
		Host:     cfg.Addr,
		User:     url.UserPassword(cfg.User, cfg.Password),
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	return u.String()
}

func waitDBUp(ctx context.Context, dbpool *pgxpool.Pool) error {
	const (
		initialBackoff = 1 * time.Second
		maxBackoff     = 10 * time.Second
	)

	var timer *time.Timer
	backoff := initialBackoff

	for {
		// Выполняем проверку подключения
		err := dbpool.Ping(ctx)
		if err == nil {
			return nil
		}
		if !isRetryableError(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		slog.Debug("connection will be retry", "error", err, "backoff", backoff)
		if timer == nil {
			timer = time.NewTimer(backoff)
		} else {
			timer.Reset(backoff)
		}

		// Ждем следующей попытки или завершения контекста
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			// Экспоненциальный рост с ограничением
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func openDB(ctx context.Context, cfg *config.DB) (_ *pgxpool.Pool, err error) {
	dsn := getDBConnectString(cfg)

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	defer func() {
		if err != nil {
			dbpool.Close()
		}
	}()

	if cfg.WaitUpTimeout > 0 {
		waitCtx, cancel := context.WithTimeout(ctx, cfg.WaitUpTimeout)
		defer cancel()

		waitErr := waitDBUp(waitCtx, dbpool)
		if waitErr == nil {
			return dbpool, nil
		}

		if !errors.Is(waitErr, context.DeadlineExceeded) || ctx.Err() != nil {
			return nil, fmt.Errorf("wait for DB failed: %w", waitErr)
		}

		slog.Debug("wait timeout exceeded, performing final ping")
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("initial ping failed: %w", err)
	}

	return dbpool, nil
}

func isRetryableError(err error) bool {
	// Если контекст завершен, не повторяем
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}

	// Сетевые ошибки
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() || errors.Is(netErr, syscall.ECONNREFUSED) {
			return true
		}
		return false
	}

	// Ошибки PostgreSQL
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Список фатальных ошибок из официальной документации:
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		fatalCodes := []string{
			"28P01", // Invalid password
			"3D000", // Database does not exist
			"42501", // Insufficient privilege
			"42601", // Syntax error
			"42703", // Undefined column
			"42P01", // Undefined table
		}
		return !slices.Contains(fatalCodes, pgErr.Code)
	}

	// По умолчанию считаем ошибку повторимой
	slog.Debug("isRetryableError: unknown error", "error", err, "errorType", reflect.TypeOf(err))
	return true
}
