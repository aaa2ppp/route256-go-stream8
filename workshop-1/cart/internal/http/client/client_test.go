package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientDoRequest(t *testing.T) {
	type Request struct {
		Message string `json:"message"`
	}

	type Response struct {
		Status string `json:"status"`
	}

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		request        interface{}
		expectStatus   int
		expectError    bool
		expectResponse interface{}
	}{
		{
			name: "successful request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(Response{Status: "ok"})
			},
			request:        Request{Message: "test"},
			expectStatus:   http.StatusOK,
			expectResponse: &Response{Status: "ok"},
		},
		{
			name: "server error response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			request:      Request{Message: "test"},
			expectStatus: http.StatusInternalServerError,
			expectError:  false, // Ошибка только если не смогли прочитать ответ
		},
		{
			name: "invalid json response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("invalid json"))
			},
			request:      Request{Message: "test"},
			expectStatus: http.StatusOK,
			expectError:  true,
		},
		{
			name: "request marshaling error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			request:      make(chan int), // Не маршалится в JSON
			expectStatus: 0,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			client := newClient(ts.URL, 1*time.Second)
			var resp Response
			status, err := client.doRequest(context.Background(), "/test", tt.request, &resp)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error %v, got %v", tt.expectError, err)
			}

			if status != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, status)
			}

			if tt.expectResponse != nil {
				expected := tt.expectResponse.(*Response)
				if resp != *expected {
					t.Errorf("expected response %v, got %v", *expected, resp)
				}
			}
		})
	}
}

func TestClientHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", ct)
		}
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("expected Accept: application/json, got %s", accept)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
	}))
	defer ts.Close()

	client := newClient(ts.URL, 1*time.Second)
	status, err := client.doRequest(context.Background(), "/test", struct{}{}, &struct{}{})
	if err != nil || status != http.StatusOK {
		t.Errorf("unexpected error or status: %v, %d", err, status)
	}
}

func TestClientTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := newClient(ts.URL, 100*time.Millisecond)
	_, err := client.doRequest(context.Background(), "/test", struct{}{}, &struct{}{})

	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := newClient(ts.URL, 1*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err := client.doRequest(ctx, "/test", struct{}{}, &struct{}{})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context canceled error, got %v", err)
	}
}

func TestClientNon2xxResponseHandling(t *testing.T) {
	type Response struct {
		Error string
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{Error: "invalid request"})
	}))
	defer ts.Close()

	client := newClient(ts.URL, 1*time.Second)
	var resp Response
	status, err := client.doRequest(context.Background(), "/test", struct{}{}, &resp)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if status != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", status)
	}

	expected := Response{}
	if resp != expected {
		t.Errorf("expected empty response, got %+v", resp)
	}
}

func TestClientRequestPayload(t *testing.T) {
	type Request struct {
		Message string `json:"message"`
	}

	t.Run("valid request body", func(t *testing.T) {
		var (
			receivedMethod  string
			receivedPath    string
			receivedHeaders http.Header
			receivedBody    []byte
		)

		// Создаем тестовый сервер, который будет сохранять параметры запроса
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedMethod = r.Method
			receivedPath = r.URL.Path
			receivedHeaders = r.Header.Clone()

			var err error
			receivedBody, err = io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("error reading request body: %v", err)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		defer ts.Close()

		// Создаем клиент и выполняем запрос
		c := newClient(ts.URL, time.Second)
		req := Request{Message: "test message"}
		var resp interface{}

		status, err := c.doRequest(context.Background(), "/api/v1/test", req, &resp)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Проверяем метод и путь
		if receivedMethod != http.MethodPost {
			t.Errorf("expected method POST, got %s", receivedMethod)
		}
		if receivedPath != "/api/v1/test" {
			t.Errorf("expected path /api/v1/test, got %s", receivedPath)
		}

		// Проверяем заголовки
		if ct := receivedHeaders.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", ct)
		}
		if accept := receivedHeaders.Get("Accept"); accept != "application/json" {
			t.Errorf("expected Accept: application/json, got %s", accept)
		}

		// Проверяем тело запроса
		expectedBody, _ := json.Marshal(req)
		if !bytes.Equal(receivedBody, expectedBody) {
			t.Errorf("request body mismatch:\n expected: %s\n received: %s",
				expectedBody, receivedBody)
		}

		// Проверяем статус ответа
		if status != http.StatusOK {
			t.Errorf("expected status 200, got %d", status)
		}
	})

	t.Run("empty request body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != "null" {
				t.Errorf("expected 'null' body, got %q", body)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		}))
		defer ts.Close()

		c := newClient(ts.URL, time.Second)
		status, err := c.doRequest(context.Background(), "/empty", nil, &struct{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != http.StatusOK {
			t.Errorf("expected status 200, got %d", status)
		}
	})
}
