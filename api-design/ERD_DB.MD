# ERD & Database Schema - Mini Bank Application

## Tech Stack: Oracle DB (GORM Oracle Driver)

**Database:** Oracle XE 21c+
**ORM:** GORM with Oracle driver
**Connection:** `gorm.io/driver/oracle`

---

## Tables

1. **users** - User profiles
2. **accounts** - Bank accounts
3. **transactions** - Transaction history

---

## 1. Table: USERS

```go
type User struct {
    ID           uint           `gorm:"primaryKey" json:"user_id"`
    Email        string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
    PasswordHash string         `gorm:"size:255;not null" json:"-"`
    FullName     string         `gorm:"size:255;not null" json:"full_name"`
    PhoneNumber  string         `gorm:"uniqueIndex;size:20;not null" json:"phone_number"`
    Address      string         `gorm:"size:500" json:"address,omitempty"`
    Role         string         `gorm:"size:20;default:user;not null" json:"role"`
    IsActive     bool           `gorm:"default:true" json:"is_active"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Oracle DDL:**
```sql
CREATE TABLE users (
    id           NUMBER(20) PRIMARY KEY,
    email        VARCHAR2(255) NOT NULL UNIQUE,
    password_hash VARCHAR2(255) NOT NULL,
    full_name    VARCHAR2(255) NOT NULL,
    phone_number VARCHAR2(20) NOT NULL UNIQUE,
    address      VARCHAR2(500),
    role         VARCHAR2(20) DEFAULT 'user' NOT NULL,
    is_active    NUMBER(1) DEFAULT 1 NOT NULL,
    created_at   TIMESTAMP DEFAULT SYSTIMESTAMP,
    updated_at   TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT chk_user_role CHECK (role IN ('user', 'admin'))
);
```

**Indexes:**
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_role ON users(role);
```

**Sample Row:**
```json
{
  "id": 1,
  "email": "john.doe@example.com",
  "password_hash": "$2a$10$N9qo8uLOickgx2ZMRZoMye...",
  "full_name": "John Doe",
  "phone_number": "+6281234567890",
  "address": "Jl. Merdeka No. 123, Jakarta",
  "role": "user",
  "is_active": true,
  "created_at": "2026-03-25T14:00:00Z",
  "updated_at": "2026-03-25T14:00:00Z"
}
```

---

## 2. Table: ACCOUNTS

```go
type Account struct {
    ID            uint           `gorm:"primaryKey" json:"account_id"`
    UserID        uint           `gorm:"not null;index" json:"user_id"`
    AccountNumber string         `gorm:"uniqueIndex;size:50;not null" json:"account_number"`
    Balance       float64        `gorm:"type:number(18,2);default:0;not null" json:"balance"`
    Currency      string         `gorm:"size:3;default:IDR" json:"currency"`
    AccountType   string         `gorm:"size:20;default:savings" json:"account_type"`
    IsActive      bool           `gorm:"default:true" json:"is_active"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Oracle DDL:**
```sql
CREATE TABLE accounts (
    id             NUMBER(20) PRIMARY KEY,
    user_id        NUMBER(20) NOT NULL,
    account_number VARCHAR2(50) NOT NULL UNIQUE,
    balance        NUMBER(18,2) DEFAULT 0 NOT NULL,
    currency       VARCHAR2(3) DEFAULT 'IDR',
    account_type   VARCHAR2(20) DEFAULT 'savings',
    is_active      NUMBER(1) DEFAULT 1 NOT NULL,
    created_at     TIMESTAMP DEFAULT SYSTIMESTAMP,
    updated_at     TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT chk_balance_non_negative CHECK (balance >= 0)
);
```

**Indexes:**
```sql
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_number ON accounts(account_number);
CREATE INDEX idx_accounts_active ON accounts(is_active);
```

**Relationships:**
- One-to-One with USERS (user_id -> id)
- One-to-Many with TRANSACTIONS (id -> account_id)

---

## 3. Table: TRANSACTIONS

```go
type Transaction struct {
    ID                 uint           `gorm:"primaryKey" json:"transaction_id"`
    UserID             uint           `gorm:"not null;index" json:"user_id"`
    AccountID          uint           `gorm:"not null;index" json:"account_id"`
    TransactionType    string         `gorm:"size:20;not null" json:"type"`
    Amount             float64        `gorm:"type:number(18,2);not null" json:"amount"`
    Fee                float64        `gorm:"type:number(18,2);default:0" json:"fee"`
    TotalAmount        float64        `gorm:"type:number(18,2);not null" json:"total_amount"`
    BalanceBefore      float64        `gorm:"type:number(18,2);not null" json:"balance_before"`
    BalanceAfter       float64        `gorm:"type:number(18,2);not null" json:"balance_after"`
    Currency           string         `gorm:"size:3;default:IDR" json:"currency"`
    
    // Transfer
    ToAccountID        *uint          `json:"to_account_id,omitempty"`
    ToAccountNumber     string         `gorm:"size:50" json:"to_account_number,omitempty"`
    ToAccountName      string         `gorm:"size:255" json:"to_account_name,omitempty"`
    
    // Withdraw
    DestinationBank    string         `gorm:"size:100" json:"destination_bank,omitempty"`
    DestinationAccount string         `gorm:"size:50" json:"destination_account,omitempty"`
    DestinationName    string         `gorm:"size:255" json:"destination_name,omitempty"`
    
    // Topup
    PaymentMethod      string         `gorm:"size:50" json:"payment_method,omitempty"`
    SenderBank         string         `gorm:"size:100" json:"sender_bank,omitempty"`
    SenderAccount      string         `gorm:"size:50" json:"sender_account,omitempty"`
    SenderName         string         `gorm:"size:255" json:"sender_name,omitempty"`
    
    // Common
    ReferenceID        string         `gorm:"uniqueIndex;size:100" json:"reference_id,omitempty"`
    Note               string         `gorm:"size:500" json:"note,omitempty"`
    Status             string         `gorm:"size:20;default:pending;index" json:"status"`
    CeleryTaskID       string         `gorm:"size:100;index" json:"celery_task_id,omitempty"`
    CreatedAt          time.Time      `gorm:"index" json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
```

**Oracle DDL:**
```sql
CREATE TABLE transactions (
    id                  NUMBER(20) DEFAULT transactions_seq.NEXTVAL PRIMARY KEY,
    user_id             NUMBER(20) NOT NULL,
    account_id          NUMBER(20) NOT NULL,
    transaction_type    VARCHAR2(20) NOT NULL,
    amount              NUMBER(18,2) NOT NULL,
    fee                 NUMBER(18,2) DEFAULT 0,
    total_amount        NUMBER(18,2) NOT NULL,
    balance_before      NUMBER(18,2) NOT NULL,
    balance_after       NUMBER(18,2) NOT NULL,
    currency            VARCHAR2(3) DEFAULT 'IDR',
    to_account_id       NUMBER(20),
    to_account_number   VARCHAR2(50),
    to_account_name     VARCHAR2(255),
    destination_bank    VARCHAR2(100),
    destination_account VARCHAR2(50),
    destination_name    VARCHAR2(255),
    payment_method      VARCHAR2(50),
    sender_bank         VARCHAR2(100),
    sender_account      VARCHAR2(50),
    sender_name         VARCHAR2(255),
    reference_id        VARCHAR2(100) UNIQUE,
    note                VARCHAR2(500),
    status              VARCHAR2(20) DEFAULT 'pending',
    celery_task_id      VARCHAR2(100),
    created_at          TIMESTAMP DEFAULT SYSTIMESTAMP,
    updated_at          TIMESTAMP DEFAULT SYSTIMESTAMP,
    CONSTRAINT chk_tx_type CHECK (transaction_type IN ('TOPUP', 'WITHDRAW', 'TRANSFER_IN', 'TRANSFER_OUT')),
    CONSTRAINT chk_tx_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    CONSTRAINT chk_amount_positive CHECK (amount > 0)
);
```

**Indexes:**
```sql
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_reference ON transactions(reference_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_user_created ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_celery_task ON transactions(celery_task_id);
```

**Sample Rows:**

*Topup:*
```json
{
  "id": 12,
  "user_id": 1,
  "account_id": 1,
  "transaction_type": "TOPUP",
  "amount": 1000000.00,
  "fee": 0.00,
  "total_amount": 1000000.00,
  "balance_before": 0.00,
  "balance_after": 1000000.00,
  "payment_method": "bank_transfer",
  "sender_bank": "BCA",
  "sender_account": "1234567890",
  "sender_name": "John Doe",
  "reference_id": "TOPUP-20260325-001",
  "status": "completed",
  "celery_task_id": "abc-123-xyz",
  "created_at": "2026-03-25T14:30:00Z"
}
```

*Transfer Out:*
```json
{
  "id": 14,
  "user_id": 1,
  "account_id": 1,
  "transaction_type": "TRANSFER_OUT",
  "amount": 250000.00,
  "fee": 0.00,
  "total_amount": 250000.00,
  "balance_before": 497500.00,
  "balance_after": 247500.00,
  "to_account_id": 2,
  "to_account_number": "9880987654321",
  "to_account_name": "Jane Doe",
  "reference_id": "TRF-20260325-001",
  "note": "Payment for lunch",
  "status": "completed",
  "celery_task_id": "ghi-789-xyz",
  "created_at": "2026-03-25T14:40:00Z"
}
```

---

## Relationships Diagram

```
┌─────────────┐       ┌─────────────┐       ┌─────────────────┐
│   USERS     │       │  ACCOUNTS   │       │  TRANSACTIONS   │
├─────────────┤       ├─────────────┤       ├─────────────────┤
│ id (PK)     │──1:1──│ user_id(FK) │──1:N──│ account_id (FK)│
│ email       │       │ account_    │       │ user_id (FK)   │
│ password_   │       │ number      │       │ type           │
│ hash        │       │ balance     │       │ amount         │
│ full_name   │       │ currency    │       │ status         │
│ phone_      │       │ account_    │       │ celery_task_id │
│ number      │       │ type        │       │ reference_id   │
│ role        │       │ is_active   │       │ created_at     │
│ is_active   │       └─────────────┘       └─────────────────┘
└─────────────┘
```

---

## Sequences (Oracle)

```sql
CREATE SEQUENCE users_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
CREATE SEQUENCE accounts_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
CREATE SEQUENCE transactions_seq START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE;
```

---

## Celery + Redis Structure

```
Redis: redis://localhost:6379/0
├── Broker for Celery tasks
└── Result backend: redis://localhost:6379/1

Queues:
  topup        -> tasks.process_topup
  withdraw     -> tasks.process_withdraw
  transfer     -> tasks.process_transfer
  notification -> tasks.send_notification
```

---

## Type Mapping (Go <-> Oracle)

| Go Type | Oracle Type |
|---------|-------------|
| uint | NUMBER(20) |
| float64 | NUMBER(18,2) |
| string | VARCHAR2(n) |
| time.Time | TIMESTAMP |
| bool | NUMBER(1) |

---

**Version:** 2.0.0
**Last Updated:** 2026-03-25
**Stack:** Go Echo + Oracle DB + Celery
