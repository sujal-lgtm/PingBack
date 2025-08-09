package handlers

import (
	"fmt"
	"net/http"
	"pingback/pkg/utils"
)

func HealthCheck(w http.ResponseWriter, r*http.Request){
	utils.Info("Health Check OK")
	fmt.Fprintf(w, "200 OK")
}