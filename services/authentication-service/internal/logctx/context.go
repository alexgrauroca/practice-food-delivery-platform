package logctx

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
)

type requestInfo struct {
	RequestID string
	Host      string
	RealIP    string
}

func (r requestInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("request_id", r.RequestID)
	enc.AddString("host", r.Host)
	enc.AddString("real_ip", r.RealIP)
	return nil
}

func WithRequestInfo(ctx context.Context, requestID, host, realIP string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	ctx = context.WithValue(ctx, hostKey, host)
	ctx = context.WithValue(ctx, realIPKey, realIP)
	return ctx
}

func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func HostFromContext(ctx context.Context) string {
	if host, ok := ctx.Value(hostKey).(string); ok {
		return host
	}
	return ""
}

func RealIPFromContext(ctx context.Context) string {
	if realIP, ok := ctx.Value(realIPKey).(string); ok {
		return realIP
	}
	return ""
}

func LoggerWithRequestInfo(ctx context.Context, logger *zap.Logger) *zap.Logger {
	reqInfo := requestInfo{
		RequestID: RequestIDFromContext(ctx),
		Host:      HostFromContext(ctx),
		RealIP:    RealIPFromContext(ctx),
	}

	return logger.With(zap.Object("request", reqInfo))
}
