# Build stage
FROM golang:1.25-alpine AS builder

# Copy clients
WORKDIR /app/clients/authentication-service
COPY ./clients/authentication-service/ ./

# Copy and build the service
WORKDIR /app/services/customer-service
COPY ./services/customer-service/go.mod ./services/customer-service/go.sum ./
RUN go mod download

COPY ./services/customer-service/ ./

RUN go build -o customer-service ./cmd/main.go

# Final stage
FROM alpine:latest

COPY --from=builder /app/services/customer-service/customer-service .
EXPOSE 8080

CMD ["./customer-service"]
