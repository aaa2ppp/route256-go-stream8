package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPong(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	pong(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "pong" {
		t.Errorf("expected body 'pong', got '%s'", w.Body.String())
	}
}
