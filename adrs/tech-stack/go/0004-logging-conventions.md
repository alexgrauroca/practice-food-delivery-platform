# 4. Logging Conventions

## Status

Accepted

Date: 2025-06-25

## Context

As we build domain services in Go, we need consistent logging practices that:

- Provide structured, machine-parseable logs
- Include appropriate context for troubleshooting
- Support different log levels for various environments
- Correlate logs across service boundaries
- Balance verbosity with performance considerations
- Follow industry best practices for observability

Inconsistent logging approaches lead to several challenges:

- Difficulty tracing requests across service boundaries
- Inconsistent log levels and formats
- Missing critical context in production issues
- Performance overhead from excessive logging
- Difficulty filtering logs for specific events or errors

## Decision

We will adopt zap (go.uber.org/zap) as our standard logging library with the following conventions:

### 1. Structured Logging

- Use structured, JSON-formatted logs in production
- Include contextual fields with all log messages
- Use consistent field naming across services

```go
logger.Info("Customer created successfully", 
    zap.String("customer_id", customer.ID),
    zap.String("email", customer.Email))
```

### 2. Log Levels

- **Debug**: Detailed information for debugging (disabled in production)
- **Info**: Normal application behavior, key lifecycle events
- **Warn**: Potential issues that don't prevent operation
- **Error**: Errors that affect a specific operation but allow the service to continue
- **Fatal**: Critical errors that prevent the service from continuing

### 3. Request Context Propagation

- Include request ID in all logs related to a request
- Pass logger with request context through the call chain
- Use our logctx package to extract request information

```go
logctx.LoggerWithRequestInfo(ctx, r.logger).Error("Failed to find customer", zap.Error(err))
```

### 4. Standard Log Fields

| Field Name | Description | Example |
|------------|-------------|--------|
| request_id | Unique request identifier | "req_abc123" |
| timestamp | ISO 8601 timestamp | "2025-06-25T14:35:05Z" |
| level | Log level | "info", "error" |
| message | Human-readable message | "Customer created successfully" |
| service | Service name | "authentication-service" |
| method | HTTP method | "POST" |
| path | Request path | "/v1.0/customers/register" |
| duration_ms | Operation duration in ms | 127 |
| user_id | User identifier (when available) | "usr_123" |
| error | Error details | {"message":"document not found"} |

### 5. Error Logging

- Always include error objects with zap.Error()
- Add context about the operation that failed
- Include relevant IDs and parameters (sanitized of sensitive data)

```go
logger.Error("Failed to create customer",
    zap.Error(err),
    zap.String("email", email),
    zap.String("operation", "CreateCustomer"))
```

### 6. Sensitive Data Handling

- Never log passwords, tokens, or credentials
- Mask or truncate personally identifiable information (PII)
- Use dedicated sanitization methods for logging user data

### 7. Logger Initialization

- Initialize a single root logger in the application entry point
- Configure log levels from environment variables
- Pass logger instances to components via constructors

```go
// In main.go
logger, err := zap.NewProduction()
if err != nil {
    log.Fatalf("can't initialize zap logger: %v", err)
}
defer logger.Sync()
```

## Consequences

### Positive

- Consistent, structured logs across all services
- Improved troubleshooting with detailed context
- Better correlation of events across service boundaries
- Efficient log searching and filtering
- Reduced time to identify and resolve issues

### Negative

- Small performance overhead from structured logging
- Learning curve for proper logging practices
- Need for discipline in the following conventions

### Neutral

- Need to update existing services to follow conventions
- Regular review of logging patterns for effectiveness

## Implementation Notes

### Middleware for Request Context

Implement middleware to add request information to the context:

```go
func RequestInfoMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Generate request ID if not present
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
            c.Header("X-Request-ID", requestID)
        }

        // Add request info to context
        ctx := logctx.WithRequestInfo(c.Request.Context(), logctx.RequestInfo{
            ID:        requestID,
            Method:    c.Request.Method,
            Path:      c.Request.URL.Path,
            UserAgent: c.Request.UserAgent(),
            ClientIP:  c.ClientIP(),
        })

        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### Helper for Contextualized Logging

Implement a package for extracting request context:

```go
// Package logctx provides utilities for contextual logging
package logctx

import (
    "context"
    "go.uber.org/zap"
)

type RequestInfo struct {
    ID        string
    Method    string
    Path      string
    UserAgent string
    ClientIP  string
}

type ctxKey string
const requestInfoKey ctxKey = "requestInfo"

func WithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
    return context.WithValue(ctx, requestInfoKey, info)
}

func LoggerWithRequestInfo(ctx context.Context, logger *zap.Logger) *zap.Logger {
    info, ok := ctx.Value(requestInfoKey).(RequestInfo)
    if !ok {
        return logger
    }

    return logger.With(
        zap.String("request_id", info.ID),
        zap.String("method", info.Method),
        zap.String("path", info.Path),
    )
}
```

## Related Documents

- [Service Structure](./0001-service-structure.md)
- [Interface Design Principles](./0002-interface-design-principles.md)
- [Dependency Injection](./0003-dependency-injection.md)

## Contributors

- Àlex Grau Roca
