# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS base
RUN apk add --no-cache git curl
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# ─── Development stage (with air hot-reload) ──────────────────────────────────────
FROM base AS development
RUN go install github.com/air-verse/air@v1.61.0
COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# ─── Build stage ──────────────────────────────────────────────────────────────────
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/server ./cmd/server

# ─── Production stage (minimal image) ────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS production
COPY --from=builder /bin/server /server
COPY --from=builder /app/internal/infrastructure/persistence/migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/server"]
