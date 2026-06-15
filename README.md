```bash
3. Handler ต้อง return err

ถ้าอยากให้ error log กลางทำงานเต็ม ๆ handler ควรเป็นแบบนี้:

items, err := h.service.Search(c.UserContext(), req)
if err != nil {
	return err
}

ไม่ใช่:

if err != nil {
	return response.Error(c, err)
}

เพราะถ้า handler เรียก response.Error เอง error จะไม่ลอยไปถึง Fiber ErrorHandler

สรุปมาตรฐานใหม่:

handler/service/repo
return err ขึ้นมา

Fiber ErrorHandler
log error กลาง + response.Error

อันนี้ product-grade แบบง่ายแล้วครับ.
```


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

Verify build:

```bash
go build ./...
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
SERVICE=backend
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

### Verify Build

```bash
go build ./...
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
curl http://localhost:8080/readyz
```

Expected responses:

```json
{
  "data": {
    "status": "ok"
  }
}
```

```json
{
  "data": {
    "status": "ready"
  }
}
```

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

This endpoint should not depend on:

* Database
* Redis
* External APIs

Response:

```json
{
  "data": {
    "status": "ok"
  }
}
```

### Readiness Probe

Endpoint:

```http
GET /readyz
```

Purpose:

```text
Check whether the application is ready to receive traffic.
```

This endpoint may validate:

* Database
* Redis
* Elasticsearch
* External Services

Response:

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

### Current Implementation

Current behavior:

```text
/healthz = process health check
/readyz  = dummy readiness check
```

Future dependency checks should be implemented in:

```go
health.Repository.Ready()
```

## Development Reference

### Architecture Flow

```text
HTTP Request
    ↓
Middleware
    ↓
Handler
    ↓
Service
    ↓
Repository
    ↓
Database / External Service
    ↓
Response
```

Layer responsibilities:

| Layer      | Responsibility                 |
| ---------- | ------------------------------ |
| Handler    | HTTP request/response handling |
| Service    | Business logic                 |
| Repository | Data access                    |
| Response   | API response contract          |

Guidelines:

* Handler should not contain business logic.
* Service should not contain HTTP-specific logic.
* Repository should not contain business logic.
* Response formatting should use the shared response package.

### Module Structure

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

### DTO

DTO is used for request and response models.

Example:

```go
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
```

Guidelines:

* Request DTO should contain validation tags.
* Response DTO should be returned to clients.
* Database models should not be exposed directly through APIs.
* Request and Response DTO should be defined in dto.go.

### Request Validation

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

### Success Response

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

### Success Response With Pagination

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

### Business Error

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

### Common Errors

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

### API Versioning Strategy

API versioning is not implemented yet.

If public API contracts are introduced, versioning will be applied using URL path:

```text
/api/v1/*
/api/v2/*
```

### RestAPI
```text
GET /api/v1/hd/master
POST /api/v1/hd/search
GET /api/v1/hd/documents/:documentId
POST /api/v1/hd/documents/:documentId/download
POST /api/v1/hd/send-email
```

### Middleware

Registered globally:

```text
RequestID
Logger
Recover
```

#### RequestID

Generate a unique request identifier for every request.

Header:

```http
X-Request-Id
```

#### Logger

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

#### Recover

Recover panic and return a standard error response.

Response:

```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "internal server error"
}
```

### Logging

Logs are written to stdout in JSON format.

Example:

```json
{
  "service": "backend",
  "env": "local",
  "requestId": "123",
  "method": "GET",
  "path": "/healthz",
  "status": 200,
  "latencyMs": 5
}
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
├── validator/
└── ...
```
