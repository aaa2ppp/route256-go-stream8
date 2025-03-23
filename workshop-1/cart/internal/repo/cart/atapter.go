package cart

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/repo/cart/queries"
	"route256/cart/internal/service"

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

// Add implements service.CartStorage.
func (a Adapter) Add(ctx context.Context, req model.AddCartItemRequest) error {
	return a.queries.Add(ctx, queries.AddParams{
		UserID: int64(req.UserID),
		Sku:    int32(req.SKU),
		Count:  int32(req.Count),
	})
}

// Clear implements service.CartStorage.
func (a Adapter) Clear(ctx context.Context, userID model.UserID) error {
	return a.queries.Clear(ctx, int64(userID))
}

// Delete implements service.CartStorage.
func (a Adapter) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	return a.queries.Delete(ctx, queries.DeleteParams{
		UserID: int64(req.UserID),
		Sku:    int32(req.SKU),
	})
}

// List implements service.CartStorage.
func (a Adapter) List(ctx context.Context, userID model.UserID) ([]model.CartItem, error) {
	list, err := a.queries.List(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	items := make([]model.CartItem, 0, len(list))
	for _, item := range list {
		items = append(items, model.CartItem{
			SKU:   model.SKU(item.Sku),
			Count: uint16(item.Count),
		})
	}
	return items, nil
}

var _ service.CartStorage = Adapter{}
