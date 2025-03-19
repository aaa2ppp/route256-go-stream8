package cart

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/service"
)

type Adapter struct {
	*Queries
}

func NewQueriesAdapter(q *Queries) Adapter {
	return Adapter{q}
}

// Add implements service.CartStorage.
func (a Adapter) Add(ctx context.Context, req model.AddCartItemRequest) error {
	n := len(req.Items)
	user_ids := make([]int64, 0, n)
	skus := make([]int32, 0, n)
	counts := make([]int32, 0, n)
	for _, item := range req.Items {
		user_ids = append(user_ids, int64(req.UserID))
		skus = append(skus, int32(item.SKU))
		counts = append(counts, int32(item.Count))
	}
	return a.Queries.AddArrays(ctx, AddArraysParams{
		Column1: user_ids,
		Column2: skus,
		Column3: counts,
	})
}

// Clear implements service.CartStorage.
func (a Adapter) Clear(ctx context.Context, userID model.UserID) error {
	return a.Queries.Clear(ctx, int64(userID))
}

// Delete implements service.CartStorage.
func (a Adapter) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	return a.Queries.Delete(ctx, DeleteParams{
		UserID: int64(req.UserID),
		Sku:    int32(req.SKU),
	})
}

// List implements service.CartStorage.
func (a Adapter) List(ctx context.Context, userID model.UserID) ([]model.CartItem, error) {
	list, err := a.Queries.List(ctx, int64(userID))
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
