package handlers

import (
	"encoding/json"
	"net/http"

	"pingback/internal/services"
	"pingback/pkg/utils"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	store *services.Store
}

func NewEventHandler(store *services.Store) *EventHandler {
	return &EventHandler{store: store}
}

func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		utils.Error("Failed to encode events: " + err.Error())
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get events")
		return
	}
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	event, found := h.store.GetByID(id)
	if !found {
		utils.Error("Event not found: " + id)
		utils.WriteJSONError(w, http.StatusNotFound, "Event not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		utils.Error("Failed to encode event: " + err.Error())
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get event")
	}
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if deleted := h.store.Delete(id); !deleted {
		utils.Error("Failed to delete event or event not found: " + id)
		utils.WriteJSONError(w, http.StatusNotFound, "Event not found or already deleted")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Event deleted successfully"})
}
