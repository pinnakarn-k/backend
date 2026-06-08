# Backend

Backend service built with Go and Fiber.

## Requirements

* Go 1.24+
* Git

## Environment

Create `.env` from `.env.example`

```env
SERVICE=k-backend
ENV=local
PORT=8080
```

Environment values:

| Key     | Description           |
| ------- | --------------------- |
| SERVICE | Service name          |
| ENV     | local, dev, uat, prod |
| PORT    | HTTP server port      |

## Installation

Install dependencies:

```bash
go mod tidy
```

## Run

Start application:

```bash
go run ./cmd/app
```

Default address:

```text
http://localhost:8080
```

## Health Check

```http
GET /healthz
GET /readyz
```

## Project Structure

```text
cmd/
└── app/
    └── main.go

internal/
├── app/
├── apperror/
├── config/
├── health/
├── logger/
├── middleware/
├── response/
└── validator/
```

## Middleware

### RequestID

Generate request identifier for every request.

Header:

```http
X-Request-Id
```

### Logger

Structured JSON logging using slog.

Logged fields:

* requestId
* method
* path
* status
* latencyMs
* ip
* service
* env

### Recover

Recover panic and return standard error response.

Example:

```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "internal server error"
}
```

## Response Contract

### Success

```json
{
  "data": {}
}
```

### Success With Pagination

```json
{
  "data": [],
  "meta": {
    "pagination": {
      "page": 1,
      "perPage": 20,
      "total": 100,
      "totalPages": 5
    }
  }
}
```

### Error

```json
{
  "code": "VALIDATION_ERROR",
  "message": "validation error",
  "errors": [
    {
      "field": "email",
      "message": "email is required"
    }
  ]
}
```

## Validation

Validation is handled using:

```text
github.com/go-playground/validator/v10
```

Example:

```go
if err := validator.Validate(req); err != nil {
    return response.Error(c, err)
}
```

## Logging

Application logs are written to stdout in JSON format.

Example:

```json
{
  "service": "k-backend",
  "env": "local",
  "requestId": "123",
  "method": "GET",
  "path": "/healthz",
  "status": 200
}
```
