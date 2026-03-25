# ARCHITECTURE.md — Mini Bank Service
> **Tujuan file ini**: Menjadi satu-satunya referensi bagi AI agent (Cursor, Copilot, Claude, dsb.) dalam memahami struktur, aturan arsitektur, dan konvensi kode pada project ini. **Baca seluruh file ini sebelum menulis atau memodifikasi kode apapun.**

---

## 1. Gambaran Umum

Project ini adalah **Mini Bank** — Go monorepo service dengan arsitektur **Clean Layered + Feature-Modular**.

### Prinsip Utama
- **Modular**: setiap domain bisnis dikapsulasi penuh dalam foldernya sendiri di `internal/modules/`
- **Clean Layered**: dependency hanya mengalir ke dalam — transport → usecase → domain ← infrastructure
- **Dependency Injection**: semua wiring dilakukan di `internal/di/`, tidak ada inisialisasi dependency di luar sana
- **Interface-first**: semua external client (DB, cache, messaging) dibungkus interface agar testable
- **Dua entrypoint**: `cmd/api` untuk HTTP server, `cmd/worker` untuk background job & consumer

### 5 Module Domain
| Module | Tanggung Jawab |
|--------|---------------|
| `auth` | Registrasi, login, logout, refresh token, verifikasi identitas |
| `account` | Buka rekening, cek saldo, daftar rekening milik user |
| `transaction` | Transfer, deposit, tarik tunai, riwayat transaksi |
| `notification` | Kirim & kelola notifikasi (push, email), tandai sudah dibaca |
| `settings` | Preferensi user, pengaturan PIN, limit transaksi, keamanan akun |

> **Penting untuk AI agent**: Jangan letakkan business logic di handler, repository, atau DTO. Jangan gunakan pola hexagonal. Komunikasi antar module hanya lewat event (messaging) atau shared interface di `pkg/` — **tidak boleh** import langsung antar module.

---

## 2. Folder Structure

```
/mini-bank
├── cmd/
│   ├── api/
│   │   └── main.go                    # Entry point HTTP API server
│   └── worker/
│       └── main.go                    # Entry point background worker & consumer
│
├── internal/
│   ├── modules/                       # Lapisan VERTICAL — satu folder per domain bisnis
│   │   │
│   │   ├── auth/                      # Module: Authentication & Authorization
│   │   │   ├── domain/
│   │   │   │   ├── entity.go          # User, Session, RefreshToken
│   │   │   │   ├── events.go          # UserRegistered, UserLoggedIn
│   │   │   │   └── repository.go      # UserRepository, SessionRepository (interface)
│   │   │   ├── usecase/
│   │   │   │   ├── register.go        # Daftar user baru
│   │   │   │   ├── login.go           # Login → JWT access + refresh token
│   │   │   │   ├── logout.go          # Invalidasi session & token
│   │   │   │   ├── refresh_token.go   # Rotate refresh token
│   │   │   │   └── verify_token.go    # Validasi JWT (dipakai middleware)
│   │   │   ├── transport/
│   │   │   │   └── http/
│   │   │   │       └── handler.go     # POST /auth/register, /auth/login, dsb.
│   │   │   └── dto/
│   │   │       ├── request.go         # RegisterRequest, LoginRequest
│   │   │       └── response.go        # AuthResponse (token pair)
│   │   │
│   │   ├── account/                   # Module: Rekening Bank
│   │   │   ├── domain/
│   │   │   │   ├── entity.go          # Account, AccountType, AccountStatus
│   │   │   │   ├── events.go          # AccountOpened, AccountSuspended
│   │   │   │   └── repository.go      # AccountRepository (interface)
│   │   │   ├── usecase/
│   │   │   │   ├── open_account.go    # Buka rekening baru
│   │   │   │   ├── get_account.go     # Detail satu rekening
│   │   │   │   ├── list_accounts.go   # Daftar rekening milik user
│   │   │   │   ├── get_balance.go     # Cek saldo terkini
│   │   │   │   └── suspend_account.go # Suspend/blokir rekening
│   │   │   ├── transport/
│   │   │   │   └── http/
│   │   │   │       └── handler.go     # GET/POST /accounts, GET /accounts/:id/balance
│   │   │   └── dto/
│   │   │       ├── request.go         # OpenAccountRequest
│   │   │       └── response.go        # AccountResponse, BalanceResponse
│   │   │
│   │   ├── transaction/               # Module: Transaksi Keuangan
│   │   │   ├── domain/
│   │   │   │   ├── entity.go          # Transaction, TransactionType, TransactionStatus
│   │   │   │   ├── events.go          # TransactionCreated, TransactionCompleted, TransactionFailed
│   │   │   │   └── repository.go      # TransactionRepository (interface)
│   │   │   ├── usecase/
│   │   │   │   ├── transfer.go        # Transfer antar rekening
│   │   │   │   ├── deposit.go         # Setor tunai / top-up
│   │   │   │   ├── withdraw.go        # Tarik tunai
│   │   │   │   ├── get_transaction.go # Detail satu transaksi
│   │   │   │   └── list_transactions.go # Riwayat transaksi dengan filter & pagination
│   │   │   ├── transport/
│   │   │   │   ├── http/
│   │   │   │   │   └── handler.go     # POST /transactions/transfer, /deposit, /withdraw
│   │   │   │   └── consumer/
│   │   │   │       └── event_consumer.go  # Consume event reversal / retry dari queue
│   │   │   └── dto/
│   │   │       ├── request.go         # TransferRequest, DepositRequest, WithdrawRequest
│   │   │       └── response.go        # TransactionResponse
│   │   │
│   │   ├── notification/              # Module: Notifikasi
│   │   │   ├── domain/
│   │   │   │   ├── entity.go          # Notification, NotificationType, NotificationChannel
│   │   │   │   ├── events.go          # NotificationSent, NotificationFailed
│   │   │   │   └── repository.go      # NotificationRepository (interface)
│   │   │   ├── usecase/
│   │   │   │   ├── send_notification.go   # Kirim notif (push / email / in-app)
│   │   │   │   ├── list_notifications.go  # Daftar notif milik user
│   │   │   │   └── mark_read.go           # Tandai notif sudah dibaca
│   │   │   ├── transport/
│   │   │   │   ├── http/
│   │   │   │   │   └── handler.go     # GET /notifications, PATCH /notifications/:id/read
│   │   │   │   └── consumer/
│   │   │   │       └── event_consumer.go  # Consume TransactionCompleted → kirim notif
│   │   │   └── dto/
│   │   │       ├── request.go
│   │   │       └── response.go
│   │   │
│   │   └── settings/                  # Module: Pengaturan & Preferensi User
│   │       ├── domain/
│   │       │   ├── entity.go          # UserSettings, TransactionLimit, SecurityConfig
│   │       │   ├── events.go          # PINChanged, LimitUpdated
│   │       │   └── repository.go      # SettingsRepository (interface)
│   │       ├── usecase/
│   │       │   ├── get_settings.go    # Ambil semua pengaturan user
│   │       │   ├── update_pin.go      # Ubah PIN transaksi
│   │       │   ├── update_limit.go    # Ubah limit transfer harian
│   │       │   ├── update_notification_pref.go  # Pilih channel notif (push/email)
│   │       │   └── toggle_feature.go  # Aktif/nonaktif fitur (biometrik, dsb.)
│   │       ├── transport/
│   │       │   └── http/
│   │       │       └── handler.go     # GET/PATCH /settings, POST /settings/pin
│   │       └── dto/
│   │           ├── request.go
│   │           └── response.go
│   │
│   ├── infrastructure/                # Lapisan HORIZONTAL — teknologi eksternal (shared)
│   │   ├── db/
│   │   │   └── postgres_client.go     # Koneksi PostgreSQL (utama untuk banking)
│   │   ├── cache/
│   │   │   └── redis_client.go        # Session, token blacklist, rate limiter
│   │   ├── messaging/
│   │   │   └── rabbitmq_client.go     # Event bus antar module
│   │   ├── email/
│   │   │   └── email_client.go        # SMTP / provider email
│   │   ├── fcm/
│   │   │   └── fcm_client.go          # Firebase push notification
│   │   └── repository/                # Implementasi konkret semua domain repository
│   │       ├── user_repo.go           # Implements auth/domain.UserRepository
│   │       ├── session_repo.go        # Implements auth/domain.SessionRepository
│   │       ├── account_repo.go        # Implements account/domain.AccountRepository
│   │       ├── transaction_repo.go    # Implements transaction/domain.TransactionRepository
│   │       ├── notification_repo.go   # Implements notification/domain.NotificationRepository
│   │       └── settings_repo.go       # Implements settings/domain.SettingsRepository
│   │
│   ├── config/
│   │   └── config.go                  # Struct Config lengkap, load dari env/yaml
│   │
│   ├── di/                            # Dependency Injection — semua wiring di sini
│   │   ├── container.go               # AppContainer struct — satu tempat semua dependency
│   │   ├── provider_core.go           # Init DB, Redis, RabbitMQ, Logger
│   │   ├── provider_auth.go           # Init semua layer module auth
│   │   ├── provider_account.go        # Init semua layer module account
│   │   ├── provider_transaction.go    # Init semua layer module transaction
│   │   ├── provider_notification.go   # Init semua layer module notification
│   │   ├── provider_settings.go       # Init semua layer module settings
│   │   ├── provider_api.go            # Init HTTP router + semua handler
│   │   └── provider_worker.go         # Init consumer / background worker
│   │
│   ├── middleware/
│   │   ├── auth.go                    # JWT validation — gunakan auth/usecase.VerifyToken
│   │   ├── pin.go                     # PIN verification untuk transaksi sensitif
│   │   ├── rate_limiter.go            # Rate limit per user/IP via Redis
│   │   ├── logger.go                  # Request/response logging
│   │   └── recovery.go                # Panic recovery
│   │
│   └── api/
│       └── router.go                  # Global route aggregator — daftarkan semua handler
│
├── pkg/                               # Shared utilities — boleh dipakai semua layer
│   ├── jwt/
│   │   └── jwt.go                     # Generate & parse JWT
│   ├── response/
│   │   └── response.go                # Helper format HTTP response standar
│   ├── crypto/
│   │   └── hash.go                    # Bcrypt hash (password, PIN)
│   ├── validator/
│   │   └── validator.go               # Input validation helper
│   ├── pagination/
│   │   └── pagination.go              # Struct & helper pagination
│   └── util/
│       └── util.go                    # Fungsi umum (generate ID, format uang, dsb.)
│
├── static/
│   ├── email_templates/
│   │   ├── welcome.html               # Email selamat datang saat register
│   │   ├── transaction.html           # Notif transaksi berhasil
│   │   └── otp.html                   # Email OTP verifikasi
│   └── firebase/
│       └── fcm-key.json
│
├── build/
│   ├── Dockerfile.api
│   ├── Dockerfile.worker
│   └── docker-compose.yaml
│
├── configs/
│   ├── config.yaml
│   └── config.local.yaml
│
├── migrations/                        # SQL migration files (urut)
│   ├── 001_create_users.sql
│   ├── 002_create_sessions.sql
│   ├── 003_create_accounts.sql
│   ├── 004_create_transactions.sql
│   ├── 005_create_notifications.sql
│   └── 006_create_settings.sql
│
├── test/
│   ├── integration/
│   └── e2e/
│
├── go.mod
├── go.sum
├── ARCHITECTURE.md                    # ← file ini
└── README.md
```

---

## 3. Domain Model per Module

### Module: `auth`

```go
// internal/modules/auth/domain/entity.go
package domain

type User struct {
    ID           string
    FullName     string
    Email        string
    Phone        string
    PasswordHash string
    PINHash      string
    IsVerified   bool
    IsActive     bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Session struct {
    ID           string
    UserID       string
    RefreshToken string
    UserAgent    string
    IP           string
    ExpiresAt    time.Time
    CreatedAt    time.Time
}

// Method bisnis
func (u *User) IsEligibleToLogin() bool {
    return u.IsActive && u.IsVerified
}
```

```go
// internal/modules/auth/domain/repository.go
package domain

type UserRepository interface {
    Save(ctx context.Context, u *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, u *User) error
}

type SessionRepository interface {
    Save(ctx context.Context, s *Session) error
    FindByRefreshToken(ctx context.Context, token string) (*Session, error)
    DeleteByUserID(ctx context.Context, userID string) error
    DeleteByID(ctx context.Context, id string) error
}
```

---

### Module: `account`

```go
// internal/modules/account/domain/entity.go
package domain

type AccountType  string
type AccountStatus string

const (
    AccountTypeSavings AccountType = "savings"
    AccountTypeGiro    AccountType = "giro"

    AccountStatusActive    AccountStatus = "active"
    AccountStatusSuspended AccountStatus = "suspended"
    AccountStatusClosed    AccountStatus = "closed"
)

type Account struct {
    ID            string
    UserID        string
    AccountNumber string        // 16-digit unik
    Type          AccountType
    Status        AccountStatus
    Balance       int64         // dalam satuan sen (hindari float)
    Currency      string        // "IDR"
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// Method bisnis
func (a *Account) IsOperational() bool {
    return a.Status == AccountStatusActive
}

func (a *Account) Debit(amount int64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    if a.Balance < amount {
        return errors.New("insufficient balance")
    }
    a.Balance -= amount
    return nil
}

func (a *Account) Credit(amount int64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    a.Balance += amount
    return nil
}
```

```go
// internal/modules/account/domain/repository.go
package domain

type AccountRepository interface {
    Save(ctx context.Context, a *Account) error
    FindByID(ctx context.Context, id string) (*Account, error)
    FindByAccountNumber(ctx context.Context, number string) (*Account, error)
    FindByUserID(ctx context.Context, userID string) ([]*Account, error)
    Update(ctx context.Context, a *Account) error
    UpdateBalance(ctx context.Context, id string, newBalance int64) error
}
```

---

### Module: `transaction`

```go
// internal/modules/transaction/domain/entity.go
package domain

type TransactionType   string
type TransactionStatus string

const (
    TxTypeTransfer  TransactionType = "transfer"
    TxTypeDeposit   TransactionType = "deposit"
    TxTypeWithdraw  TransactionType = "withdraw"

    TxStatusPending   TransactionStatus = "pending"
    TxStatusCompleted TransactionStatus = "completed"
    TxStatusFailed    TransactionStatus = "failed"
    TxStatusReversed  TransactionStatus = "reversed"
)

type Transaction struct {
    ID                string
    ReferenceNumber   string            // kode unik untuk tracing
    Type              TransactionType
    Status            TransactionStatus
    FromAccountID     string            // kosong jika deposit
    ToAccountID       string            // kosong jika withdraw
    Amount            int64             // dalam satuan sen
    Fee               int64
    Note              string
    FailureReason     string
    CreatedAt         time.Time
    CompletedAt       *time.Time
}

// Method bisnis
func (t *Transaction) Complete() {
    now := time.Now()
    t.Status = TxStatusCompleted
    t.CompletedAt = &now
}

func (t *Transaction) Fail(reason string) {
    t.Status = TxStatusFailed
    t.FailureReason = reason
}
```

```go
// internal/modules/transaction/domain/repository.go
package domain

type TransactionRepository interface {
    Save(ctx context.Context, t *Transaction) error
    FindByID(ctx context.Context, id string) (*Transaction, error)
    FindByReferenceNumber(ctx context.Context, ref string) (*Transaction, error)
    FindByAccountID(ctx context.Context, accountID string, filter Filter) ([]*Transaction, int64, error)
    Update(ctx context.Context, t *Transaction) error
}

type Filter struct {
    Type      *TransactionType
    Status    *TransactionStatus
    StartDate *time.Time
    EndDate   *time.Time
    Page      int
    Limit     int
}
```

---

### Module: `notification`

```go
// internal/modules/notification/domain/entity.go
package domain

type NotificationType    string
type NotificationChannel string

const (
    NotifTypeTransaction  NotificationType = "transaction"
    NotifTypeSecurity     NotificationType = "security"
    NotifTypePromo        NotificationType = "promo"
    NotifTypeSystem       NotificationType = "system"

    ChannelPush  NotificationChannel = "push"
    ChannelEmail NotificationChannel = "email"
    ChannelInApp NotificationChannel = "in_app"
)

type Notification struct {
    ID        string
    UserID    string
    Type      NotificationType
    Channel   NotificationChannel
    Title     string
    Body      string
    Metadata  map[string]string   // data tambahan (mis: transaction_id)
    IsRead    bool
    SentAt    *time.Time
    CreatedAt time.Time
}

func (n *Notification) MarkAsRead() {
    n.IsRead = true
}
```

---

### Module: `settings`

```go
// internal/modules/settings/domain/entity.go
package domain

type UserSettings struct {
    ID                     string
    UserID                 string
    DailyTransferLimit     int64          // dalam satuan sen
    NotificationChannels   []string       // ["push", "email"]
    IsBiometricEnabled     bool
    IsTransactionPINActive bool
    Language               string         // "id", "en"
    UpdatedAt              time.Time
}

// Method bisnis
func (s *UserSettings) IsChannelEnabled(channel string) bool {
    for _, c := range s.NotificationChannels {
        if c == channel {
            return true
        }
    }
    return false
}

func (s *UserSettings) IsWithinDailyLimit(amount int64) bool {
    return amount <= s.DailyTransferLimit
}
```

---

## 4. Layer Architecture

Setiap module mengikuti layer ini. Dependency hanya boleh mengarah **ke dalam**:

```
┌─────────────────────────────────────┐
│     Transport (HTTP / Consumer)     │  ← parsing I/O saja
└──────────────────┬──────────────────┘
                   ↓
┌─────────────────────────────────────┐
│             Use Case                │  ← logika aplikasi & orkestrasi
└──────────────────┬──────────────────┘
                   ↓
┌─────────────────────────────────────┐
│              Domain                 │  ← entity + interface (ZERO dependency)
└──────────────────▲──────────────────┘
                   │ implements
┌─────────────────────────────────────┐
│          Infrastructure             │  ← DB, cache, messaging, email, FCM
└─────────────────────────────────────┘
```

### Layer 1 — Domain
**Lokasi**: `internal/modules/{module}/domain/`

- `entity.go` → struct bisnis + method bisnis yang melekat pada entity
- `repository.go` → **interface** kontrak persistence (PORT) — bukan implementasi
- `events.go` → domain events untuk komunikasi async antar module
- **WAJIB zero external dependency** — tidak boleh import package manapun selain standard library

### Layer 2 — Use Case
**Lokasi**: `internal/modules/{module}/usecase/`

- Satu file per use case (SRP): `transfer.go`, `login.go`, dsb.
- Orkestrasi domain + repository + external service
- Menerima semua dependency via **constructor injection** dengan tipe **interface**
- **Tidak boleh** akses DB/Redis/MQ langsung
- **Tidak boleh** tahu soal HTTP (fiber.Ctx, http.Request, dsb.)
- Boleh publish domain event ke messaging untuk komunikasi ke module lain

```go
// Contoh: internal/modules/transaction/usecase/transfer.go
package usecase

// TransferUseCase menangani alur transfer antar rekening.
type TransferUseCase struct {
    txRepo      domain.TransactionRepository      // interface
    accountRepo accountdomain.AccountRepository   // interface dari pkg/contract
    settingsRepo settingsdomain.SettingsRepository // interface dari pkg/contract
    publisher   messaging.Publisher               // interface
    logger      logger.Logger
}

func NewTransferUseCase(
    txRepo domain.TransactionRepository,
    accountRepo accountdomain.AccountRepository,
    settingsRepo settingsdomain.SettingsRepository,
    publisher messaging.Publisher,
    logger logger.Logger,
) *TransferUseCase {
    return &TransferUseCase{
        txRepo: txRepo, accountRepo: accountRepo,
        settingsRepo: settingsRepo, publisher: publisher, logger: logger,
    }
}

func (uc *TransferUseCase) Execute(ctx context.Context, req dto.TransferRequest) (*dto.TransactionResponse, error) {
    // 1. Validasi limit harian dari settings
    // 2. Cek rekening asal & tujuan
    // 3. Debit rekening asal (domain method)
    // 4. Credit rekening tujuan (domain method)
    // 5. Simpan transaksi
    // 6. Publish event TransactionCompleted → notification module consume
    // 7. Return response
}
```

### Layer 3 — Transport
**Lokasi**: `internal/modules/{module}/transport/`

**HTTP Handler** — hanya: parse request → validate → call usecase → format response

```go
// Contoh: internal/modules/transaction/transport/http/handler.go
package http

type TransactionHandler struct {
    transferUC  *usecase.TransferUseCase
    depositUC   *usecase.DepositUseCase
    withdrawUC  *usecase.WithdrawUseCase
    listUC      *usecase.ListTransactionsUseCase
}

func (h *TransactionHandler) Transfer(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)  // dari JWT middleware
    var req dto.TransferRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, err)
    }
    req.FromUserID = userID
    result, err := h.transferUC.Execute(c.Context(), req)
    if err != nil {
        return response.Error(c, err)
    }
    return response.OK(c, result)
}

func (h *TransactionHandler) RegisterRoutes(r fiber.Router, authMW, pinMW fiber.Handler) {
    tx := r.Group("/transactions", authMW)
    tx.Post("/transfer",  pinMW, h.Transfer)   // butuh PIN
    tx.Post("/deposit",   h.Deposit)
    tx.Post("/withdraw",  pinMW, h.Withdraw)   // butuh PIN
    tx.Get("/",           h.List)
    tx.Get("/:id",        h.GetDetail)
}
```

**Consumer** — hanya: consume message → unmarshal → call usecase

```go
// Contoh: internal/modules/notification/transport/consumer/event_consumer.go
package consumer

// NotificationConsumer mendengarkan event dari module lain dan men-trigger notifikasi.
type NotificationConsumer struct {
    sendUC *usecase.SendNotificationUseCase
    mq     messaging.Consumer
    logger logger.Logger
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
    // Consume event TransactionCompleted yang di-publish oleh transaction module
    return c.mq.Consume(ctx, "transaction.completed", c.handleTransactionCompleted)
}

func (c *NotificationConsumer) handleTransactionCompleted(ctx context.Context, msg []byte) error {
    var event events.TransactionCompleted
    if err := json.Unmarshal(msg, &event); err != nil {
        return err
    }
    return c.sendUC.Execute(ctx, dto.SendNotificationRequest{
        UserID: event.UserID,
        Type:   "transaction",
        Title:  "Transfer Berhasil",
        Body:   fmt.Sprintf("Transfer Rp%s telah berhasil.", formatMoney(event.Amount)),
    })
}
```

### Layer 4 — Infrastructure (Implementasi)
**Lokasi**: `internal/infrastructure/repository/`

- Implementasi konkret semua interface di `domain/repository.go`
- Boleh import `database/sql`, `go.mongodb.org`, `redis`, dsb.
- Tidak boleh diimport langsung oleh domain atau usecase

```go
// Contoh: internal/infrastructure/repository/account_repo.go
package repository

// PostgresAccountRepo mengimplementasi account/domain.AccountRepository menggunakan PostgreSQL.
type PostgresAccountRepo struct {
    db *sql.DB
}

func NewPostgresAccountRepo(db *sql.DB) *PostgresAccountRepo {
    return &PostgresAccountRepo{db: db}
}

func (r *PostgresAccountRepo) UpdateBalance(ctx context.Context, id string, newBalance int64) error {
    _, err := r.db.ExecContext(ctx,
        `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`,
        newBalance, id,
    )
    return err
}
```

---

## 5. Dependency Injection

**Lokasi**: `internal/di/`

Semua `New*()` dipanggil dari sini. Tidak boleh ada inisialisasi dependency di tempat lain.

```go
// internal/di/container.go
package di

type AppContainer struct {
    // Core
    DB       *sql.DB
    Redis    *redis.Client
    RabbitMQ *amqp.Connection
    Logger   logger.Logger

    // Auth module
    UserRepo         authdomain.UserRepository
    SessionRepo      authdomain.SessionRepository
    RegisterUC       *authusecase.RegisterUseCase
    LoginUC          *authusecase.LoginUseCase
    LogoutUC         *authusecase.LogoutUseCase
    RefreshTokenUC   *authusecase.RefreshTokenUseCase
    VerifyTokenUC    *authusecase.VerifyTokenUseCase
    AuthHandler      *authhttp.AuthHandler

    // Account module
    AccountRepo      accountdomain.AccountRepository
    OpenAccountUC    *accountusecase.OpenAccountUseCase
    GetAccountUC     *accountusecase.GetAccountUseCase
    ListAccountsUC   *accountusecase.ListAccountsUseCase
    GetBalanceUC     *accountusecase.GetBalanceUseCase
    AccountHandler   *accounthttp.AccountHandler

    // Transaction module
    TransactionRepo     txdomain.TransactionRepository
    TransferUC          *txusecase.TransferUseCase
    DepositUC           *txusecase.DepositUseCase
    WithdrawUC          *txusecase.WithdrawUseCase
    ListTransactionsUC  *txusecase.ListTransactionsUseCase
    TransactionHandler  *txhttp.TransactionHandler
    TransactionConsumer *txconsumer.TransactionConsumer

    // Notification module
    NotificationRepo      notifdomain.NotificationRepository
    SendNotificationUC    *notifusecase.SendNotificationUseCase
    ListNotificationsUC   *notifusecase.ListNotificationsUseCase
    MarkReadUC            *notifusecase.MarkReadUseCase
    NotificationHandler   *notifhttp.NotificationHandler
    NotificationConsumer  *notifconsumer.NotificationConsumer

    // Settings module
    SettingsRepo          settingsdomain.SettingsRepository
    GetSettingsUC         *settingsusecase.GetSettingsUseCase
    UpdatePINUC           *settingsusecase.UpdatePINUseCase
    UpdateLimitUC         *settingsusecase.UpdateLimitUseCase
    UpdateNotifPrefUC     *settingsusecase.UpdateNotificationPrefUseCase
    SettingsHandler       *settingshttp.SettingsHandler
}
```

```go
// internal/di/provider_transaction.go — contoh wiring satu module
package di

func ProvideTransaction(c *AppContainer) {
    // Repository
    c.TransactionRepo = repository.NewPostgresTransactionRepo(c.DB)

    // Use cases — inject interface, bukan concrete
    c.TransferUC = txusecase.NewTransferUseCase(
        c.TransactionRepo,
        c.AccountRepo,      // dari provider_account.go
        c.SettingsRepo,     // dari provider_settings.go
        messaging.NewRabbitMQPublisher(c.RabbitMQ),
        c.Logger,
    )
    c.DepositUC  = txusecase.NewDepositUseCase(c.TransactionRepo, c.AccountRepo, c.Logger)
    c.WithdrawUC = txusecase.NewWithdrawUseCase(c.TransactionRepo, c.AccountRepo, c.SettingsRepo, c.Logger)
    c.ListTransactionsUC = txusecase.NewListTransactionsUseCase(c.TransactionRepo, c.Logger)

    // Handler
    c.TransactionHandler = txhttp.NewTransactionHandler(
        c.TransferUC, c.DepositUC, c.WithdrawUC, c.ListTransactionsUC,
    )

    // Consumer
    c.TransactionConsumer = txconsumer.NewTransactionConsumer(
        messaging.NewRabbitMQConsumer(c.RabbitMQ),
        c.Logger,
    )
}
```

---

## 6. Komunikasi Antar Module

Module **tidak boleh** saling import secara langsung. Komunikasi dilakukan melalui dua cara:

### a) Domain Event via Message Queue
Digunakan untuk side-effect asinkron (fire and forget):

```
transaction module  →  publish "transaction.completed"
                              ↓
notification module  →  consume → kirim notif ke user
```

```go
// Event yang di-publish (definisikan di pkg/events/ agar bisa dipakai kedua module)
// pkg/events/transaction_events.go
package events

const TopicTransactionCompleted = "transaction.completed"

type TransactionCompleted struct {
    TransactionID string    `json:"transaction_id"`
    UserID        string    `json:"user_id"`
    Amount        int64     `json:"amount"`
    Type          string    `json:"type"`
    CompletedAt   time.Time `json:"completed_at"`
}
```

### b) Shared Contract Interface via `pkg/contract/`
Digunakan saat satu module perlu data dari module lain secara **sinkron** (misal: transaction butuh cek limit dari settings):

```go
// pkg/contract/account_contract.go
// Definisikan interface minimal yang dibutuhkan module lain
package contract

type AccountReader interface {
    FindByID(ctx context.Context, id string) (*AccountSummary, error)
    FindByAccountNumber(ctx context.Context, number string) (*AccountSummary, error)
}

type AccountSummary struct {
    ID            string
    AccountNumber string
    Balance       int64
    Status        string
    UserID        string
}
```

> **Aturan**: Jangan import `internal/modules/account/domain` dari dalam `internal/modules/transaction/`. Selalu gunakan interface di `pkg/contract/` sebagai jembatan.

---

## 7. Aturan Wajib (Rules for AI Agent)

### Yang HARUS dilakukan

| Rule | Keterangan |
|------|-----------|
| ✅ Handler hanya call usecase | Parse request → validate → call UC → format response |
| ✅ Usecase hanya orkestrasi | Tidak akses DB/Redis/MQ langsung |
| ✅ Domain tetap pure | Tidak import apapun selain standard library |
| ✅ Repository via interface | Implementasi di `infrastructure/repository/` |
| ✅ Constructor injection | Semua dependency masuk via parameter constructor |
| ✅ DI hanya di `internal/di/` | Semua `New*()` dipanggil dari sana |
| ✅ Gunakan interface untuk external client | Repo, FCM, email, MQ — bukan concrete struct |
| ✅ DTO hanya untuk I/O | Tidak pakai DTO sebagai domain entity |
| ✅ Satu file per use case | SRP — fokus satu alur bisnis per file |
| ✅ Balance dalam satuan sen (`int64`) | Hindari `float64` untuk uang |
| ✅ Transaksi menggunakan DB transaction | Transfer harus atomic (debit + credit dalam satu tx) |
| ✅ Event antar module via `pkg/events/` | Definisi event di pkg, bukan di dalam module |
| ✅ Logger dari infrastructure | Selalu inject, tidak pakai `fmt.Println` |
| ✅ Komentar pada struct dan method publik | Wajib untuk semua exported symbol |

### Yang DILARANG

| Larangan | Alasan |
|----------|--------|
| ❌ Business logic di handler | Melanggar SRP |
| ❌ DB query di usecase | Usecase tidak boleh tahu detail persistence |
| ❌ Import infra di domain | Domain harus zero-dependency |
| ❌ Import langsung antar module | Gunakan `pkg/contract/` atau event |
| ❌ `float64` untuk nilai uang | Gunakan `int64` (satuan sen) |
| ❌ Pola hexagonal / ports-adapters | Proyek ini clean layered |
| ❌ `new()` / inisialisasi di luar DI | Semua wiring di `internal/di/` |
| ❌ Global variable untuk dependency | Constructor injection |
| ❌ Hardcode konfigurasi | Ambil dari `config.Config` |
| ❌ Transfer tanpa DB transaction | Harus atomic, gunakan `sql.Tx` |
| ❌ Simpan plain text password/PIN | Selalu hash dengan `pkg/crypto` |

---

## 8. Konvensi Penamaan

### File
```
entity.go                     → domain entity + method bisnis
repository.go                 → repository interface (domain layer)
events.go                     → domain events
{action}.go                   → usecase: transfer.go, login.go, open_account.go
handler.go                    → HTTP handler per module
event_consumer.go             → async consumer per module
request.go                    → DTO input
response.go                   → DTO output
{entity}_repo.go              → implementasi repository: account_repo.go
provider_{module}.go          → DI provider: provider_transaction.go
```

### Package
```go
package domain      // internal/modules/*/domain/
package usecase     // internal/modules/*/usecase/
package http        // internal/modules/*/transport/http/
package consumer    // internal/modules/*/transport/consumer/
package dto         // internal/modules/*/dto/
package repository  // internal/infrastructure/repository/
package di          // internal/di/
package contract    // pkg/contract/
package events      // pkg/events/
```

### Struct & Interface
```go
// Interface — nama tanpa prefix "I", tanpa suffix "Interface"
type UserRepository interface { ... }
type AccountRepository interface { ... }

// Implementasi konkret — prefix teknologi
type PostgresUserRepo struct { ... }
type PostgresAccountRepo struct { ... }
type RedisSessionRepo struct { ... }

// Use case — suffix "UseCase"
type TransferUseCase struct { ... }
type LoginUseCase struct { ... }
type OpenAccountUseCase struct { ... }

// Handler — suffix "Handler"
type TransactionHandler struct { ... }
type AuthHandler struct { ... }

// Consumer — suffix "Consumer"
type NotificationConsumer struct { ... }
type TransactionConsumer struct { ... }
```

---

## 9. Endpoint API per Module

### Auth — `/api/v1/auth`
```
POST   /auth/register          Daftar akun baru
POST   /auth/login             Login, return access + refresh token
POST   /auth/logout            Logout, invalidasi session     [JWT required]
POST   /auth/refresh           Rotate refresh token
```

### Account — `/api/v1/accounts`
```
POST   /accounts               Buka rekening baru             [JWT required]
GET    /accounts               Daftar rekening milik user     [JWT required]
GET    /accounts/:id           Detail rekening                [JWT required]
GET    /accounts/:id/balance   Cek saldo                      [JWT required]
```

### Transaction — `/api/v1/transactions`
```
POST   /transactions/transfer  Transfer antar rekening        [JWT + PIN required]
POST   /transactions/deposit   Setor / top-up                 [JWT required]
POST   /transactions/withdraw  Tarik tunai                    [JWT + PIN required]
GET    /transactions           Riwayat transaksi (filter+paging) [JWT required]
GET    /transactions/:id       Detail transaksi               [JWT required]
```

### Notification — `/api/v1/notifications`
```
GET    /notifications          Daftar notif user              [JWT required]
PATCH  /notifications/:id/read Tandai sudah dibaca            [JWT required]
PATCH  /notifications/read-all Tandai semua sudah dibaca      [JWT required]
```

### Settings — `/api/v1/settings`
```
GET    /settings               Ambil semua pengaturan         [JWT required]
PATCH  /settings/pin           Ubah PIN transaksi             [JWT required]
PATCH  /settings/limit         Ubah limit transfer harian     [JWT required]
PATCH  /settings/notification  Atur channel notifikasi        [JWT required]
PATCH  /settings/features      Toggle fitur (biometrik, dsb.) [JWT required]
```

---

## 10. Database Schema (Ringkasan)

Gunakan **PostgreSQL**. Semua ID menggunakan `UUID`. Semua nilai uang dalam `BIGINT` (satuan sen).

```sql
-- Auth
users        (id, full_name, email, phone, password_hash, pin_hash, is_verified, is_active, created_at, updated_at)
sessions     (id, user_id, refresh_token, user_agent, ip, expires_at, created_at)

-- Account
accounts     (id, user_id, account_number, type, status, balance, currency, created_at, updated_at)

-- Transaction
transactions (id, reference_number, type, status, from_account_id, to_account_id, amount, fee, note, failure_reason, created_at, completed_at)

-- Notification
notifications (id, user_id, type, channel, title, body, metadata, is_read, sent_at, created_at)

-- Settings
user_settings (id, user_id, daily_transfer_limit, notification_channels, is_biometric_enabled, is_transaction_pin_active, language, updated_at)
```

---

## 11. Entry Points

### `cmd/api/main.go`
```
main()
  → load config
  → build AppContainer (di.New())
  → register routes (router.Setup(container))
  → start HTTP server
  → handle graceful shutdown
```

### `cmd/worker/main.go`
```
main()
  → load config
  → build AppContainer (di.New())
  → start NotificationConsumer.Start()
  → start TransactionConsumer.Start()  (jika ada retry/reversal)
  → block until shutdown signal
```

---

## 12. Testing Convention

```
internal/modules/{module}/usecase/{usecase}_test.go  → unit test (mock semua interface)
internal/modules/{module}/domain/{entity}_test.go    → unit test method bisnis domain
test/integration/                                     → integration test (hit DB asli)
test/e2e/                                             → end-to-end test full flow
```

Contoh unit test transfer:
```go
func TestTransfer_Success(t *testing.T) {
    mockTxRepo      := new(mocks.TransactionRepository)
    mockAccountRepo := new(mocks.AccountRepository)
    mockSettingsRepo := new(mocks.SettingsRepository)
    mockPublisher   := new(mocks.Publisher)

    uc := usecase.NewTransferUseCase(mockTxRepo, mockAccountRepo, mockSettingsRepo, mockPublisher, logger.Noop())

    mockSettingsRepo.On("FindByUserID", mock.Anything, "user-1").
        Return(&settingsdomain.UserSettings{DailyTransferLimit: 100_000_000}, nil)
    mockAccountRepo.On("FindByID", mock.Anything, "acc-from").
        Return(&contract.AccountSummary{Balance: 50_000_000, Status: "active"}, nil)
    // ... setup mocks lainnya

    _, err := uc.Execute(context.Background(), dto.TransferRequest{
        FromAccountID: "acc-from",
        ToAccountID:   "acc-to",
        Amount:        10_000_000,
    })

    assert.NoError(t, err)
    mockTxRepo.AssertExpectations(t)
}

func TestTransfer_InsufficientBalance(t *testing.T) { ... }
func TestTransfer_ExceedDailyLimit(t *testing.T) { ... }
func TestTransfer_AccountSuspended(t *testing.T) { ... }
```

---

## 13. Checklist Sebelum Generate Kode

- [ ] File ini masuk layer mana? (domain / usecase / transport / infrastructure / di)
- [ ] Ada business logic yang tidak sengaja masuk ke handler atau repository?
- [ ] Dependency di-inject via constructor dengan tipe interface?
- [ ] Domain layer bebas dari import infrastructure dan transport?
- [ ] DI provider di-update untuk komponen baru?
- [ ] Komunikasi ke module lain sudah lewat event atau `pkg/contract/`?
- [ ] Nilai uang menggunakan `int64` (bukan `float64`)?
- [ ] Operasi yang melibatkan debit+credit menggunakan DB transaction?
- [ ] Password dan PIN di-hash sebelum disimpan?
- [ ] Penamaan file dan package sesuai konvensi di Section 8?
- [ ] Ada unit test minimal untuk use case yang dibuat?
- [ ] Semua exported struct dan method sudah ada komentar?

---

*Update file ini setiap kali ada perubahan arsitektur signifikan. Versi terakhir selalu menjadi acuan utama AI agent.*