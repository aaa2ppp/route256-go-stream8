package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/config"
	"route256/cart/internal/model"
	"route256/cart/pkg/http/middleware"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProduct_GetInfo(t *testing.T) {
	const requestTimeout = 100 * time.Millisecond
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		request       model.SKU
		token         string
		expectedResp  model.GetProductResponse
		expectedError error
	}{
		{
			name: "successful product info",
			handler: func(w http.ResponseWriter, r *http.Request) {
				var req productGetInfoRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				require.Equal(t, "token123", req.Token)
				require.Equal(t, model.SKU(456), req.SKU)

				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(productGetInfoResponse{
					Name:  "Product 1",
					Price: 1000,
				})
			},
			request: 456,
			token:   "token123",
			expectedResp: model.GetProductResponse{
				Name:  "Product 1",
				Price: 1000,
			},
			expectedError: nil,
		},
		{
			name: "request error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(requestTimeout * 2)
				w.WriteHeader(http.StatusOK)
			},
			request:       456,
			token:         "token123",
			expectedResp:  model.GetProductResponse{},
			expectedError: model.ErrInternalError,
		},
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			request:       456,
			token:         "token123",
			expectedResp:  model.GetProductResponse{},
			expectedError: model.ErrNotFound,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			request:       456,
			token:         "token123",
			expectedResp:  model.GetProductResponse{},
			expectedError: model.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			cfg := &config.HTTPProductClient{
				BaseURL:        ts.URL,
				RequestTimeout: requestTimeout,
				GetEndpoint:    "/product/info",
			}
			client := NewProduct(cfg)

			ctx := context.Background()
			if tt.token != "" {
				ctx = middleware.ContextWithAuthToken(ctx, tt.token)
			}
			resp, err := client.GetInfo(ctx, tt.request)

			require.Equal(t, tt.expectedResp, resp)
			require.True(t, errors.Is(err, tt.expectedError))
		})
	}
}
