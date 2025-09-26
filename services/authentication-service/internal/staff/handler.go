package staff

import (
	"github.com/gin-gonic/gin"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Handler manages HTTP requests for auth-customer-related operations.
type Handler struct {
	logger         log.Logger
	service        Service
	authMiddleware auth.Middleware
}

// NewHandler creates a new instance of Handler.
func NewHandler(logger log.Logger, service Service, authMiddleware auth.Middleware) *Handler {
	return &Handler{
		logger:         logger,
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registers the customer-related HTTP routes.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
}
