package middleware

import (
	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		host := c.Request.Host
		realIP := c.GetHeader("X-Real-IP")
		if realIP == "" {
			realIP = c.ClientIP()
		}

		ctx := logctx.WithRequestInfo(c.Request.Context(), requestID, host, realIP)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
