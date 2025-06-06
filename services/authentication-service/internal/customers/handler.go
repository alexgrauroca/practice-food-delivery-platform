package customers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RegisterCustomerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type RegisterCustomerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1.0/customers", h.RegisterCustomer)
}

func (h *Handler) RegisterCustomer(c *gin.Context) {
	h.logger.Info("RegisterCustomer handler called")

	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		errResp := h.getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	//TODO implement the use case logic for registering a customer

	//TODO remove the mocked response
	resp := RegisterCustomerResponse{
		ID:        "fake-id",
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	h.logger.Info("Customer registered successfully", zap.Any("customer", resp))
	c.JSON(http.StatusCreated, resp)
}

// getErrorResponseFromValidationErr gets the ErrorResponse based on the error type returned from the validation
func (h *Handler) getErrorResponseFromValidationErr(err error) ErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return ErrorResponse{Error: ve.Error()}
	}
	return ErrorResponse{Error: "Invalid request"}
}
