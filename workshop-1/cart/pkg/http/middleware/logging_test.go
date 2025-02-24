package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Logging(slog.Default(), handler)
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestLogging_RequestLogging(t *testing.T) {
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, &slog.HandlerOptions{Level: slog.LevelDebug}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := Logging(log, handler)
	middleware.ServeHTTP(w, req)

	logStr := logOutput.String()

	if !strings.Contains(logStr, `http request begin`) {
		t.Error("expected log entry for request begin")
	}

	if !strings.Contains(logStr, `http request end`) {
		t.Error("expected log entry for request end")
	}
}

func TestLogging_LoggerInContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := GetLoggerFromContext(r.Context())
		if log == nil {
			t.Error("expected logger in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := Logging(slog.Default(), handler)
	middleware.ServeHTTP(w, req)
}

func TestLogging_PanicRecovery(t *testing.T) {
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, nil))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := Logging(log, handler)
	middleware.ServeHTTP(w, req)

	logStr := logOutput.String()
	if !strings.Contains(logStr, `*** panic recovered ***`) {
		t.Error("expected log entry for panic recovery")
	}

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestLogging_StatusCodeHook(t *testing.T) {
	var logOutput bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logOutput, &slog.HandlerOptions{Level: slog.LevelDebug}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	middleware := Logging(log, handler)
	middleware.ServeHTTP(w, req)

	logStr := logOutput.String()
	if !strings.Contains(logStr, `statusCode=404`) {
		t.Error("expected log entry for status code 404")
	}
}
