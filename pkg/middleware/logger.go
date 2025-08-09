package middleware

import (
	"net/http"
	"pingback/pkg/utils"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		utils.Info("Started " + r.Method + " " + r.RequestURI)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		utils.Info("Completed " + r.Method + " " + r.RequestURI + " in " + duration.String())
	})
}
