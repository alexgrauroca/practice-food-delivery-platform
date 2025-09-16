# Build stage
FROM golang:1.25-alpine AS builder

# Copy pkg
WORKDIR /app/pkg
COPY ./pkg/ ./

# Copy clients
WORKDIR /app/clients/authentication-service
COPY ./clients/authentication-service/ ./

# Copy and build the service
WORKDIR /app/services/authentication-service
COPY ./services/authentication-service/go.mod ./services/authentication-service/go.sum ./
RUN go mod download

COPY ./services/authentication-service/ ./

RUN go build -o authentication-service ./cmd/main.go

# Final stage
FROM alpine:latest

COPY --from=builder /app/services/authentication-service/authentication-service .
EXPOSE 8080

CMD ["./authentication-service"]
