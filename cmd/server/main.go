package main

import (
	"fmt"
	"net/http"
	"pingback/internal/config"
	"pingback/internal/handlers"
	"pingback/internal/services"
	"pingback/pkg/middleware"
	"pingback/pkg/utils"

	"github.com/gorilla/mux"
)

func main() {

	cfg := config.LoadConfig()
	r := mux.NewRouter()

	r.Use(middleware.Logger)

	store := services.NewStore()
	forwarder := &services.Forwarder{}
	replayer := services.NewReplayer(store, forwarder)

	captureHandler := handlers.NewCaptureHandler(store)
	replayHandler := handlers.NewReplayHandler(replayer)
	eventHandler := handlers.NewEventHandler(store)

	r.HandleFunc("/ping", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/capture", captureHandler.Capture).Methods("POST")
	r.HandleFunc("/replay", replayHandler.Replay).Methods("POST")
	r.HandleFunc("/replay/{id}", replayHandler.ReplayByID).Methods("POST")
	r.HandleFunc("/events", eventHandler.ListEvents).Methods("GET")
	r.HandleFunc("/events/{id}", eventHandler.GetEvent).Methods("GET")
	r.HandleFunc("/events/{id}", eventHandler.DeleteEvent).Methods("DELETE")

	fmt.Printf("Server running on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		utils.Error("Error in starting server")
		fmt.Println("Error starting server", err)
	}
}
