# Quotum
Rate limiting backend server implemented in Go with configurable metrics and storage, supporting multiple algorithms. Built for scalability, efficiency, and easy integration into modern infrastructure.

## Features
- Per-user rate limiting via API
- Well-structured and clean codebase, designed for easy debugging and extensibility
- Supports multiple rate limiting algorithms
  - Fixed window
  - Sliding window
  - Token bucket
  - Sliding bucket
- In-memory and storage based options
- API key authentication
- Endpoints for health check and management
- Dockerized for easy deployment

## Getting started
**Prerequisites:**
- Go 1.22+

**Clone the repository**
```text
git clone https://github.com/vuphan121/quotum.git
```

**Run local**

***Add .env file***
```text
RATE=5
INTERVAL=1s
STORE=memory

LOGGING=true

API_KEY=abcde
```
***Run command***
```text
go mod tidy
go run main.go
```

**Run with docker**
```text
docker-compose up --build
```

## Authentication
**Request Header**
```text
Authorization: Bearer <api-key>
```
**Example CURL request**
```text
curl -X POST http://localhost:8080/request \
  -H "Authorization: Bearer <api-key>" \
  -H "Content-Type: application/json" \
  -d '{"userID": "user_1"}'
```


## API Endpoints

### `POST /request/{userID}`
Registers a request and returns whether it's allowed.

**Response**
```json
{
  "allowed": true
} 
```
or
```json
{
  "allowed": false
} 
```

If `allowed: false`, the user has exceeded the rate limit. To check when the user will be unblocked, use [`GET /status/{userID}`](#2-get-statususerid) or [`GET /banlist`](#3-get-banlist) endpoints.

### `GET /status/{userID}`
Returns the status of the user.

**Response**
```json
{
  "status": "ok"
} 
```
or 
```json
{
  "status": "banned",
  "bannedUntil": "2025-01-01T17:00:00Z"
} 
```

### `GET /banlist`
Returns a list of banned users and when they will be unbanned.

**Response**
```json
{
  "bannedUsers": [
    {
      "userID": "user_1",
      "bannedUntil": "2025-01-01T17:00:00Z"
    },
    {
      "userID": "user_2",
      "bannedUntil": "2025-01-01T18:00:00Z"
    }
  ]
}
```

### `GET /health`
Health check endpoint for monitoring different metrics.

**Response**
```json
{
  "status": "ok",
  "uptimeSeconds": 7200,
  "requestsProcessed": 10000,
  "storage": "redis",
  "storageStatus": "connected"
}
```



## Contributing
Pull requests are welcome. Open an issue if you want to discuss an idea.

## License
This project is licensed under the MIT License.

