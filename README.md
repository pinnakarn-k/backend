# K Backend

Backend service built with Go and Fiber.

---

# Prerequisites

Required software:

* Go 1.24+
* Git

Verify installation:

```bash
go version
git --version
```

---

# Getting Started

## 1. Clone Repository

```bash
git clone <repository-url>
cd backend
```

## 2. Create Environment File

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

---

## 3. Install Dependencies

```bash
go mod tidy
```

---

## 4. Run Application

```bash
go run ./cmd/app
```

Application will start on:

```text
http://localhost:8080
```

---

# Verify Application

Health check:

```http
GET /healthz
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

# Development Workflow

## Add New Module

Example:

```text
internal/
└── user/
    ├── handler.go
    ├── service.go
    ├── repository.go
    ├── route.go
    ├── request.go
    ├── response.go
    └── error.go
```

Register routes in:

```go
internal/app/app.go
```

---

## Request Validation

Example DTO:

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

---

## Success Response

```json
{
  "data": {}
}
```

Example:

```go
return response.Success(c, user)
```

---

## Validation Error Response

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

## Business Error Response

```json
{
  "code": "USER_NOT_FOUND",
  "message": "user not found"
}
```

Example:

```go
return response.Error(c, user.ErrUserNotFound)
```

---

# Middleware

Registered globally:

```go
RequestID
Logger
Recover
```

Responsibilities:

| Middleware | Purpose                               |
| ---------- | ------------------------------------- |
| RequestID  | Generate request identifier           |
| Logger     | HTTP access logging                   |
| Recover    | Recover panic and return 500 response |

---

# Logging

Logs are written to stdout in JSON format.

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

---

# Project Structure

```text
cmd/
└── app/

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
