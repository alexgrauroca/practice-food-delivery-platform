# Build stage
FROM golang:1.25-alpine AS builder

# Copy pkg
WORKDIR /app/pkg
COPY ./pkg/ ./

# Copy clients
WORKDIR /app/clients/authentication-service
COPY ./clients/authentication-service/ ./

# Copy and build the service
WORKDIR /app/services/restaurant-service
COPY ./services/restaurant-service/go.mod ./services/restaurant-service/go.sum ./
RUN go mod download

COPY ./services/restaurant-service/ ./

RUN go build -o restaurant-service ./cmd/main.go

# Final stage
FROM alpine:latest

COPY --from=builder /app/services/restaurant-service/restaurant-service .
EXPOSE 8080

CMD ["./restaurant-service"]
