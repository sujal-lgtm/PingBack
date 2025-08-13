package handlers

import (
	"encoding/json"
	"net/http"

	"pingback/internal/models"
	"pingback/internal/services"
	"pingback/pkg/utils"
)

type CaptureHandler struct {
	store *services.Store
}

func NewCaptureHandler(store *services.Store) *CaptureHandler {
	return &CaptureHandler{store: store}
}

func (h *CaptureHandler) Capture(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.Error("Invalid JSON in capture: " + err.Error())
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if event.Source == "" {
		utils.Error("Empty source in capture request")
		utils.WriteJSONError(w, http.StatusBadRequest, "Source cannot be empty")
		return
	}

	if event.Payload == "" {
		utils.Error("Empty payload in capture request")
		utils.WriteJSONError(w, http.StatusBadRequest, "Payload cannot be empty")
		return
	}

	event.ID = utils.GenerateID("evt")
	h.store.Save(event)
	utils.Info("Captured event ID: " + event.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"event_id": event.ID})
}
