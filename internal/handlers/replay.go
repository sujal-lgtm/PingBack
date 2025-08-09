package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"pingback/internal/services"
	"pingback/pkg/utils"

	"github.com/gorilla/mux"
)

type ReplayRequest struct {
	EventID   string `json:"event_id,omitempty"` // optional for ReplayByID, mandatory for Replay
	TargetURL string `json:"target_url"`
}

type ReplayHandler struct {
	replayer *services.Replayer
}

func NewReplayHandler(replayer *services.Replayer) *ReplayHandler {
	return &ReplayHandler{
		replayer: replayer,
	}
}

// Replay handles POST /replay with JSON body containing event_id and target_url
func (h *ReplayHandler) Replay(w http.ResponseWriter, r *http.Request) {
	var req ReplayRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Invalid JSON in replay: " + err.Error())
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.EventID == "" {
		utils.Error("Replay request missing event_id")
		utils.WriteJSONError(w, http.StatusBadRequest, "event_id is required")
		return
	}

	if req.TargetURL == "" {
		utils.Error("Replay request missing target_url")
		utils.WriteJSONError(w, http.StatusBadRequest, "target_url is required")
		return
	}

	if _, err := url.ParseRequestURI(req.TargetURL); err != nil {
		utils.Error("Invalid target URL while replaying request")
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid target_url")
		return
	}

	utils.Info("Replaying event ID: " + req.EventID + " to " + req.TargetURL)

	if err := h.replayer.ReplayEvent(req.EventID, req.TargetURL); err != nil {
		utils.Error("Failed to replay event: " + err.Error())
		utils.WriteJSONError(w, http.StatusBadGateway, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event replayed successfully!"})
}

// ReplayByID handles POST /replay/{id} with JSON body containing only target_url
func (h *ReplayHandler) ReplayByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["id"]

	if eventID == "" {
		utils.Error("ReplayByID request missing event ID in URL")
		utils.WriteJSONError(w, http.StatusBadRequest, "event ID is required in URL path")
		return
	}

	var req ReplayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error("Invalid JSON in replay by ID: " + err.Error())
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.TargetURL == "" {
		utils.Error("ReplayByID request missing target_url")
		utils.WriteJSONError(w, http.StatusBadRequest, "target_url is required")
		return
	}

	if _, err := url.ParseRequestURI(req.TargetURL); err != nil {
		utils.Error("Invalid target URL while replaying request by ID")
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid target_url")
		return
	}

	utils.Info("Replaying event ID: " + eventID + " to " + req.TargetURL)

	if err := h.replayer.ReplayEvent(eventID, req.TargetURL); err != nil {
		utils.Error("Failed to replay event by ID: " + err.Error())
		utils.WriteJSONError(w, http.StatusBadGateway, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event replayed successfully!"})
}
