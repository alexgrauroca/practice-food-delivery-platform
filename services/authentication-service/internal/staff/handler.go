package staff

import (
	"time"

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
	authRouter := router.Group("/v1.0/auth")
	{
		authRouter.POST("/staff", h.RegisterStaff)
	}
}

// RegisterStaffRequest represents the request payload for registering a new staff user.
type RegisterStaffRequest struct {
	StaffID  string `json:"staff_id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// RegisterStaffResponse represents the response returned after successfully registering a new staff user.
type RegisterStaffResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterStaff handles the registration of a new staff user.
func (h *Handler) RegisterStaff(c *gin.Context) {

}
