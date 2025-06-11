// Package middleware provides HTTP middleware components for handling cross-cutting
// concerns in the authentication service
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
)

// RequestInfoMiddleware is a middleware that attaches request-specific information to the context.
func RequestInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		realIP := c.GetHeader("X-Real-IP")
		if realIP == "" {
			realIP = c.ClientIP()
		}

		info := logctx.RequestInfo{
			RequestID: requestID,
			Host:      c.Request.Host,
			RealIP:    realIP,
			UserAgent: c.Request.UserAgent(),
		}

		ctx := logctx.WithRequestInfo(c.Request.Context(), info)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
