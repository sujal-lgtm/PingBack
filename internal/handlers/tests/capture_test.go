package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pingback/internal/handlers"
	"pingback/internal/services"
)

func TestCaptureHandler_Success(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewCaptureHandler(store)

	payload := `{"source":"test_source","payload":"{\"key\":\"value\"}"}`
	req := httptest.NewRequest("POST", "/capture", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Capture(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", resp.StatusCode)
	}
}

func TestCaptureHandler_EmptySource(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewCaptureHandler(store)

	payload := `{"source":"","payload":"{\"key\":\"value\"}"}`
	req := httptest.NewRequest("POST", "/capture", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Capture(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestCaptureHandler_EmptyPayload(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewCaptureHandler(store)

	payload := `{"source":"test_source","payload":""}`
	req := httptest.NewRequest("POST", "/capture", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Capture(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestCaptureHandler_InvalidJSON(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewCaptureHandler(store)

	payload := `{"source":"test_source","payload":{invalid json}`
	req := httptest.NewRequest("POST", "/capture", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Capture(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}
