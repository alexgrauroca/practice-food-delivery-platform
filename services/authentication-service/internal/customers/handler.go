package customers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	CodeValidationError       = "VALIDATION_ERROR"
	CodeInvalidRequest        = "INVALID_REQUEST"
	CodeCustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"
	CodeInternalError         = "INTERNAL_ERROR"
	CodeInvalidCredentials    = "INVALID_CREDENTIALS"

	MsgValidationError          = "validation failed"
	MsgInvalidRequest           = "invalid request"
	MsgCustomerAlreadyExists    = "customer already exists"
	MsgFailedToRegisterCustomer = "failed to register the customer"
	MsgInvalidCredentials       = "invalid credentials"
	MsgFailedToLoginCustomer    = "failed to login the customer"
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

type LoginCustomerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginCustomerResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"` // the number of seconds until the token expires
	TokenType string `json:"token_type"`
}

type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
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
	router.POST("/v1.0/customers/register", h.RegisterCustomer)
	router.POST("/v1.0/customers/login", h.LoginCustomer)
}

func (h *Handler) RegisterCustomer(c *gin.Context) {
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("RegisterCustomer handler called")

	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Failed to bind request", zap.Error(err))
		errResp := getErrorResponseFromValidationErr(err)
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
			logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
				Warn("Customer already exists", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, newErrorResponse(CodeCustomerAlreadyExists, MsgCustomerAlreadyExists))
			return
		}
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
			Error("Failed to register customer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgFailedToRegisterCustomer))
		return
	}

	resp := RegisterCustomerResponse{
		ID:        output.ID,
		Email:     output.Email,
		Name:      output.Name,
		CreatedAt: output.CreatedAt,
	}
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
		Info("Customer registered successfully", zap.Any("customer", resp))
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) LoginCustomer(c *gin.Context) {
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("LoginCustomer handler called")

	var req LoginCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Failed to bind request", zap.Error(err))
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := LoginCustomerInput{
		Email:    req.Email,
		Password: req.Password,
	}
	output, err := h.service.LoginCustomer(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
				Warn("Invalid credentials provided", zap.String("email", req.Email))
			c.JSON(http.StatusUnauthorized, newErrorResponse(CodeInvalidCredentials, MsgInvalidCredentials))
			return
		}
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
			Error("Failed to login customer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgFailedToLoginCustomer))
		return
	}

	resp := LoginCustomerResponse{
		Token:     output.Token,
		ExpiresIn: output.ExpiresIn,
		TokenType: output.TokenType,
	}
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
		Info("Customer logged in successfully")
	c.JSON(http.StatusOK, resp)
}

func newErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Details: make([]string, 0),
	}
}

// getErrorResponseFromValidationErr gets the ErrorResponse based on the error type returned from the validation
func getErrorResponseFromValidationErr(err error) ErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		errResp := newErrorResponse(CodeValidationError, MsgValidationError)
		details := make([]string, 0)

		for _, fe := range ve {
			details = append(details, getValidationErrorDetail(fe))
		}
		errResp.Details = details

		return errResp
	}
	return newErrorResponse(CodeInvalidRequest, MsgInvalidRequest)
}

// getValidationErrorDetail returns a detailed error message based on the field error
func getValidationErrorDetail(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		if field == "password" {
			return field + " must be a valid password with at least 8 characters long"
		}
		return field + " must be at least " + fe.Param() + " characters long" //notest
	default: //notest
		return field + " is invalid"
	}
}
