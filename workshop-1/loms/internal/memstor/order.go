package memstor

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/service"
	"slices"
	"sync"
)

type Order struct {
	mu      sync.RWMutex
	idCount model.OrderID
	items   map[model.OrderID]model.Order
}

func NewOrder() *Order {
	return &Order{
		items: map[model.OrderID]model.Order{},
	}
}

func (o *Order) newOrderID() model.OrderID {
	o.idCount++
	return o.idCount
}

// Create implements service.OrderStorage.
func (o *Order) Create(_ context.Context, req model.CreateOrderRequest) (model.CreateOrderResponse, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	orderID := o.newOrderID()
	status := model.OrderStatusNew

	o.items[orderID] = model.Order{
		Status: status,
		UserID: req.UserID,
		Items:  slices.Clone(req.Items),
	}

	return model.CreateOrderResponse{
		OrderID: orderID,
		Status:  status,
	}, nil
}

// GetByID implements service.OrderStorage.
func (o *Order) GetByID(_ context.Context, orderID model.OrderID) (model.Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, exists := o.items[orderID]
	if !exists {
		return model.Order{}, model.ErrNotFound
	}

	return order, nil
}

// SetStatus implements service.OrderStorage.
func (o *Order) SetStatus(_ context.Context, req model.SetOrderStatusRequest) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	order, exists := o.items[req.OrderID]
	if !exists {
		return model.ErrNotFound
	}

	order.Status = req.Status
	o.items[req.OrderID] = order
	return nil
}

var _ service.OrderStorage = &Order{}
