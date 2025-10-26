package staff

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/auth"
	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

const (
	// CodeStaffAlreadyExists represents the error code indicating the staff already exists in the system.
	CodeStaffAlreadyExists = "STAFF_ALREADY_EXISTS"
	// MsgStaffAlreadyExists represents the error message indicating that the staff already exists in the system.
	MsgStaffAlreadyExists = "staff already exists"
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
}

// RegisterStaffResponse represents the response returned after successfully registering a new staff user.
type RegisterStaffResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterStaff handles the registration of a new staff user.
func (h *Handler) RegisterStaff(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("RegisterStaff handler called")

	var req RegisterStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RegisterStaffInput(req)

	output, err := h.service.RegisterStaff(ctx, input)
	if err != nil {
		if errors.Is(err, ErrStaffAlreadyExists) {
			logger.Warn("Staff already exists", log.Field{Key: "email", Value: req.Email})
			c.JSON(http.StatusConflict, customhttp.NewErrorResponse(CodeStaffAlreadyExists, MsgStaffAlreadyExists))
			return
		}
		logger.Error("Failed to register staff", err)
		c.JSON(
			http.StatusInternalServerError,
			customhttp.NewErrorResponse(customhttp.CodeInternalError, customhttp.MsgInternalError),
		)
		return
	}

	resp := RegisterStaffResponse(output)
	logger.Info("Staff registered successfully", log.Field{Key: "staff", Value: resp})
	c.JSON(http.StatusCreated, resp)
}
