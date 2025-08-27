#  Analytics Engine

A production-ready analytics platform that mimics the user behavior tracking systems used by major tech companies like Netflix and YouTube. This system captures, processes, and visualizes user engagement data in real-time.

## Architecture Overview

The system is built with a microservices architecture:

- **Go Backend**: High-performance API for event ingestion
- **Redis**: Real-time event streaming and buffering
- **PostgreSQL**: Persistent data storage and aggregation
- **Grafana**: Advanced analytics and visualization
- **HTML/TailwindCSS**: Clean, responsive dashboard UI

## Tech Stack

- **Backend**: Go 1.21+
- **Database**: PostgreSQL 15
- **Cache/Streaming**: Redis
- **Visualization**: Grafana
- **Frontend**: HTML5 with TailwindCSS
- **Containerization**: Docker + Docker Compose

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later
- Git

## Quick Start

Clone and setup the project:
```bash
git clone <your-repo>
cd analytics-demo
```

Start the infrastructure:
```bash
docker-compose up -d
```

Initialize the Go application:
```bash
go mod init analytics-demo
go mod tidy
go run main.go
```

Access the applications:
- **Main Dashboard**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)

## Project Structure

```
analytics-demo/
├── main.go                 # Main application entry point
├── go.mod                 # Go module dependencies
├── docker-compose.yml     # Container orchestration
├── index.html            # Dashboard frontend
├── handlers/             # HTTP request handlers
│   ├── events.go         # Event ingestion endpoint
│   └── metrics.go        # Metrics API endpoint
├── models/               # Data structures
│   └── event.go          # Event model definition
├── workers/              # Background processors
│   └── processor.go      # Redis to PostgreSQL worker
├── config/               # Configuration files
│   └── database.go       # Database connection setup
└── postgres-init/        # Database initialization
    └── init.sql          # Schema setup scripts
```

## Key Features

### Real-time Event Processing
- High-throughput event ingestion via HTTP API
- Redis streams for buffering and real-time processing
- Background workers for reliable data persistence

### Advanced Analytics
- Materialized views for fast aggregations
- Real-time metrics endpoint for dashboard updates
- PostgreSQL optimized for time-series data

### Professional Visualization
- Grafana integration for enterprise-grade dashboards
- TailwindCSS for modern, responsive UI
- Interactive demo interface with test event generation

### Scalable Architecture
- Docker containerization for easy deployment
- Microservices design for horizontal scaling
- Production-ready error handling and logging

## API Endpoints

### POST /event

Accepts user engagement events in JSON format:
```json
{
  "user_id": "user_123",
  "action": "play",
  "element": "video_player",
  "duration": 12.5,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /metrics

Returns real-time analytics data:
```json
{
  "active_users": 15,
  "events_per_min": 42.5,
  "avg_duration": 8.2,
  "top_elements": [
    {"element": "video_player", "count": 120},
    {"element": "like_button", "count": 85}
  ]
}
```

## Database Schema

### Events Table
```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255),
    action VARCHAR(50),
    element VARCHAR(100),
    duration DOUBLE PRECISION,
    timestamp TIMESTAMPTZ
);
```

### Materialized Views
Pre-aggregated data for fast query performance:
- User engagement metrics
- Activity heatmaps
- Element interaction rankings

## Monitoring and Analytics

The system provides multiple levels of monitoring:

- **Real-time Web Dashboard**: Live metrics updated every 3 seconds
- **Grafana Analytics**: Professional time-series visualizations
- **Database Insights**: SQL-based custom reporting
- **Event Stream Monitoring**: Redis stream health checks

## Performance Characteristics

- Handles thousands of events per second with Redis buffering
- Sub-second response times for metrics API
- Efficient PostgreSQL queries with proper indexing
- Automatic connection pooling and management

## Development Guide

### Adding New Event Types
1. Extend the Event model in `models/event.go`
2. Update the event handler in `handlers/events.go`
3. Add corresponding database migrations if needed

### Creating New Visualizations
1. Add new SQL queries in `handlers/metrics.go`
2. Update the frontend dashboard in `index.html`
3. Create new Grafana panels as needed

### Scaling the System
1. Add more Go instances behind a load balancer
2. Scale Redis with clustering for higher throughput
3. Use PostgreSQL read replicas for analytics queries
4. Implement Redis persistence for data durability

## Troubleshooting

### Common Issues
- **Port conflicts**: Ensure ports 8080, 3000, 5432, and 6379 are available
- **Database connection errors**: Wait for PostgreSQL to fully initialize
- **Redis connection issues**: Check Docker container status

### Debug Mode
Run without daemon mode to see detailed logs:
```bash
docker-compose up
go run main.go
```

### Reset Everything
```bash
docker-compose down
docker volume rm analytics-demo_postgres_data
docker-compose up -d
```

## Production Considerations

- Implement proper authentication for APIs
- Add SSL/TLS encryption for data in transit
- Set up database backups and monitoring
- Configure Grafana authentication and user management
- Implement rate limiting for event ingestion
- Add comprehensive logging and alerting