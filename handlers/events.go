package handlers

import (
	// "context"
	"encoding/json"
	"net/http"
	"time"

	"analytics-engine/models"

	"github.com/redis/go-redis/v9"
)

func EventHandler(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.Event
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		// Store in Redis stream
		eventJSON, _ := json.Marshal(event)
		err := rdb.XAdd(r.Context(), &redis.XAddArgs{
			Stream: "events:live",
			Values: map[string]interface{}{"event": string(eventJSON)},
		}).Err()

		if err != nil {
			http.Error(w, "Failed to process event", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Event processed"))
	}
}