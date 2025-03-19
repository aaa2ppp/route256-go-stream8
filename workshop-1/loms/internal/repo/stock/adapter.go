package stock

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/service"
)

type Adapter struct {
	*Queries
}

// GetBySKU implements storage.
func (a Adapter) GetBySKU(ctx context.Context, sku model.SKU) (_ model.Stock, err error) {
	qresp, err := a.Queries.GetBySKU(ctx, int64(sku))
	if err != nil {
		return model.Stock{}, err
	}
	return model.Stock{
		SKU:      sku,
		Count:    qresp.Count,
		Reserved: qresp.Reserved,
	}, nil
}

// Reserve implements storage.
func (a Adapter) Reserve(ctx context.Context, items []model.OrderItem) (err error) {
	qitems := make([]ReserveParams, 0, len(items))
	for _, item := range items {
		qitems = append(qitems, ReserveParams{
			Sku:      int64(item.SKU),
			Reserved: int64(item.Count)},
		)
	}
	b := a.Queries.Reserve(ctx, qitems)
	b.QueryRow(nil) // TODO: check error for each SKU
	return b.Close()
}

// ReserveCancel implements storage.
func (a Adapter) ReserveCancel(ctx context.Context, items []model.OrderItem) error {
	qitems := make([]ReserveCancelParams, 0, len(items))
	for _, item := range items {
		qitems = append(qitems, ReserveCancelParams{
			Sku:      int64(item.SKU),
			Reserved: int64(item.Count)},
		)
	}
	b := a.Queries.ReserveCancel(ctx, qitems)
	b.QueryRow(nil) // TODO: check error for each SKU
	return b.Close()
}

// ReserveRemove implements storage.
func (a Adapter) ReserveRemove(ctx context.Context, items []model.OrderItem) error {
	qitems := make([]ReserveRemoveParams, 0, len(items))
	for _, item := range items {
		qitems = append(qitems, ReserveRemoveParams{
			Sku:   int64(item.SKU),
			Count: int64(item.Count)},
		)
	}
	b := a.Queries.ReserveRemove(ctx, qitems)
	b.QueryRow(nil) // TODO: check error for each SKU
	return b.Close()
}

var _ service.StockStorage = Adapter{}
