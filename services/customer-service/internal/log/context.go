package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey string

const (
	requestIDKey ctxKey = "requestId"
	hostKey      ctxKey = "host"
	realIPKey    ctxKey = "realIp"
	userAgentKey ctxKey = "userAgent"
)

// RequestInfo represents metadata about a request, including its ID, host, and client's real IP address.
type RequestInfo struct {
	RequestID string
	Host      string
	RealIP    string
	UserAgent string
}

// MarshalLogObject serializes the RequestInfo fields into the provided zapcore.ObjectEncoder for structured logging.
func (r RequestInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("request_id", r.RequestID)
	enc.AddString("host", r.Host)
	enc.AddString("real_ip", r.RealIP)
	enc.AddString("user_agent", r.UserAgent)
	return nil
}

// WithRequestInfo adds request-specific information to the provided context and returns it.
func WithRequestInfo(ctx context.Context, info RequestInfo) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, info.RequestID)
	ctx = context.WithValue(ctx, hostKey, info.Host)
	ctx = context.WithValue(ctx, realIPKey, info.RealIP)
	ctx = context.WithValue(ctx, userAgentKey, info.UserAgent)
	return ctx
}

// RequestIDFromContext extracts the request ID from the provided context. Returns an empty string if not found.
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// HostFromContext extracts the host string from the provided context. Returns an empty string if the value is not present.
func HostFromContext(ctx context.Context) string {
	if host, ok := ctx.Value(hostKey).(string); ok {
		return host
	}
	return ""
}

// RealIPFromContext extracts the client's real IP address from the provided context. Returns an empty string if not found.
func RealIPFromContext(ctx context.Context) string {
	if realIP, ok := ctx.Value(realIPKey).(string); ok {
		return realIP
	}
	return ""
}

// UserAgentFromContext retrieves the user agent string from the given context if it exists, otherwise returns an empty string.
func UserAgentFromContext(ctx context.Context) string {
	if userAgent, ok := ctx.Value(userAgentKey).(string); ok {
		return userAgent
	}
	return ""
}

// LoggerWithRequestInfo enriches the provided logger with request metadata derived from the context.
func LoggerWithRequestInfo(ctx context.Context, logger *zap.Logger) *zap.Logger {
	reqInfo := RequestInfo{
		RequestID: RequestIDFromContext(ctx),
		Host:      HostFromContext(ctx),
		RealIP:    RealIPFromContext(ctx),
		UserAgent: UserAgentFromContext(ctx),
	}

	return logger.With(zap.Object("request", reqInfo))
}
