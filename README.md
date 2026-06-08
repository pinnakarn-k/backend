# K Backend

Backend service built with Go and Fiber.

## Quick Start

```powershell
Copy-Item .env.example .env
go mod tidy
go run ./cmd/app
```

Application will start at:

```text
http://localhost:8080
```

Health check:

```powershell
curl http://localhost:8080/healthz
```

## Prerequisites

Required software:

* Go 1.24+
* Git

Verify installation:

```powershell
go version
git --version
```

## Getting Started

### 1. Clone Repository

```powershell
git clone <repository-url>
cd backend
```

### 2. Create Environment File

Create `.env` from `.env.example`

#### PowerShell

```powershell
Copy-Item .env.example .env
```

#### Command Prompt

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

| Key     | Description           |
| ------- | --------------------- |
| SERVICE | Service name          |
| ENV     | local, dev, uat, prod |
| PORT    | HTTP server port      |

### 3. Install Dependencies

```powershell
go mod tidy
```

### 4. Run Application

```powershell
go run ./cmd/app
```

## Verify Application

Health check:

```http
GET /healthz
GET /readyz
```

Example:

```powershell
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

## Development Workflow

### Add New Module

Example module structure:

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

Register routes in:

```go
internal/app/app.go
```

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

### Success Response

```json
{
  "data": {}
}
```

Example:

```go
return response.Success(c, user)
```

### Validation Error Response

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

### Business Error Response

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

## Middleware

Registered globally:

```text
RequestID
Logger
Recover
```

| Middleware | Responsibility                                   |
| ---------- | ------------------------------------------------ |
| RequestID  | Generate request identifier                      |
| Logger     | HTTP access logging                              |
| Recover    | Recover panic and return standard error response |

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
