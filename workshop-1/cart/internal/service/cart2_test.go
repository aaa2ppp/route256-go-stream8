package service

import (
	"context"
	"errors"
	"route256/cart/internal/model"
	"slices"
	"testing"
)

type mock2CartStorage struct {
	items []model.CartItem
	err   error
}

// Add implements CartStorage.
func (m *mock2CartStorage) Add(ctx context.Context, req model.AddCartItemRequest) error {
	return m.err
}

// Clear implements CartStorage.
func (m *mock2CartStorage) Clear(ctx context.Context, userID model.UserID) error {
	return m.err
}

// Delete implements CartStorage.
func (m *mock2CartStorage) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	return m.err
}

// List implements CartStorage.
func (m *mock2CartStorage) List(ctx context.Context, userID model.UserID) ([]model.CartItem, error) {
	return slices.Clone(m.items), m.err
}

var _ CartStorage = &mock2CartStorage{}

type mock2OrderStorage struct {
	order model.OrderID
	count uint64
	err   error
}

// CreateOrder implements OrderStorage.
func (m *mock2OrderStorage) CreateOrder(ctx context.Context, req model.OrderCreateRequest) (model.OrderID, error) {
	return m.order, m.err
}

// GetStockInfo implements OrderStorage.
func (m *mock2OrderStorage) GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, err error) {
	return m.count, m.err
}

var _ OrderStorage = &mock2OrderStorage{}

type mock2ProductStorage struct {
	name  string
	price uint32
	err   error
}

// GetInfo implements ProductStorage.
func (m *mock2ProductStorage) GetInfo(ctx context.Context, sku model.SKU) (model.GetProductResponse, error) {
	return model.GetProductResponse{Name: m.name, Price: m.price}, m.err
}

var _ ProductStorage = &mock2ProductStorage{}

func TestCart2_Add(t *testing.T) {
	errOtherError := errors.New("other error")
	_ = errOtherError

	type args struct {
		ctx context.Context
		req model.AddCartItemRequest
	}
	tests := []struct {
		name        string
		args        args
		cartStor    *mock2CartStorage
		orderStor   *mock2OrderStorage
		productStor *mock2ProductStorage
		wantErr     error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: nil},
			orderStor:   &mock2OrderStorage{count: 10, err: nil},
			cartStor:    &mock2CartStorage{err: nil},
		},
		{
			name: "product not found",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: model.ErrNotFound},
			orderStor:   &mock2OrderStorage{count: 10, err: nil},
			cartStor:    &mock2CartStorage{err: nil},
			wantErr:     model.ErrNotFound,
		},
		{
			name: "product other error",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: errOtherError},
			orderStor:   &mock2OrderStorage{count: 10, err: nil},
			cartStor:    &mock2CartStorage{err: nil},
			wantErr:     errOtherError,
		},
		{
			name: "not enough stock",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: nil},
			orderStor:   &mock2OrderStorage{count: 1, err: nil},
			cartStor:    &mock2CartStorage{err: nil},
			wantErr:     model.ErrPreconditionFailed,
		},
		{
			name: "stock not found",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: nil},
			orderStor:   &mock2OrderStorage{count: 10, err: model.ErrNotFound},
			cartStor:    &mock2CartStorage{err: nil},
			wantErr:     model.ErrPreconditionFailed,
		},
		{
			name: "stock other error",
			args: args{
				ctx: context.Background(),
				req: model.AddCartItemRequest{
					UserID: 1,
					SKU:    1,
					Count:  2,
				},
			},
			productStor: &mock2ProductStorage{err: nil},
			orderStor:   &mock2OrderStorage{count: 10, err: errOtherError},
			cartStor:    &mock2CartStorage{err: nil},
			wantErr:     errOtherError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := NewCart(tt.cartStor, tt.orderStor, tt.productStor)
			if err := cart.Add(tt.args.ctx, tt.args.req); !errors.Is(err, tt.wantErr) {
				t.Errorf("Cart.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
