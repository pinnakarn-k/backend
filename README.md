```summary

// ===== Repo layer (เหมือนเดิม ไม่ต้องแก้) =====
type ProductRepo struct {
    ID          int              `db:"id"`
    Name        string           `db:"name"`
    Price       *decimal.Decimal `db:"price"`
    Cost        *decimal.Decimal `db:"cost"`
    CreatedBy   string           `db:"created_by"`
    UpdatedAt   time.Time        `db:"updated_at"`
    DeletedFlag bool             `db:"deleted_flag"`
}

// ===== 1. View/Web response =====
type ProductResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price string `json:"price"` // "1,234.56"
}

func ToProductResponse(r *ProductRepo) *ProductResponse {
    return &ProductResponse{
        ID:    r.ID,
        Name:  r.Name,
        Price: FormatPriceWithComma(r.Price),
    }
}

// ===== 2. Export (ไม่ต้องมี struct ใหม่ ใช้ [][]any ตรงๆ) =====
func ToExportRow(r *ProductRepo) []any {
    return []any{
        r.ID,
        r.Name,
        FormatPriceWithComma(r.Price), // เรียก func เดิม ไม่เขียนซ้ำ
    }
}

func ToExportRows(list []*ProductRepo) [][]any {
    rows := make([][]any, 0, len(list))
    for _, r := range list {
        rows = append(rows, ToExportRow(r))
    }
    return rows
}


// ===== Repo layer: map ตรงตาม column ใน DB =====
type ProductRepo struct {
    ID          int              `db:"id"`
    Name        string           `db:"name"`
    Price       *decimal.Decimal `db:"price"`
    Cost        *decimal.Decimal `db:"cost"`
    CreatedBy   string           `db:"created_by"`
    UpdatedAt   time.Time        `db:"updated_at"`
    DeletedFlag bool             `db:"deleted_flag"`
    // ... อีก 10-20 field ตาม table
}

// ===== Service/Response layer: เอาแค่ที่ต้องพ่นออก =====
type ProductResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price string `json:"price"` // format แล้ว เช่น "1,234.56"
}

func ToProductResponse(r *ProductRepo) *ProductResponse {
    return &ProductResponse{
        ID:    r.ID,
        Name:  r.Name,
        Price: FormatPriceWithComma(r.Price),
    }
}


งั้นผมสรุปนะ 
1.repo คืน struct ตาม ผลลัพธ์ของ sql (ยังไม่ต้องมี tag json)
2.service ก็มาวนลูป และเก็บใน struct ของตัวเอง อันนี้คือส่วนที่จะส่งไปให้ f (ติด tag json)

DB → repo.FindAll() → []PositionRow → service แปลง → []PositionResponse → Handler → JSON
```

```repo
// Query A — ดึงทั้งตาราง
type PositionRow struct {
    RowNo   int
    DataDt  time.Time
    ShrCd   string
    RefType string
    Vol     int
    Price   decimal.Decimal
    Amt     decimal.Decimal
}

// Query B — JOIN กับอีก table ได้ field เพิ่ม
type PositionWithNameRow struct {
    RowNo    int
    ShrCd    string
    ShrName  string  // มาจากอีก table
    RefType  string
    Amt      decimal.Decimal
}

// Query C — ดึงแค่บาง column สำหรับ summary
type PositionSummaryRow struct {
    ShrCd   string
    TotalAmt decimal.Decimal
}
```

```service
// Response Struct
type PositionResponse struct {
    RowNo   int    `json:"row_no"`
    DataDt  string `json:"data_dt"`
    ShrCd   string `json:"shr_cd"`
    RefType string `json:"ref_type"`
    Vol     int    `json:"vol"`
    Price   string `json:"price"`
    Amt     string `json:"amt"`
}

// แปลง RefType
func toRefTypeLabel(refType string) string {
    switch refType {
    case "L":
        return "Long"
    case "S":
        return "Short"
    default:
        return refType // fallback ส่งค่าเดิมถ้าไม่รู้จัก
    }
}

// Mapper
func toPositionResponse(p PositionRow) PositionResponse {
    return PositionResponse{
        RowNo:   p.RowNo,
        DataDt:  p.DataDt.Format("2006-01-02"),
        ShrCd:   p.ShrCd,
        RefType: toRefTypeLabel(p.RefType), // แปลงตรงนี้
        Vol:     p.Vol,
        Price:   p.Price.String(),
        Amt:     p.Amt.String(),
    }
}

// Service
func (s *positionService) GetAll() ([]PositionResponse, error) {
    rows, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    result := make([]PositionResponse, len(rows))
    for i, row := range rows {
        result[i] = toPositionResponse(row)
    }

    return result, nil
}
```

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

```text
ใช่ครับ 🎯

จริง ๆ นี่คือแนวที่ Go backend ส่วนใหญ่ใช้กันเลย

repo
↓ return err

service
↓ return err

handler
↓ return err

Fiber ErrorHandler
↓
log
↓
response json

ตัวอย่างจริง

Repo
func (r *repo) Search(...) error {

	_, err := r.db.QueryContext(...)

	if err != nil {
		return fmt.Errorf(
			"query transaction: %w",
			err,
		)
	}

	return nil
}
Service
func (s *service) Search(...) error {

	err := s.repo.Search(...)

	if err != nil {
		return err
	}

	return nil
}

หรือ

if errors.Is(err, sql.ErrNoRows) {
	return ErrTransactionNotFound
}
Handler
func (h *Handler) Search(
	c *fiber.Ctx,
) error {

	var req SearchRequest

	...

	result, err := h.service.Search(
		c.UserContext(),
		req,
	)

	if err != nil {
		return err
	}

	return response.Success(
		c,
		result,
	)
}
ErrorHandler กลาง
func ErrorHandler(...) fiber.ErrorHandler {

	return func(
		c *fiber.Ctx,
		err error,
	) error {

		log.Error(...)

		return response.Error(
			c,
			err,
		)
	}
}

ข้อดีคือ

คุณไม่ต้องเขียน

if err != nil {
	log.Error(...)
	return response.Error(...)
}

ทุก handler

🤣

ผมชอบเรียกว่า

Business Layer
ไม่รู้เรื่อง HTTP

HTTP Layer
ไม่รู้เรื่อง Business

ErrorHandler
เป็นคนเชื่อม

และที่สำคัญมาก

พออีก 2 เดือนคุณมี

20 handlers
50 services
100 repos

คุณยังมี

Error Log = 1 จุด
Error Response = 1 จุด

ดูแลง่ายมาก

ดังนั้นตอบตรง ๆ

ทุกชั้น return err ไหลมาถึง error handler กลางจบเลย

✅ ใช่ครับ

และผมว่าจาก architecture ที่คุณกำลังทำอยู่

context
validator
apperror
response
middleware

ตัว ErrorHandler กลางนี่แหละ คือชิ้นส่วนที่ขาดอยู่พอดีครับ 🚀
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

```text
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"time"
)

func main() {
	requester := "CUSTOMER_CODE_HERE" // requester = customer code
	application := "XXXXX"            // TODO: ขอค่าจริงจากเจ้าของ API
	key := ""                         // TODO: ใส่ UAT key จริง
	contentType := "application/json"

	timestamp := time.Now().Format("20060102150405") // yyyyMMddHHmmss
	iv := "KS" + timestamp                           // 16 bytes

	printHeaderSet("1) Standard Base64", base64.StdEncoding, timestamp, requester, application, key, iv, contentType)
	printHeaderSet("2) URL Safe Base64 with padding", base64.URLEncoding, timestamp, requester, application, key, iv, contentType)
	printHeaderSet("3) URL Safe Base64 without padding", base64.RawURLEncoding, timestamp, requester, application, key, iv, contentType)

	fmt.Println()
	fmt.Println("=== Debug ===")
	fmt.Println("timestamp:", timestamp)
	fmt.Println("iv:", iv)
	fmt.Println("iv length:", len(iv))
	fmt.Println("key length:", len(key))
}

func printHeaderSet(
	title string,
	encoder *base64.Encoding,
	timestamp, requester, application, key, iv, contentType string,
) {
	pretoken := encoder.EncodeToString([]byte(timestamp))

	token, err := encryptAESCBCPKCS5(timestamp, key, iv, encoder)
	if err != nil {
		fmt.Println("===", title, "===")
		fmt.Println("ERROR:", err)
		fmt.Println()
		return
	}

	fmt.Println("===", title, "===")
	fmt.Println("requester:", requester)
	fmt.Println("application:", application)
	fmt.Println("pretoken:", pretoken)
	fmt.Println("token:", token)
	fmt.Println("Content-Type:", contentType)
	fmt.Println()
}

func encryptAESCBCPKCS5(
	plainText string,
	key string,
	iv string,
	encoder *base64.Encoding,
) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainBytes := []byte(plainText)
	plainBytes = pkcs5Padding(plainBytes, aes.BlockSize)

	cipherText := make([]byte, len(plainBytes))

	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, plainBytes)

	return encoder.EncodeToString(cipherText), nil
}

func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
```