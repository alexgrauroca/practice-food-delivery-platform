package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

type contextKey string

const (
	authHeader               = "Authorization"
	bearerPrefix             = "Bearer "
	subjectCtxKey contextKey = "token-subject"
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
	return func(c *gin.Context) {
		claims, err := m.validateToken(c)
		if err != nil {
			m.handleAuthError(c, err)
			return
		}

		if claims.Role != string(RoleCustomer) {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				newErrorResponse(CodeForbiddenError, MessageForbiddenError),
			)
			return
		}

		c.Set(string(subjectCtxKey), claims.Subject)
		ctx := context.WithValue(c.Request.Context(), subjectCtxKey, claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (m *middleware) validateToken(c *gin.Context) (*Claims, error) {
	token, err := extractBearerToken(c)
	if err != nil {
		return nil, err
	}

	output, err := m.service.ValidateAccessToken(c.Request.Context(), ValidateAccessTokenInput{AccessToken: token})
	if err != nil {
		return nil, err
	}

	return output.Claims, nil

}

func (m *middleware) handleAuthError(c *gin.Context, err error) {
	code := CodeUnauthorizedError
	msg := MessageUnauthorizedError
	if errors.Is(err, ErrTokenExpired) {
		code = CodeForbiddenError
		msg = MessageForbiddenError
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, newErrorResponse(code, msg))
}

func extractBearerToken(c *gin.Context) (string, error) {
	header := c.GetHeader(authHeader)
	if header == "" {
		return "", ErrAuthHeaderMissing
	}

	if !strings.HasPrefix(header, bearerPrefix) {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(header, bearerPrefix), nil
}
