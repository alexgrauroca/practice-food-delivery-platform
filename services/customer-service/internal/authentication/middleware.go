package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/log"
)

const (
	authHeader    = "Authorization"
	bearerPrefix  = "Bearer "
	subjectCtxKey = "token-subject"
)

// Middleware defines the interface for authentication-related middleware functions used
// to secure and validate requests.
//
//go:generate mockgen -destination=./mocks/middleware_mock.go -package=authentication_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/customer-service/internal/authentication Middleware
type Middleware interface {
	RequireCustomer() gin.HandlerFunc
}

type middleware struct {
	logger  log.Logger
	service Service
}

// NewMiddleware creates a new instance of Middleware.
// It handles authentication-related operations and requests validation.
func NewMiddleware(logger log.Logger, service Service) Middleware {
	return &middleware{
		logger:  logger,
		service: service,
	}
}

func (m *middleware) RequireCustomer() gin.HandlerFunc {
	// TODO: review the implementation
	return func(c *gin.Context) {
		claims, err := m.validateToken(c)
		if err != nil {
			m.handleAuthError(c, err)
			return
		}

		if claims.Role != string(RoleCustomer) {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(subjectCtxKey, claims.Subject)
		c.Next()
	}
}

func GetSubject(ctx context.Context) (string, bool) {
	// TODO: review the implementation
	v := ctx.Value(subjectCtxKey)
	if v == nil {
		return "", false
	}
	subject, ok := v.(string)

	return subject, ok
}

func (m *middleware) validateToken(c *gin.Context) (*Claims, error) {
	// TODO: review the implementation
	token, err := extractBearerToken(c.GetHeader(authHeader))
	if err != nil {
		return nil, err
	}

	output, err := m.service.ValidateAccessToken(c.Request.Context(), ValidateAccessTokenInput{
		AccessToken: token,
	})
	if err != nil {
		return nil, err
	}

	return output.Claims, nil

}

func (m *middleware) handleAuthError(c *gin.Context, err error) {
	// TODO: review the implementation
	status := http.StatusUnauthorized
	code := "UNAUTHORIZED"
	msg := "Invalid credentials"

	if errors.Is(err, ErrExpiredToken) {
		status = http.StatusForbidden
		code = "EXPIRED_TOKEN"
		msg = "Token has expired"
	}

	c.AbortWithStatusJSON(status, gin.H{
		"code":    code,
		"message": msg,
	})
}

func extractBearerToken(header string) (string, error) {
	// TODO: review the implementation
	if header == "" {
		return "", errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(header, bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(header, bearerPrefix), nil
}
