# API Design - Mini Bank Application

## Tech Stack

| Component | Technology |
|-----------|------------|
| Web Framework | Go Echo v4 (Labstack) |
| Database | Oracle DB (via GORM Oracle Driver) |
| Message Broker | Celery with Redis broker |
| Authentication | JWT (HS256, 24h expiry) |
| Language | Go 1.21+ |

**Base URL:** `http://localhost:8080/api/v1`

---

## Authentication

### JWT Bearer Token

```
Authorization: Bearer <access_token>
```

**Claims:**
```go
type JWTClaims struct {
    UserID    uint   `json:"user_id"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    jwt.RegisteredClaims
}
```

---

## Endpoints

### POST /auth/register

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123",
  "full_name": "John Doe",
  "phone_number": "+6281234567890",
  "address": "Jl. Merdeka No. 123, Jakarta"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "account_number": "9881234567890",
    "created_at": "2026-03-25T14:00:00Z"
  }
}
```

**Validation:**
- email: valid format, unique
- password: min 8 chars, 1 uppercase, 1 number
- full_name: min 2, max 100 chars
- phone_number: format +62xxx, unique

---

### POST /auth/login

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "user": {
      "user_id": 1,
      "email": "user@example.com",
      "full_name": "John Doe",
      "account_number": "9881234567890"
    }
  }
}
```

---

### GET /accounts/balance

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "account_number": "9881234567890",
    "balance": 500000.00,
    "currency": "IDR",
    "last_updated": "2026-03-25T14:30:00Z"
  }
}
```

---

### POST /accounts/topup

**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "amount": 1000000.00,
  "reference_id": "TOPUP-20260325-001",
  "payment_method": "bank_transfer",
  "payment_details": {
    "bank_name": "BCA",
    "account_number": "1234567890",
    "sender_name": "John Doe"
  }
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Topup successful",
  "data": {
    "transaction_id": 12,
    "type": "TOPUP",
    "amount": 1000000.00,
    "balance_before": 0.00,
    "balance_after": 1000000.00,
    "reference_id": "TOPUP-20260325-001",
    "celery_task_id": "abc-123-xyz",
    "status": "pending",
    "created_at": "2026-03-25T14:30:00Z"
  }
}
```

**Validation:**
- amount: min 10000, max 100000000 IDR
- reference_id: optional, max 50 chars, unique
- payment_method: bank_transfer | e_wallet

---

### POST /accounts/withdraw

**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "amount": 500000.00,
  "reference_id": "WD-20260325-001",
  "destination": {
    "bank_name": "BCA",
    "account_number": "9876543210",
    "account_name": "Jane Smith"
  }
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "transaction_id": 13,
    "type": "WITHDRAW",
    "amount": 500000.00,
    "fee": 2500.00,
    "total_amount": 502500.00,
    "balance_before": 1000000.00,
    "balance_after": 497500.00,
    "reference_id": "WD-20260325-001",
    "celery_task_id": "def-456-xyz",
    "status": "pending",
    "created_at": "2026-03-25T14:35:00Z"
  }
}
```

**Notes:**
- Fee: 0.5% (min 2500 IDR)

---

### POST /accounts/transfer

**Headers:** `Authorization: Bearer <token>`

**Request:**
```json
{
  "to_account_number": "9880987654321",
  "amount": 250000.00,
  "reference_id": "TRF-20260325-001",
  "note": "Payment for lunch"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "transaction_id": 14,
    "type": "TRANSFER_OUT",
    "amount": 250000.00,
    "fee": 0.00,
    "balance_before": 497500.00,
    "balance_after": 247500.00,
    "to_account_number": "9880987654321",
    "to_account_name": "Jane Doe",
    "reference_id": "TRF-20260325-001",
    "celery_task_id": "ghi-789-xyz",
    "status": "pending",
    "created_at": "2026-03-25T14:40:00Z"
  }
}
```

---

### GET /transactions

**Headers:** `Authorization: Bearer <token>`

**Query params:** `page`, `limit`, `type`, `start_date`, `end_date`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "transactions": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total_items": 45,
      "total_pages": 3
    }
  }
}
```

---

### GET /transactions/:id

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "transaction_id": 14,
    "type": "TRANSFER_OUT",
    "amount": 250000.00,
    "fee": 0.00,
    "balance_before": 497500.00,
    "balance_after": 247500.00,
    "to_account_number": "9880987654321",
    "to_account_name": "Jane Doe",
    "reference_id": "TRF-20260325-001",
    "note": "Payment for lunch",
    "status": "completed",
    "celery_task_id": "ghi-789-xyz",
    "created_at": "2026-03-25T14:40:00Z",
    "updated_at": "2026-03-25T14:40:01Z"
  }
}
```

---

## Error Response Format

```json
{
  "success": false,
  "error": "ERROR_CODE",
  "message": "Human readable message",
  "details": [
    {"field": "email", "message": "Invalid email format"}
  ]
}
```

**Error Codes:**
| Code | HTTP | Description |
|------|------|-------------|
| VALIDATION_ERROR | 400 | Input validation failed |
| UNAUTHORIZED | 401 | Invalid/expired token |
| FORBIDDEN | 403 | Access denied |
| NOT_FOUND | 404 | Resource not found |
| INSUFFICIENT_BALANCE | 400 | Not enough balance |
| DUPLICATE_ENTRY | 409 | Data already exists |
| RATE_LIMITED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |

---

## Celery Integration

### Queue Structure

```
topup        -> process_topup task
withdraw     -> process_withdraw task
transfer     -> process_transfer task
notification -> send_notification task
```

### Processing Flow

1. Request -> validate -> create transaction (status: `pending`)
2. Publish Celery task to Redis broker
3. Worker picks up task -> processes -> updates transaction status
4. Send notification task

### Celery Tasks

```python
@app.task(name='tasks.process_topup')
def process_topup(transaction_id: int) -> dict

@app.task(name='tasks.process_withdraw')
def process_withdraw(transaction_id: int) -> dict

@app.task(name='tasks.process_transfer')
def process_transfer(transaction_id: int) -> dict

@app.task(name='tasks.send_notification')
def send_notification(user_id: int, message: str) -> dict
```

---

## Rate Limiting

- 60 requests/minute per IP (configurable)
- Headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset

---

**Version:** 2.0.0
**Last Updated:** 2026-03-25
**Stack:** Go Echo + Oracle DB + Celery
