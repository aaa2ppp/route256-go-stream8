package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"route256/cart/internal/model"
	"strings"
	"testing"
)

func TestCartAddItem(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		token          string
		body           string
		addFunc        func(ctx context.Context, req model.AddCartItemRequest) error
		expectedStatus int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `{"user": 1, "sku": 123, "count": 2}`,
			addFunc: func(ctx context.Context, req model.AddCartItemRequest) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		// {
		// 	name:   "unauthorized",
		// 	method: http.MethodPost,
		// 	token:  "",
		// 	body:   `{"user": 1, "sku": 123, "count": 2}`,
		// 	addFunc: func(ctx context.Context, req model.AddCartItemRequest) error {
		// 		return nil
		// 	},
		// 	expectedStatus: http.StatusUnauthorized,
		// },
		{
			name:   "invalid body",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `invalid-json`,
			addFunc: func(ctx context.Context, req model.AddCartItemRequest) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `{"user": 1, "sku": 123, "count": 2}`,
			addFunc: func(ctx context.Context, req model.AddCartItemRequest) error {
				return errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "precondition failed",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `{"user": 1, "sku": 123, "count": 2}`,
			addFunc: func(ctx context.Context, req model.AddCartItemRequest) error {
				return model.ErrPreconditionFailed
			},
			expectedStatus: http.StatusPreconditionFailed,
		},
	}

	defer func(l *slog.Logger) {
		slog.SetDefault(l)
	}(slog.Default())
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/cart/add", strings.NewReader(tt.body))
			req.Header.Set("X-AuthToken", tt.token)
			w := httptest.NewRecorder()

			handler := CartAddItem(tt.addFunc)
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCartDeleteItem(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		deleteFunc     func(ctx context.Context, req model.DeleteCartItemRequest) error
		expectedStatus int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   `{"user": 1, "sku": 123}`,
			deleteFunc: func(ctx context.Context, req model.DeleteCartItemRequest) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid body",
			method: http.MethodPost,
			body:   `invalid-json`,
			deleteFunc: func(ctx context.Context, req model.DeleteCartItemRequest) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			body:   `{"user": 1, "sku": 123}`,
			deleteFunc: func(ctx context.Context, req model.DeleteCartItemRequest) error {
				return errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/cart/delete", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler := CartDeleteItem(tt.deleteFunc)
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCartList(t *testing.T) {

	// TODO: check response

	tests := []struct {
		name           string
		method         string
		token          string
		body           string
		listFunc       func(ctx context.Context, userID model.UserID) (model.CartListResponse, error)
		expectedStatus int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `{"user": 1}`,
			listFunc: func(ctx context.Context, userID model.UserID) (model.CartListResponse, error) {
				return model.CartListResponse{}, nil
			},
			expectedStatus: http.StatusOK,
		},
		// {
		// 	name:   "unauthorized",
		// 	method: http.MethodPost,
		// 	token:  "",
		// 	body:   `{"user": 1}`,
		// 	listFunc: func(ctx context.Context, userID model.UserID) (model.CartListResponse, error) {
		// 		return model.CartListResponse{}, nil
		// 	},
		// 	expectedStatus: http.StatusUnauthorized,
		// },
		{
			name:   "invalid body",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `invalid-json`,
			listFunc: func(ctx context.Context, userID model.UserID) (model.CartListResponse, error) {
				return model.CartListResponse{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			token:  "valid-token",
			body:   `{"user": 1}`,
			listFunc: func(ctx context.Context, userID model.UserID) (model.CartListResponse, error) {
				return model.CartListResponse{}, errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/cart/list", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler := CartList(tt.listFunc)
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCartClear(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		clearFunc      func(ctx context.Context, userID model.UserID) error
		expectedStatus int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   `{"user": 1}`,
			clearFunc: func(ctx context.Context, userID model.UserID) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid body",
			method: http.MethodPost,
			body:   `invalid-json`,
			clearFunc: func(ctx context.Context, userID model.UserID) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			body:   `{"user": 1}`,
			clearFunc: func(ctx context.Context, userID model.UserID) error {
				return errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/cart/clear", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler := CartClear(tt.clearFunc)
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestCartCheckout(t *testing.T) {

	// TODO: check response

	tests := []struct {
		name           string
		method         string
		body           string
		checkoutFunc   func(ctx context.Context, userID model.UserID) (model.OrderID, error)
		expectedStatus int
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   `{"user": 1}`,
			checkoutFunc: func(ctx context.Context, userID model.UserID) (model.OrderID, error) {
				return 123, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "invalid body",
			method: http.MethodPost,
			body:   `invalid-json`,
			checkoutFunc: func(ctx context.Context, userID model.UserID) (model.OrderID, error) {
				return 0, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "internal error",
			method: http.MethodPost,
			body:   `{"user": 1}`,
			checkoutFunc: func(ctx context.Context, userID model.UserID) (model.OrderID, error) {
				return 0, errors.New("internal error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/cart/checkout", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler := CartCheckout(tt.checkoutFunc)
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
