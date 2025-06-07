package customers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	ErrorMsgCustomerAlreadyExists    = "Customer already exists"
	ErrorMsgFailedToRegisterCustomer = "Failed to register customer"
	ErrorMsgInvalidRequest           = "Invalid request"
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
	logger  *zap.Logger
	service Service
}

func NewHandler(logger *zap.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1.0/customers", h.RegisterCustomer)
}

func (h *Handler) RegisterCustomer(c *gin.Context) {
	//TODO review how can I include context info in the logs
	h.logger.Info("RegisterCustomer handler called")

	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Failed to bind request", zap.Error(err))
		errResp := h.getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RegisterCustomerInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	output, err := h.service.RegisterCustomer(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrCustomerAlreadyExists) {
			h.logger.Warn("Customer already exists", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, ErrorResponse{Error: ErrorMsgCustomerAlreadyExists})
			return
		}
		h.logger.Error("Failed to register customer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrorMsgFailedToRegisterCustomer})
		return
	}

	resp := RegisterCustomerResponse{
		ID:        output.ID,
		Email:     output.Email,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
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
	return ErrorResponse{Error: ErrorMsgInvalidRequest}
}
