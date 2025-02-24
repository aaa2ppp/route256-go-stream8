package service

import (
	"context"
	"errors"
	"route256/cart/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCartStorage реализует интерфейс CartStorage
type mockCartStorage struct {
	mock.Mock
}

func (m *mockCartStorage) Add(ctx context.Context, req model.AddCartItemRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockCartStorage) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockCartStorage) List(ctx context.Context, userID model.UserID) ([]model.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.CartItem), args.Error(1)
}

func (m *mockCartStorage) Clear(ctx context.Context, userID model.UserID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// mockOrderStorage реализует интерфейс OrderStorage
type mockOrderStorage struct {
	mock.Mock
}

func (m *mockOrderStorage) CreateOrder(ctx context.Context, req model.OrderCreateRequest) (model.OrderID, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(model.OrderID), args.Error(1)
}

func (m *mockOrderStorage) GetStockInfo(ctx context.Context, sku model.SKU) (uint64, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(uint64), args.Error(1)
}

// mockProductStorage реализует интерфейс ProductStorage
type mockProductStorage struct {
	mock.Mock
}

func (m *mockProductStorage) GetInfo(ctx context.Context, req model.GetProductRequest) (model.GetProductResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(model.GetProductResponse), args.Error(1)
}

func TestCart_Add(t *testing.T) {
	ctx := context.Background()
	req := model.AddCartItemRequest{
		Token: "token",
		Items: []model.CartItem{
			{SKU: 123, Count: 2},
		},
	}

	t.Run("success", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockProductStorage.On("GetInfo", ctx, model.GetProductRequest{Token: "token", SKU: 123}).Return(model.GetProductResponse{}, nil)
		mockOrderStorage.On("GetStockInfo", ctx, model.SKU(123)).Return(uint64(10), nil)
		mockCartStorage.On("Add", ctx, req).Return(nil)

		err := cartService.Add(ctx, req)
		assert.NoError(t, err)

		mockProductStorage.AssertExpectations(t)
		mockOrderStorage.AssertExpectations(t)
		mockCartStorage.AssertExpectations(t)
	})

	t.Run("product not found", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockProductStorage.On("GetInfo", ctx, model.GetProductRequest{Token: "token", SKU: 123}).Return(model.GetProductResponse{}, model.ErrNotFound)

		err := cartService.Add(ctx, req)
		assert.ErrorIs(t, err, model.ErrPreconditionFailed)

		mockProductStorage.AssertExpectations(t)
	})

	t.Run("not enough stock", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockProductStorage.On("GetInfo", ctx, model.GetProductRequest{Token: "token", SKU: 123}).Return(model.GetProductResponse{}, nil)
		mockOrderStorage.On("GetStockInfo", ctx, model.SKU(123)).Return(uint64(1), nil)

		err := cartService.Add(ctx, req)
		assert.ErrorIs(t, err, model.ErrPreconditionFailed)

		mockProductStorage.AssertExpectations(t)
		mockOrderStorage.AssertExpectations(t)
	})
}
func TestCart_Add2(t *testing.T) {
	errUnknownError := errors.New("unknown error")
	tests := []struct {
		name                   string
		ctx                    context.Context
		req                    model.AddCartItemRequest
		mockOrder_GetStockInfo func(*mockOrderStorage, context.Context, model.SKU)
		mockProduct_GetInfo    func(*mockProductStorage, context.Context, model.GetProductRequest)
		mockCart_Add           func(*mockCartStorage, context.Context, model.AddCartItemRequest)
		wantErr                error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, nil)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				m.On("GetStockInfo", ctx, sku).Return(uint64(10), nil)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				m.On("Add", ctx, req).Return(nil)
			},
		},
		{
			name: "product not found",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, model.ErrNotFound)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				//m.On("GetStockInfo", ctx, sku).Return(uint64(10), nil)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				//m.On("Add", ctx, req).Return(nil)
			},
			wantErr: model.ErrPreconditionFailed,
		},
		{
			name: "product unknown error",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, errUnknownError)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				//m.On("GetStockInfo", ctx, sku).Return(uint64(10), nil)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				//m.On("Add", ctx, req).Return(nil)
			},
			wantErr: errUnknownError,
		},
		{
			name: "not enough stock",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, nil)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				m.On("GetStockInfo", ctx, sku).Return(uint64(1), nil)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				// m.On("Add", ctx, req).Return(nil)
			},
			wantErr: model.ErrPreconditionFailed,
		},
		{
			name: "stock not found",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, nil)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				m.On("GetStockInfo", ctx, sku).Return(uint64(10), model.ErrNotFound)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				// m.On("Add", ctx, req).Return(nil)
			},
			wantErr: model.ErrPreconditionFailed,
		},
		{
			name: "stock unknown error",
			ctx:  context.Background(),
			req: model.AddCartItemRequest{
				Token: "token",
				Items: []model.CartItem{{SKU: 123, Count: 2}},
			},
			mockProduct_GetInfo: func(m *mockProductStorage, ctx context.Context, req model.GetProductRequest) {
				m.On("GetInfo", ctx, req).Return(model.GetProductResponse{}, nil)
			},
			mockOrder_GetStockInfo: func(m *mockOrderStorage, ctx context.Context, sku model.SKU) {
				m.On("GetStockInfo", ctx, sku).Return(uint64(10), errUnknownError)
			},
			mockCart_Add: func(m *mockCartStorage, ctx context.Context, req model.AddCartItemRequest) {
				// m.On("Add", ctx, req).Return(nil)
			},
			wantErr: errUnknownError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockProduct = &mockProductStorage{}
				mockOrder   = &mockOrderStorage{}
				mockCart    = &mockCartStorage{}
			)
			tt.mockProduct_GetInfo(mockProduct, tt.ctx, model.GetProductRequest{
				Token: tt.req.Token,
				SKU:   tt.req.Items[0].SKU,
			})
			tt.mockOrder_GetStockInfo(mockOrder, tt.ctx, tt.req.Items[0].SKU)
			tt.mockCart_Add(mockCart, tt.ctx, tt.req)

			cart := NewCart(mockCart, mockOrder, mockProduct)

			err := cart.Add(tt.ctx, tt.req)
			assert.ErrorIs(t, err, tt.wantErr)

			mockCart.AssertExpectations(t)
			mockOrder.AssertExpectations(t)
			mockProduct.AssertExpectations(t)
		})
	}

}

func TestCart_Delete(t *testing.T) {

	ctx := context.Background()
	req := model.DeleteCartItemRequest{
		UserID: 1,
		SKU:    123,
	}

	t.Run("success", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}
		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockCartStorage.On("Delete", ctx, req).Return(nil)

		err := cartService.Delete(ctx, req)
		assert.NoError(t, err)

		mockCartStorage.AssertExpectations(t)
	})

	t.Run("item not found", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}
		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockCartStorage.On("Delete", ctx, req).Return(model.ErrNotFound)

		err := cartService.Delete(ctx, req)
		assert.NoError(t, err)

		mockCartStorage.AssertExpectations(t)
	})

	t.Run("other error", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}
		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockCartStorage.On("Delete", ctx, req).Return(errors.New("some error"))

		err := cartService.Delete(ctx, req)
		v := assert.Error(t, err)
		_ = v

		mockCartStorage.AssertExpectations(t)
	})
}

func TestCart_List(t *testing.T) {
	ctx := context.Background()
	req := model.CartListRequest{
		UserID: 1,
		Token:  "token",
	}

	t.Run("success", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		cartItems := []model.CartItem{
			{SKU: 123, Count: 2},
		}
		mockCartStorage.On("List", ctx, model.UserID(1)).Return(cartItems, nil)
		mockProductStorage.On("GetInfo", ctx, model.GetProductRequest{Token: "token", SKU: 123}).Return(model.GetProductResponse{Name: "Product 123", Price: 100}, nil)

		resp, err := cartService.List(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, uint32(200), resp.TotalPrice)

		mockCartStorage.AssertExpectations(t)
		mockProductStorage.AssertExpectations(t)
	})

	t.Run("cart is empty", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockCartStorage.On("List", ctx, model.UserID(1)).Return([]model.CartItem(nil), model.ErrNotFound)

		resp, err := cartService.List(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 0)

		mockCartStorage.AssertExpectations(t)
	})

	t.Run("product not found", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		cartItems := []model.CartItem{
			{SKU: 123, Count: 2},
		}
		mockCartStorage.On("List", ctx, model.UserID(1)).Return(cartItems, nil)
		mockProductStorage.On("GetInfo", ctx, model.GetProductRequest{Token: "token", SKU: 123}).Return(model.GetProductResponse{}, model.ErrNotFound)
		mockCartStorage.On("Delete", ctx, model.DeleteCartItemRequest{UserID: 1, SKU: 123}).Return(nil)

		resp, err := cartService.List(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 0)

		mockCartStorage.AssertExpectations(t)
		mockProductStorage.AssertExpectations(t)
	})
}

func TestCart_Checkout(t *testing.T) {
	ctx := context.Background()
	userID := model.UserID(1)

	t.Run("success", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		cartItems := []model.CartItem{
			{SKU: 123, Count: 2},
		}
		mockCartStorage.On("List", ctx, userID).Return(cartItems, nil)
		mockOrderStorage.On("CreateOrder", ctx, model.OrderCreateRequest{UserID: userID, Items: cartItems}).Return(model.OrderID(1), nil)
		mockCartStorage.On("Clear", ctx, userID).Return(nil)

		orderID, err := cartService.Checkout(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, model.OrderID(1), orderID)

		mockCartStorage.AssertExpectations(t)
		mockOrderStorage.AssertExpectations(t)
	})

	t.Run("cart is empty", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		mockCartStorage.On("List", ctx, userID).Return([]model.CartItem(nil), model.ErrNotFound)

		_, err := cartService.Checkout(ctx, userID)
		assert.ErrorIs(t, err, model.ErrNotFound)

		mockCartStorage.AssertExpectations(t)
	})

	t.Run("create order failed", func(t *testing.T) {
		mockCartStorage := &mockCartStorage{}
		mockOrderStorage := &mockOrderStorage{}
		mockProductStorage := &mockProductStorage{}

		cartService := NewCart(mockCartStorage, mockOrderStorage, mockProductStorage)

		cartItems := []model.CartItem{
			{SKU: 123, Count: 2},
		}
		mockCartStorage.On("List", ctx, userID).Return(cartItems, nil)
		mockOrderStorage.On("CreateOrder", ctx, model.OrderCreateRequest{UserID: userID, Items: cartItems}).Return(model.OrderID(0), errors.New("create order failed"))

		_, err := cartService.Checkout(ctx, userID)
		assert.Error(t, err)

		mockCartStorage.AssertExpectations(t)
		mockOrderStorage.AssertExpectations(t)
	})
}
