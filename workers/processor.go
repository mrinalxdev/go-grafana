package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"analytics-engine/models"

	"github.com/redis/go-redis/v9"
)

func ProcessEvents(rdb *redis.Client, db *sql.DB) {
	for {
		// Read from Redis stream
		result, err := rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"events:live", "0"},
			Block:   0,
		}).Result()
		
		if err != nil {
			log.Printf("Error reading from stream: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range result {
			for _, message := range stream.Messages {
				// Parse the event
				eventJSON := message.Values["event"].(string)
				var event models.Event
				err := json.Unmarshal([]byte(eventJSON), &event)
				if err != nil {
					log.Printf("Error unmarshaling event: %v", err)
					continue
				}

				// Store in PostgreSQL
				_, err = db.Exec(`
					INSERT INTO events (user_id, action, element, duration, timestamp)
					VALUES ($1, $2, $3, $4, $5)`,
					event.UserID, event.Action, event.Element, event.Duration, event.Timestamp)
				
				if err != nil {
					log.Printf("DB insert error: %v", err)
					continue
				}

				// Acknowledge processing of the event
				rdb.XAck(context.Background(), "events:live", "events-group", message.ID)
			}
		}
	}
}