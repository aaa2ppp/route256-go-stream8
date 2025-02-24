package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOrder_CreateOrder(t *testing.T) {
	const requestTimeout = 100 * time.Millisecond
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		request       model.OrderCreateRequest
		expectedID    model.OrderID
		expectedError error
	}{
		{
			name: "successful order creation",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var req createOrderRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				require.Equal(t, model.UserID(123), req.User)
				require.Len(t, req.Items, 2)
				require.Equal(t, model.SKU(456), req.Items[0].SKU)
				require.Equal(t, uint16(2), req.Items[0].Count)
				require.Equal(t, model.SKU(789), req.Items[1].SKU)
				require.Equal(t, uint16(1), req.Items[1].Count)

				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(createOrderResponse{OrderID: 42})
			},
			request: model.OrderCreateRequest{
				UserID: 123,
				Items: []model.OrderItem{
					{SKU: 456, Count: 2},
					{SKU: 789, Count: 1},
				},
			},
			expectedID:    42,
			expectedError: nil,
		},
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(requestTimeout * 2)
				w.WriteHeader(http.StatusOK)
			},
			request: model.OrderCreateRequest{
				UserID: 123,
				Items:  []model.OrderItem{{SKU: 456, Count: 2}},
			},
			expectedID:    0,
			expectedError: model.ErrInternalError,
		},
		{
			name: "precondition failed",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusPreconditionFailed)
			},
			request: model.OrderCreateRequest{
				UserID: 123,
				Items:  []model.OrderItem{{SKU: 456, Count: 2}},
			},
			expectedID:    0,
			expectedError: model.ErrPreconditionFailed,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			request: model.OrderCreateRequest{
				UserID: 123,
				Items:  []model.OrderItem{{SKU: 456, Count: 2}},
			},
			expectedID:    0,
			expectedError: model.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			cfg := &config.OrderClient{
				BaseURL:             ts.URL,
				RequestTimeout:      requestTimeout,
				CreateOrderEndpoint: "/order/create",
			}
			client := NewOrder(cfg)

			orderID, err := client.CreateOrder(context.Background(), tt.request)

			require.Equal(t, tt.expectedID, orderID)
			require.True(t, errors.Is(err, tt.expectedError))
		})
	}
}

func TestOrder_GetStockInfo(t *testing.T) {
	const requestTimeout = 100 * time.Millisecond
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		sku           model.SKU
		expectedCount uint64
		expectedError error
	}{
		{
			name: "successful stock info",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var req getStockInfoRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				require.Equal(t, model.SKU(123), req.SKU)

				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(getStockInfoResponse{Count: 10})
			},
			sku:           123,
			expectedCount: 10,
			expectedError: nil,
		},
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(requestTimeout * 2)
				w.WriteHeader(http.StatusOK)
			},
			sku:           123,
			expectedCount: 0,
			expectedError: model.ErrInternalError,
		},
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			sku:           123,
			expectedCount: 0,
			expectedError: model.ErrNotFound,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			sku:           123,
			expectedCount: 0,
			expectedError: model.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			cfg := &config.OrderClient{
				BaseURL:              ts.URL,
				RequestTimeout:       requestTimeout,
				GetStockInfoEndpoint: "/stock/info",
			}
			client := NewOrder(cfg)

			count, err := client.GetStockInfo(context.Background(), tt.sku)

			require.True(t, errors.Is(err, tt.expectedError))
			require.Equal(t, tt.expectedCount, count)
		})
	}
}
