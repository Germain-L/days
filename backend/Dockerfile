FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go

FROM alpine:latest

# Install PostgreSQL client for migrations
RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/db/migrations ./migrations
COPY --from=builder /app/scripts/migrate.sh ./migrate.sh

# Make the migration script executable
RUN chmod +x ./migrate.sh

EXPOSE 8080

CMD ["./main"]