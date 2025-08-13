package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pingback/internal/handlers"
	"pingback/internal/models"
	"pingback/internal/services"
)

func TestReplayHandler_Success(t *testing.T) {
	store := services.NewStore()
	forwarder := &services.Forwarder{}
	replayer := services.NewReplayer(store, forwarder)

	// First save an event
	event := models.Event{
		ID:      "evt_123",
		Source:  "test_source",
		Payload: "{\"foo\":\"bar\"}",
	}
	store.Save(event)

	handler := handlers.NewReplayHandler(replayer)

	payload := `{"event_id":"evt_123","target_url":"https://webhook.site/test"}`
	req := httptest.NewRequest("POST", "/replay", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Replay(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}
}

func TestReplayHandler_MissingEventID(t *testing.T) {
	store := services.NewStore()
	forwarder := &services.Forwarder{}
	replayer := services.NewReplayer(store, forwarder)

	handler := handlers.NewReplayHandler(replayer)

	payload := `{"target_url":"https://webhook.site/test"}`
	req := httptest.NewRequest("POST", "/replay", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Replay(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestReplayHandler_InvalidTargetURL(t *testing.T) {
	store := services.NewStore()
	forwarder := &services.Forwarder{}
	replayer := services.NewReplayer(store, forwarder)

	handler := handlers.NewReplayHandler(replayer)

	payload := `{"event_id":"evt_123","target_url":"invalid-url"}`
	req := httptest.NewRequest("POST", "/replay", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Replay(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestReplayHandler_EventNotFound(t *testing.T) {
	store := services.NewStore()
	forwarder := &services.Forwarder{}
	replayer := services.NewReplayer(store, forwarder)

	handler := handlers.NewReplayHandler(replayer)

	payload := `{"event_id":"nonexistent","target_url":"https://webhook.site/test"}`
	req := httptest.NewRequest("POST", "/replay", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Replay(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status 502 Bad Gateway, got %d", resp.StatusCode)
	}
}
