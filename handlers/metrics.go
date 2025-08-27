package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"database/sql"
)

type MetricsResponse struct {
	ActiveUsers    int     `json:"active_users"`
	EventsPerMin   float64 `json:"events_per_min"`
	AvgDuration    float64 `json:"avg_duration"`
	TopElements    []ElementMetric `json:"top_elements"`
}

type ElementMetric struct {
	Element string `json:"element"`
	Count   int    `json:"count"`
}

func MetricsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get active users in last 5 minutes
		var activeUsers int
		fiveMinAgo := time.Now().Add(-5 * time.Minute)
		err := db.QueryRow(`
			SELECT COUNT(DISTINCT user_id) 
			FROM events 
			WHERE timestamp > $1`, fiveMinAgo).Scan(&activeUsers)
		if err != nil {
			http.Error(w, "Error querying active users", http.StatusInternalServerError)
			return
		}

		// Get events per minute
		var eventsPerMin float64
		err = db.QueryRow(`
			SELECT COUNT(*) / 5.0 
			FROM events 
			WHERE timestamp > $1`, fiveMinAgo).Scan(&eventsPerMin)
		if err != nil {
			http.Error(w, "Error querying events per min", http.StatusInternalServerError)
			return
		}

		// Get average duration
		var avgDuration float64
		err = db.QueryRow(`
			SELECT AVG(duration) 
			FROM events 
			WHERE timestamp > $1 AND duration > 0`, fiveMinAgo).Scan(&avgDuration)
		if err != nil {
			http.Error(w, "Error querying avg duration", http.StatusInternalServerError)
			return
		}

		// Get top interactive elements
		rows, err := db.Query(`
			SELECT element, COUNT(*) as count 
			FROM events 
			WHERE timestamp > $1 
			GROUP BY element 
			ORDER BY count DESC 
			LIMIT 5`, fiveMinAgo)
		if err != nil {
			http.Error(w, "Error querying top elements", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var topElements []ElementMetric
		for rows.Next() {
			var em ElementMetric
			if err := rows.Scan(&em.Element, &em.Count); err != nil {
				continue
			}
			topElements = append(topElements, em)
		}

		response := MetricsResponse{
			ActiveUsers:  activeUsers,
			EventsPerMin: eventsPerMin,
			AvgDuration:  avgDuration,
			TopElements:  topElements,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}