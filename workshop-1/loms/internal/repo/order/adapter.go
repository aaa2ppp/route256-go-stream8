package order

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/service"
)

type Adapter struct {
	*Queries
}

// Create implements storage.
func (a Adapter) Create(ctx context.Context, req model.CreateOrderRequest) (model.OrderID, error) {
	orderID, err := a.Queries.Create(ctx, int64(req.UserID))
	if err != nil {
		return 0, err
	}

	n := len(req.Items)
	orders := make([]int64, 0, n)
	skus := make([]int64, 0, n)
	counts := make([]int32, 0, n)
	for _, item := range req.Items {
		orders = append(orders, orderID)
		skus = append(skus, int64(item.SKU))
		counts = append(counts, int32(item.Count))
	}

	if err := a.Queries.AddItems(ctx, AddItemsParams{
		Column1: orders,
		Column2: skus,
		Column3: counts,
	}); err != nil {
		return 0, err
	}

	return model.OrderID(orderID), nil
}

// GetByID implements storage.
func (a Adapter) GetByID(ctx context.Context, orderID model.OrderID) (model.Order, error) {
	qresp, err := a.Queries.GetByID(ctx, int64(orderID))
	if err != nil {
		return model.Order{}, err
	}
	if len(qresp) == 0 {
		return model.Order{}, model.ErrNotFound
	}
	status, err := model.ParseOrderStatus(string(qresp[0].Status))
	if err != nil {
		return model.Order{}, err
	}

	resp := model.Order{
		UserID: model.UserID(qresp[0].UserID),
		Status: status,
		Items:  make([]model.OrderItem, 0, len(qresp)),
	}

	for _, item := range qresp {
		resp.Items = append(resp.Items, model.OrderItem{
			SKU:   model.SKU(item.Sku),
			Count: int(item.Count),
		})
	}

	return resp, nil
}

// SetStatus implements storage.
func (a Adapter) SetStatus(ctx context.Context, req model.SetOrderStatusRequest) error {
	return a.Queries.SetStatus(ctx, SetStatusParams{
		OrderID: int64(req.OrderID),
		Status:  OrderStatus(req.Status.String()),
	})
}

var _ service.OrderStorage = Adapter{}

type StockAdapter struct {
	*Queries
}
