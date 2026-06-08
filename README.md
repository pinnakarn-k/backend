# K Backend

Backend service built with Go and Fiber.

## Prerequisites

Required software:

* Go 1.24+
* Git

Verify installation:

```bash
go version
git --version
```

## Quick Start

Install dependencies:

```bash
go mod tidy
```

Run application:

```bash
go run ./cmd/app
```

Application will start at:

```text
http://localhost:8080
```

## Getting Started

### Clone Repository

```bash
git clone <repository-url>
cd backend
```

### Create Environment File

PowerShell

```powershell
Copy-Item .env.example .env
```

Command Prompt

```cmd
copy .env.example .env
```

Example:

```env
SERVICE=k-backend
ENV=local
PORT=8080
```

Environment values:

| Key     | Description                                 |
| ------- | ------------------------------------------- |
| SERVICE | Service name                                |
| ENV     | Runtime environment (local, dev, uat, prod) |
| PORT    | HTTP server port                            |

### Install Dependencies

```bash
go mod tidy
```

### Run Application

```bash
go run ./cmd/app
```

## Verify Application

Health check endpoints:

```http
GET /healthz
GET /readyz
```

Example:

```bash
curl http://localhost:8080/healthz
```

Expected response:

```json
{
  "data": {
    "status": "ok"
  }
}
```

---

## Health Check

### Liveness Probe

Endpoint:

```http
GET /healthz
```

Purpose:

```text
Check whether the application process is alive.
```

This endpoint should not depend on external systems such as:

* Database
* Redis
* External APIs

Expected response:

```json
{
  "data": {
    "status": "ok"
  }
}
```

---

### Readiness Probe

Endpoint:

```http
GET /readyz
```

Purpose:

```text
Check whether the application is ready to receive traffic.
```

This endpoint may validate required dependencies such as:

* Database
* Redis
* Elasticsearch
* External Services

Expected response:

```json
{
  "data": {
    "status": "ready"
  }
}
```

If a required dependency is unavailable:

```json
{
  "code": "SERVICE_UNAVAILABLE",
  "message": "service unavailable"
}
```

---

### Current Implementation

At the moment:

```text
/healthz = process health check
/readyz  = dummy readiness check
```

Future dependencies should be added to:

```go
health.Repository.Ready()
```

---

# Development Reference

## Module Structure

Every business module should follow this structure:

```text
internal/
└── user/
    ├── handler.go
    ├── service.go
    ├── repository.go
    ├── route.go
    ├── dto.go
    └── error.go
```

Responsibilities:

| File          | Responsibility         |
| ------------- | ---------------------- |
| handler.go    | HTTP layer             |
| service.go    | Business logic         |
| repository.go | Data access            |
| route.go      | Route registration     |
| dto.go        | Request / Response DTO |
| error.go      | Business errors        |

---

## Request Validation

DTO example:

```go
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}
```

Validate request:

```go
if err := validator.Validate(req); err != nil {
	return response.Error(c, err)
}
```

Validation response:

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

---

## Success Response

Response:

```json
{
  "data": {
    "id": 1,
    "name": "John"
  }
}
```

Usage:

```go
return response.Success(c, user)
```

---

## Success Response With Pagination

Response:

```json
{
  "data": [
    {
      "id": 1,
      "name": "John"
    }
  ],
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

Usage:

```go
return response.SuccessWithPagination(
	c,
	users,
	response.Pagination{
		Page:       1,
		PerPage:    20,
		Total:      100,
		TotalPages: 5,
	},
)
```

---

## Business Error

Create business errors inside the module.

Example:

```go
package user

import (
	"net/http"

	"backend/internal/apperror"
)

var (
	ErrUserNotFound = apperror.New(
		http.StatusNotFound,
		"USER_NOT_FOUND",
		"user not found",
	)

	ErrUserAlreadyExists = apperror.New(
		http.StatusConflict,
		"USER_ALREADY_EXISTS",
		"user already exists",
	)
)
```

Usage:

```go
return response.Error(
	c,
	user.ErrUserNotFound,
)
```

Response:

```json
{
  "code": "USER_NOT_FOUND",
  "message": "user not found"
}
```

---

## Common Errors

Available common errors:

```go
apperror.ErrBadRequest
apperror.ErrUnauthorized
apperror.ErrForbidden
apperror.ErrInternalServer
apperror.ErrServiceUnavailable
```

Usage:

```go
return response.Error(
	c,
	apperror.ErrUnauthorized,
)
```

---

## Middleware

Registered globally:

```text
RequestID
Logger
Recover
```

### RequestID

Generate a unique request identifier for every request.

Header:

```http
X-Request-Id
```

### Logger

Log every HTTP request.

Logged fields:

```text
requestId
method
path
status
latencyMs
ip
service
env
```

### Recover

Recover panic and return a standard error response.

Response:

```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "internal server error"
}
```

---

## Logging

Logs are written to stdout in JSON format.

Example:

```json
{
  "service": "k-backend",
  "env": "local",
  "requestId": "123",
  "method": "GET",
  "path": "/healthz",
  "status": 200,
  "latencyMs": 5
}
```

---

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
