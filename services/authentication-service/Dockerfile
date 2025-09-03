# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/main.go

# Final stage
FROM alpine:latest

COPY --from=builder /app/auth-service .
EXPOSE 8080

CMD ["./auth-service"]
