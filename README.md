# Portfolio

Personal portfolio and blog built with Go and Angular.

## Tech Stack

- **Frontend:** Angular 17, Tailwind CSS, shadcn/ui (Spartan), GSAP
- **Backend:** Go 1.22 (Chi Router, Clean Architecture)
- **Database:** PostgreSQL (sqlc)
- **Caching & Sessions:** Redis
- **Search:** Typesense
- **Messaging:** RabbitMQ & Kafka

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js (v18+)
- Go (1.22+)

### Setup

1. Copy the environment variables:
```bash
cp .env.example .env
```
Update `.env` with your desired configuration if needed.

2. Start the infrastructure services (Postgres, Redis, Typesense, RabbitMQ, Kafka):
```bash
docker compose up -d
```

### Running Locally

**Backend API:**
```bash
cd backend
go mod tidy

# Run with hot reload (requires air)
air
# Or run without hot reload
go run cmd/server/main.go
```
The API runs on `http://localhost:8080`.

**Frontend:**
```bash
cd frontend
npm install
npm run start
```
The web app runs on `http://localhost:4200`.

### Database Migrations
```bash
cd backend
migrate -path migrations -database "$DATABASE_URL" up
```
