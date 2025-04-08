# URL Shortener Service

URL shortening service with automatic expiration, developed in Go.

## Features

- URL shortening with unique codes
- Automatic URL expiration after 24 hours
- Automatic redirection to original URLs
- URL and protocol validation
- Metrics and monitoring with Prometheus and Grafana
- Redis storage
- RESTful API
- Unit tests
- Complete documentation

## Architecture

The project follows the CMD (Command Query Responsibility Segregation) architecture with the following layers:

```
.
├── cmd/                    # Application entry points
│   └── server/            # HTTP server
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and routes
│   ├── domain/           # Entities and business rules
│   ├── repository/       # Persistence layer
│   ├── service/          # Business logic
│   └── metrics/          # Metrics and monitoring
├── pkg/                   # Public reusable code
└── test/                 # Integration tests
```

## API Endpoints

### 1. Shorten URL
```bash
POST /shorten
Content-Type: application/json

{
    "url": "https://www.example.com"
}
```

Response:
```json
{
    "short_url": "http://url.li/Ab3Cd4Ef"
}
```

### 2. Redirect to Original URL
```bash
GET /:shortURL
```

Automatically redirects to the original URL.

### 3. Get URL Information
```bash
GET /info/:shortURL
```

Response:
```json
{
    "short_url": "http://url.li/Ab3Cd4Ef",
    "original_url": "https://www.example.com",
    "expires_at": "2024-02-21T15:04:05Z"
}
```

### 4. Delete URL
```bash
DELETE /:shortURL
```

Response:
```json
{
    "message": "URL deleted successfully"
}
```

### 5. Metrics
```bash
GET /metrics
```
Prometheus endpoint with service metrics.

### 6. Health Check
```bash
GET /health
```
Endpoint for service health verification.

## Available Metrics

### HTTP Metrics
- `http_requests_total`: Total HTTP requests by method, endpoint and status
- `http_request_duration_seconds`: HTTP request duration in seconds

### Service Metrics
- `url_shortening_total`: Total shortened URLs
- `url_redirects_total`: Total redirects
- `active_urls`: Current number of active URLs

## Monitoring

The service includes integration with Prometheus and Grafana for monitoring:

### Prometheus
- Endpoint: `http://localhost:9090`
- Configuration in `prometheus.yml`
- Collects metrics every 15 seconds

### Grafana
- Interface: `http://localhost:3000`
- Default login: admin/admin
- Pre-configured dashboards for:
  - Request rate
  - Average latency
  - Active URLs
  - Total shortened URLs
  - Success rate
  - Redirects

## Requirements

- Go 1.21 or higher
- Docker and Docker Compose
- Redis

## Installation

1. Clone the repository:
```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener
```

2. Install dependencies:
```bash
go mod download
```

3. Start services with Docker Compose:
```bash
docker-compose up -d
```

4. Run the application:
```bash
go run cmd/server/main.go
```

## Tests

Run unit tests:
```bash
go test ./...
```

Run integration tests:
```bash
go test ./test/...
```

### Load Testing with K6

The project includes load testing using K6. To run the tests:

1. Make sure the service is running:
```bash
docker-compose up -d
```

2. Run the load test:
```bash
docker-compose run k6 run /scripts/load-test.js
```

The load test includes:
- Progressive load simulation (50-100 virtual users)
- Testing of all API endpoints
- Performance and error metrics
- Integration with Prometheus for metric visualization

#### Test Configuration
- Total duration: 9 minutes
- Stages:
  - 1 min: Ramp up to 50 users
  - 3 min: Maintain 50 users
  - 1 min: Increase to 100 users
  - 3 min: Maintain 100 users
  - 1 min: Ramp down
- Thresholds:
  - 95% of requests must complete in less than 500ms
  - Error rate must be less than 10%

#### Viewing Results
Load test results are automatically sent to Prometheus and can be viewed in Grafana:
1. Access Grafana (http://localhost:3000)
2. Import the load test dashboard
3. View performance metrics

## Usage Example

1. Shorten a URL:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com"}'
```

2. Get URL information:
```bash
curl http://localhost:8080/info/Ab3Cd4Ef
```

3. Access the shortened URL:
```bash
curl -L http://localhost:8080/Ab3Cd4Ef
```

4. Delete a URL:
```bash
curl -X DELETE http://localhost:8080/Ab3Cd4Ef
```

## Environment Configuration

The project uses environment variables for configuration. Copy the `.env.example` file to `.env` and adjust the variables as needed:

```bash
cp .env.example .env
```

### Environment Variables

- `SERVER_PORT`: HTTP server port (default: 8080)
- `REDIS_HOST`: Redis host (default: localhost)
- `REDIS_PORT`: Redis port (default: 6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `BASE_URL`: Base URL for shortened URLs (default: http://url.li)
- `URL_DURATION`: URL expiration duration (default: 24h)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 