package memstor

import (
	"context"
	"route256/loms/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Create(t *testing.T) {
	ctx := context.Background()
	storage := NewOrder()

	req := model.CreateOrderRequest{
		UserID: 1,
		Items:  []model.OrderItem{{SKU: 123, Count: 2}},
	}

	resp, err := storage.Create(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, model.OrderID(1), resp.OrderID)
	assert.Equal(t, model.OrderStatusNew, resp.Status)

	order, err := storage.GetByID(ctx, resp.OrderID)
	assert.NoError(t, err)
	assert.Equal(t, req.UserID, order.UserID)
	assert.Equal(t, req.Items, order.Items)
	assert.Equal(t, model.OrderStatusNew, order.Status)
}

func TestOrder_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	storage := NewOrder()

	_, err := storage.GetByID(ctx, 1)
	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestOrder_SetStatus(t *testing.T) {
	ctx := context.Background()
	storage := NewOrder()

	req := model.CreateOrderRequest{
		UserID: 1,
		Items:  []model.OrderItem{{SKU: 123, Count: 2}},
	}

	resp, err := storage.Create(ctx, req)
	assert.NoError(t, err)

	newStatus := model.OrderStatusPayed
	err = storage.SetStatus(ctx, model.SetOrderStatusRequest{
		OrderID: resp.OrderID,
		Status:  newStatus,
	})
	assert.NoError(t, err)

	order, err := storage.GetByID(ctx, resp.OrderID)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, order.Status)
}

func TestOrder_SetStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	storage := NewOrder()

	err := storage.SetStatus(ctx, model.SetOrderStatusRequest{
		OrderID: 1,
		Status:  model.OrderStatusPayed,
	})
	assert.ErrorIs(t, err, model.ErrNotFound)
}
