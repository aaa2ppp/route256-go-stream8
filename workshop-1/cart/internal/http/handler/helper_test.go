package handler

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"route256/cart/internal/model"
	"testing"
)

func TestHelper_checkPOSTMethod(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedResult bool
	}{
		{
			name:           "POST method",
			method:         http.MethodPost,
			expectedResult: true,
		},
		{
			name:           "GET method",
			method:         http.MethodGet,
			expectedResult: false,
		},
		{
			name:           "PUT method",
			method:         http.MethodPut,
			expectedResult: false,
		},
		{
			name:           "DELETE method",
			method:         http.MethodDelete,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос с указанным методом
			r := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			// Создаем helper
			h := newHelper(w, r, "TestCheckPOSTMethod")

			// Вызываем метод и проверяем результат
			result := h.checkPOSTMethod()
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}

			// Если метод не POST, проверяем, что ошибка записана
			if tt.method != http.MethodPost {
				resp := w.Result()
				if resp.StatusCode != http.StatusInternalServerError {
					t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
				}
			}
		})
	}
}

func TestHelper_getAuthToken(t *testing.T) {
	tests := []struct {
		name           string
		token          string
		expectedResult string
	}{
		{
			name:           "with token",
			token:          "valid-token",
			expectedResult: "valid-token",
		},
		{
			name:           "empty token",
			token:          "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос с заголовком X-Authtoken
			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			if tt.token != "" {
				r.Header.Set("X-Authtoken", tt.token)
			}

			w := httptest.NewRecorder()

			// Создаем helper
			h := newHelper(w, r, "TestGetAuthToken")

			// Вызываем метод и проверяем результат
			result := h.getAuthToken()
			if result != tt.expectedResult {
				t.Errorf("expected %q, got %q", tt.expectedResult, result)
			}
		})
	}
}

func TestHelper_readBody(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectedBody  []byte
		expectedError bool
	}{
		{
			name:          "successful read",
			body:          `{"key": "value"}`,
			expectedBody:  []byte(`{"key": "value"}`),
			expectedError: false,
		},
		{
			name:          "empty body",
			body:          "",
			expectedBody:  []byte{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h := newHelper(w, r, "testOp")
			h.lg = slog.Default() // Используем стандартный логгер для тестов

			body, ok := h.ReadBody()
			if ok != !tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, !ok)
			}

			if !bytes.Equal(body, tt.expectedBody) {
				t.Errorf("expected body: %s, got: %s", tt.expectedBody, body)
			}
		})
	}
}

type mockValidator struct {
	Valid bool
}

func (m *mockValidator) Validate() error {
	if !m.Valid {
		return errors.New("validation failed")
	}
	return nil
}

type noValidator struct {
	Valid bool
}

func TestHelper_decodeBodyAndValidateRequest(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		req           any
		expectedError bool
	}{
		{
			name:          "successful decode and validate",
			body:          `{"Valid": true}`,
			req:           &mockValidator{Valid: true},
			expectedError: false,
		},
		{
			name:          "invalid JSON",
			body:          `invalid json`,
			req:           &mockValidator{Valid: true},
			expectedError: true,
		},
		{
			name:          "no validator",
			body:          `{"Valid": true}`,
			req:           &noValidator{Valid: true},
			expectedError: true,
		},
		{
			name:          "validation failed",
			body:          `{"Valid": false}`,
			req:           &mockValidator{Valid: false},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h := newHelper(w, r, "testOp")
			h.lg = slog.Default()

			ok := h.decodeBodyAndValidateRequest(tt.req)
			if ok == tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, !ok)
			}
		})
	}
}

func TestHelper_WriteResponse(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		resp           any
		expectedBody   string
		expectedHeader string
	}{
		{
			name:           "successful write",
			code:           http.StatusOK,
			resp:           map[string]string{"key": "value"},
			expectedBody:   `{"key":"value"}`,
			expectedHeader: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			h := newHelper(w, r, "testOp")
			h.lg = slog.Default()

			h.writeResponse(tt.code, tt.resp)

			if w.Code != tt.code {
				t.Errorf("expected status code: %d, got: %d", tt.code, w.Code)
			}

			if w.Header().Get("content-type") != tt.expectedHeader {
				t.Errorf("expected content-type: %s, got: %s", tt.expectedHeader, w.Header().Get("content-type"))
			}

			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body: %s, got: %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHelper_writeError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "httpError",
			err:          &httpError{Code: http.StatusBadRequest, Message: "bad request"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "httpError",
			err:          &httpError{Code: 418, Message: "i'm a teapot"},
			expectedCode: 418,
		},
		{
			name:         "unhandled error",
			err:          errors.New("unexpected error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "not found",
			err:          model.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "precondition failed",
			err:          model.ErrPreconditionFailed,
			expectedCode: http.StatusPreconditionFailed,
		},
		{
			name:         "internal error",
			err:          model.ErrInternalError,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			h := newHelper(w, r, "testOp")
			h.lg = slog.Default()

			h.writeError(tt.err)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code: %d, got: %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestHelper_Log(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := newHelper(w, r, "testOp")

	logger := h.log()
	if logger == nil {
		t.Error("expected logger to be initialized, got nil")
	}
}
