package handler

import (
	"encoding/json"
	"learn-ci-cd/response"
	"net/http"
	"time"
)

func HealthHandler(timeUp string, timeZone *time.Location) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		}

		currentTime := time.Now().In(timeZone).Format("2006-01-02 15:04:05")

		response := response.HealthResponse{
			Nama:      "Tunas Bimatara Chrisnanta Budiman",
			NRP:       "5025231999",
			Status:    "UP",
			TimeStamp: currentTime,
			Uptime:    timeUp,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
