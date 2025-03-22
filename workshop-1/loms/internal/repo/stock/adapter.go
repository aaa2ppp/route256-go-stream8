package stock

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/repo/stock/queries"
	"route256/loms/internal/service"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	queries.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Adapter struct {
	db      DB
	queries *queries.Queries
}

func New(db DB) *Adapter {
	return &Adapter{
		db:      db,
		queries: queries.New(db),
	}
}

// GetBySKU implements storage.
func (a Adapter) GetBySKU(ctx context.Context, sku model.SKU) (_ model.Stock, err error) {
	qresp, err := a.queries.GetBySKU(ctx, int32(sku))
	if err != nil {
		return model.Stock{}, err
	}
	return model.Stock{
		SKU:       model.SKU(qresp.Sku),
		Available: uint64(qresp.Available),
		Reserved:  uint64(qresp.Reserved),
	}, nil
}

// Reserve implements storage.
func (a Adapter) Reserve(ctx context.Context, items []model.OrderItem) error {
	queryItems := make([]queries.ReserveParams, 0, len(items))
	for _, item := range items {
		queryItems = append(queryItems, queries.ReserveParams{
			Sku:   int32(item.SKU),
			Count: int64(item.Count)},
		)
	}

	err := pgx.BeginFunc(ctx, a.db, func(tx pgx.Tx) error {
		br := a.queries.WithTx(tx).Reserve(ctx, queryItems)

		var queryErr error
		br.QueryRow(func(i int, count int64, err error) {
			if err == queries.ErrBatchAlreadyClosed {
				return
			}
			if err != nil {
				queryErr = err
				br.Close()
				return
			}
			if count == 0 {
				queryErr = model.ErrPreconditionFailed
				br.Close()
				return
			}
		})

		if queryErr != nil {
			return queryErr
		}

		if err := br.Close(); err != nil {
			return err
		}

		return nil
	})

	return err
}

// ReserveCancel implements storage.
func (a Adapter) ReserveCancel(ctx context.Context, items []model.OrderItem) error {

	queryItems := make([]queries.ReserveCancelParams, 0, len(items))
	for _, item := range items {
		queryItems = append(queryItems, queries.ReserveCancelParams{
			Sku:   int32(item.SKU),
			Count: int64(item.Count)},
		)
	}

	return pgx.BeginFunc(ctx, a.db, func(tx pgx.Tx) error {
		br := a.queries.WithTx(tx).ReserveCancel(ctx, queryItems)

		var queryErr error
		br.QueryRow(func(i int, count int64, err error) {
			if err == queries.ErrBatchAlreadyClosed {
				return
			}
			if err != nil {
				queryErr = err
				br.Close()
				return
			}
			if count == 0 {
				queryErr = model.ErrPreconditionFailed
				br.Close()
				return
			}
		})

		if queryErr != nil {
			return queryErr
		}

		if err := br.Close(); err != nil {
			return err
		}

		return nil
	})
}

// ReserveRemove implements storage.
func (a Adapter) ReserveRemove(ctx context.Context, items []model.OrderItem) error {

	queryItems := make([]queries.ReserveRemoveParams, 0, len(items))
	for _, item := range items {
		queryItems = append(queryItems, queries.ReserveRemoveParams{
			Sku:   int32(item.SKU),
			Count: int64(item.Count)},
		)
	}

	return pgx.BeginFunc(ctx, a.db, func(tx pgx.Tx) error {
		br := a.queries.WithTx(tx).ReserveRemove(ctx, queryItems)

		var queryErr error
		br.QueryRow(func(i int, count int64, err error) {
			if err == queries.ErrBatchAlreadyClosed {
				return
			}
			if err != nil {
				queryErr = err
				br.Close()
				return
			}
			if count == 0 {
				queryErr = model.ErrPreconditionFailed
				br.Close()
				return
			}
		})

		if queryErr != nil {
			return queryErr
		}

		if err := br.Close(); err != nil {
			return err
		}

		return nil
	})
}

var _ service.StockStorage = Adapter{}
