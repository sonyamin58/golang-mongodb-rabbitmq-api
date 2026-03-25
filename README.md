# Mini Bank API

REST API untuk aplikasi Mini Bank dengan fitur Register, Login, Topup, Withdraw, Transfer, Cek Saldo, dan Mutasi Rekening.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Go Echo v4 |
| Database | Oracle DB |
| ORM | GORM (Oracle driver) |
| Task Queue | Machinery (Go-native, Redis broker) |
| Authentication | JWT (HS256) |

## Features

- **Authentication**: Register, Login, JWT-based auth
- **Account**: Check balance, Topup, Withdraw, Transfer
- **Transactions**: History & detail dengan pagination
- **Async Processing**: Machinery workers (Go-native) untuk transaksi

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Oracle DB (via Docker or cloud)
- Redis

### 1. Clone & Setup

```bash
git clone https://github.com/sonyamin58/golang-mongodb-rabbitmq-api.git
cd golang-mongodb-rabbitmq-api
cp config.yaml config.local.yaml
```

### 2. Configure

Edit `config.local.yaml`:

```yaml
database:
  host: "localhost"
  port: 1521
  service_name: "XE"
  username: "system"
  password: "oracle"

redis:
  host: "localhost"
  port: 6379

machinery:
  broker: "redis://localhost:6379/0"
  result_backend: "redis://localhost:6379/1"
```

### 3. Run with Docker

```bash
# Start Oracle DB + Redis + API + Machinery Worker
docker compose up -d

# API will be available at http://localhost:8080
```

### 4. Run Locally (Development)

```bash
# Start Oracle & Redis
docker compose up -d oracle redis

# Install Go deps
make deps

# Run migrations (manual)
sqlplus system/oracle@localhost:1521/XE @migrations/001_init_schema.sql

# Start API
make run

# Start Machinery worker (new terminal)
make worker
```

## Project Structure

```
golang-mongodb-rabbitmq-api/
├── cmd/
│   ├── api/main.go          # API entry point
│   └── worker/main.go      # Machinery worker entry point
├── internal/
│   ├── config/              # Configuration loading
│   ├── handler/             # Echo HTTP handlers
│   ├── machinery/           # Machinery tasks & client
│   ├── middleware/          # JWT, rate limiting
│   ├── model/               # Database models
│   ├── repository/          # Data access layer
│   ├── service/             # Business logic
│   └── response/            # Standard response helpers
├── pkg/
│   ├── database/            # Oracle & Redis connections
│   └── validator/           # Input validation
├── migrations/
│   └── 001_init_schema.sql  # Oracle DDL scripts
├── api-design/
│   ├── API_DESIGN.md        # API specification
│   └── ERD_DB.MD           # Database schema docs
├── config.yaml
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | /api/v1/auth/register | Register | No |
| POST | /api/v1/auth/login | Login | No |
| GET | /api/v1/accounts/balance | Check balance | Yes |
| POST | /api/v1/accounts/topup | Topup | Yes |
| POST | /api/v1/accounts/withdraw | Withdraw | Yes |
| POST | /api/v1/accounts/transfer | Transfer | Yes |
| GET | /api/v1/transactions | Transaction history | Yes |
| GET | /api/v1/transactions/:id | Transaction detail | Yes |

## Makefile Commands

```bash
make deps          # Download dependencies
make build         # Build API binary
make run           # Run API server
make worker        # Run Machinery worker
make test          # Run tests
make lint          # Run linters
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
```

## Documentation

- [API Design](api-design/API_DESIGN.md)
- [Database Schema](api-design/ERD_DB.MD)
- [Architecture](ARCHITECTURE.md)

## Author

Sony Amin Gumelar
