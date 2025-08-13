package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"pingback/internal/handlers"
	"pingback/internal/models"
	"pingback/internal/services"

	"github.com/gorilla/mux"
)

func TestListEvents(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewEventHandler(store)

	// Prepopulate with event
	store.Save(models.Event{
		ID:      "evt_1",
		Source:  "test",
		Payload: "{\"foo\":\"bar\"}",
	})

	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	handler.ListEvents(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestGetEvent_Success(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewEventHandler(store)

	store.Save(models.Event{
		ID:      "evt_1",
		Source:  "test",
		Payload: "{\"foo\":\"bar\"}",
	})

	req := httptest.NewRequest("GET", "/events/evt_1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "evt_1"})
	w := httptest.NewRecorder()

	handler.GetEvent(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestGetEvent_NotFound(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewEventHandler(store)

	req := httptest.NewRequest("GET", "/events/nonexistent", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
	w := httptest.NewRecorder()

	handler.GetEvent(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", resp.StatusCode)
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewEventHandler(store)

	store.Save(models.Event{
		ID:      "evt_1",
		Source:  "test",
		Payload: "{\"foo\":\"bar\"}",
	})

	req := httptest.NewRequest("DELETE", "/events/evt_1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "evt_1"})
	w := httptest.NewRecorder()

	handler.DeleteEvent(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	store := services.NewStore()
	handler := handlers.NewEventHandler(store)

	req := httptest.NewRequest("DELETE", "/events/nonexistent", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "nonexistent"})
	w := httptest.NewRecorder()

	handler.DeleteEvent(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", resp.StatusCode)
	}
}
