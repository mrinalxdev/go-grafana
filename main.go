package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// "os"
	// "time"

	"analytics-engine/config"
	"analytics-engine/handlers"
	"analytics-engine/workers"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Connect to PostgreSQL
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize database schema
	err = initDB(db)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Start the event processor worker
	go workers.ProcessEvents(rdb, db)

	// Setup HTTP handlers
	http.Handle("/event", cors(http.HandlerFunc(handlers.EventHandler(rdb))))
	http.Handle("/metrics", cors(http.HandlerFunc(handlers.MetricsHandler(db))))
	http.Handle("/", cors(http.HandlerFunc(serveDashboard)))

	// Start server
	fmt.Println("Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// func initDB(db *sql.DB) error {
// 	// Create events table
// 	_, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS events (
// 			id SERIAL PRIMARY KEY,
// 			user_id VARCHAR(255),
// 			action VARCHAR(50),
// 			element VARCHAR(100),
// 			duration DOUBLE PRECISION,
// 			timestamp TIMESTAMPTZ
// 		)
// 	`)
// 	if err != nil {
// 		return err
// 	}

// 	// Create materialized view
// 	_, err = db.Exec(`
// 		CREATE MATERIALIZED VIEW IF NOT EXISTS user_engagement AS
// 		SELECT
// 			user_id,
// 			COUNT(*) AS total_events,
// 			SUM(duration) AS total_duration,
// 			COUNT(DISTINCT DATE(timestamp)) AS active_days
// 		FROM events
// 		GROUP BY user_id
// 	`)
// 	return err
// }

// func initDB(db *sql.DB) error {
//     // Create events table
//     _, err := db.Exec(`
//         CREATE TABLE IF NOT EXISTS events (
//             id SERIAL PRIMARY KEY,
//             user_id VARCHAR(255),
//             action VARCHAR(50),
//             element VARCHAR(100),
//             duration DOUBLE PRECISION,
//             timestamp TIMESTAMPTZ,
//             INDEX idx_user_id (user_id),
//             INDEX idx_timestamp (timestamp),
//             INDEX idx_action (action)
//         )
//     `)
//     if err != nil {
//         return fmt.Errorf("failed to create events table: %v", err)
//     }

//     // Drop existing materialized view if it exists
//     db.Exec(`DROP MATERIALIZED VIEW IF EXISTS user_engagement`)

//     // Create materialized view
//     _, err = db.Exec(`
//         CREATE MATERIALIZED VIEW user_engagement AS
//         SELECT
//             user_id,
//             COUNT(*) AS total_events,
//             SUM(duration) AS total_duration,
//             COUNT(DISTINCT DATE(timestamp)) AS active_days
//         FROM events
//         GROUP BY user_id
//     `)
//     if err != nil {
//         return fmt.Errorf("failed to create materialized view: %v", err)
//     }

//     // Create refresh function for materialized view
//     _, err = db.Exec(`
//         CREATE OR REPLACE FUNCTION refresh_user_engagement()
//         RETURNS TRIGGER AS $$
//         BEGIN
//             REFRESH MATERIALIZED VIEW user_engagement;
//             RETURN NULL;
//         END;
//         $$ LANGUAGE plpgsql
//     `)
//     if err != nil {
//         log.Printf("Note: Could not create refresh function (non-fatal): %v", err)
//     }

//     log.Println("Database initialized successfully")
//     return nil
// }

func initDB(db *sql.DB) error {
	// Create events table without indexes first
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
            user_id VARCHAR(255),
            action VARCHAR(50),
            element VARCHAR(100),
            duration DOUBLE PRECISION,
            timestamp TIMESTAMPTZ
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create events table: %v", err)
	}

	// Create indexes separately
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp)",
		"CREATE INDEX IF NOT EXISTS idx_events_action ON events(action)",
	}

	for _, indexSQL := range indexes {
		_, err = db.Exec(indexSQL)
		if err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	// Drop existing materialized view if it exists
	db.Exec(`DROP MATERIALIZED VIEW IF EXISTS user_engagement`)

	// Create materialized view
	_, err = db.Exec(`
        CREATE MATERIALIZED VIEW user_engagement AS
        SELECT 
            user_id,
            COUNT(*) AS total_events,
            SUM(duration) AS total_duration,
            COUNT(DISTINCT DATE(timestamp)) AS active_days
        FROM events
        GROUP BY user_id
    `)
	if err != nil {
		return fmt.Errorf("failed to create materialized view: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
