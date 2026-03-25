# Mini Bank API

REST API untuk aplikasi Mini Bank dengan fitur Register, Login, Topup, Withdraw, Transfer, Cek Saldo, dan Mutasi Rekening.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Go Echo v4 |
| Database | Oracle DB |
| ORM | GORM (Oracle driver) |
| Message Broker | Celery + Redis |
| Authentication | JWT (HS256) |

## Features

- **Authentication**: Register, Login, JWT-based auth
- **Account**: Check balance, Topup, Withdraw, Transfer
- **Transactions**: History & detail dengan pagination
- **Async Processing**: Celery workers untuk transaksi

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

Edit `config.local.yaml` sesuai environment:

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
```

### 3. Run with Docker

```bash
# Start Oracle DB + Redis + API + Celery
docker-compose up -d

# Run migrations
docker-compose exec api /app/bin/api migrate
```

### 4. Run Locally (Development)

```bash
# Start Oracle & Redis
docker-compose up -d oracle redis

# Install Go deps
make deps

# Run migrations (manual)
sqlplus system/oracle@localhost:1521/XE @migrations/001_init_schema.sql

# Start API
make run

# Start Celery worker (new terminal)
make celery-worker

# Start Celery beat (new terminal)
make celery-beat
```

## Project Structure

```
golang-mongodb-rabbitmq-api/
├── cmd/
│   ├── api/main.go          # API entry point
│   └── worker/main.py      # Celery worker entry point
├── internal/
│   ├── config/              # Configuration loading
│   ├── handler/             # Echo HTTP handlers
│   ├── middleware/         # JWT, rate limiting
│   ├── model/               # Database models
│   ├── repository/          # Data access layer
│   ├── service/             # Business logic
│   └── response/           # Standard response helpers
├── pkg/
│   ├── database/            # Oracle & Redis connections
│   ├── celery/              # Celery task publisher
│   └── validator/           # Input validation
├── workers/
│   ├── celery_app.py        # Celery app config
│   └── tasks.py             # Async tasks
├── migrations/
│   └── 001_init_schema.sql  # Oracle DDL scripts
├── api-design/
│   ├── API_DESIGN.md        # API specification
│   └── ERD_DB.MD           # Database schema docs
├── config.yaml
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── requirements.txt
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

## Development

### Makefile Commands

```bash
make deps          # Download dependencies
make build         # Build binary
make run           # Run API
make test          # Run tests
make lint          # Run linters
make migrate       # Run migrations
make docker-build  # Build Docker image
make docker-up     # Start Docker services
make celery-worker # Run Celery worker
make celery-flower # Run Flower UI (port 5555)
```

### Testing

```bash
make test
make test-coverage
```

## Documentation

- [API Design](api-design/API_DESIGN.md)
- [Database Schema](api-design/ERD_DB.MD)
- [Architecture](ARCHITECTURE.md)

## Author

Sony Amin Gumelar
