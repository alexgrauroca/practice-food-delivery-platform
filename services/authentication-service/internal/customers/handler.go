package customers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"go.uber.org/zap"

	"github.com/alexgrauroca/practice-food-delivery-platform/services/authentication-service/internal/logctx"
)

const (

	// CodeValidationError represents the error code for validation failures during input processing or validation checks.
	CodeValidationError = "VALIDATION_ERROR"
	// CodeInvalidRequest represents the error code for an invalid or improper request made to the system.
	CodeInvalidRequest = "INVALID_REQUEST"
	// CodeCustomerAlreadyExists represents the error code indicating the customer already exists in the system.
	CodeCustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"
	// CodeInternalError represents the error code for an unspecified internal server error encountered in the system.
	CodeInternalError = "INTERNAL_ERROR"
	// CodeInvalidCredentials represents the error code for failed authentication due to invalid login credentials.
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	// CodeInvalidRefreshToken represents the error code for an invalid or expired refresh token used in authentication processes.
	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	// CodeTokenMismatch represents an error code indicating a mismatch between the provided token and the expected value.
	CodeTokenMismatch = "TOKEN_MISMATCH"

	// MsgValidationError represents the error message for validation failures during input validation checks.
	MsgValidationError = "validation failed"
	// MsgInvalidRequest represents the error message for an invalid or improperly formed request.
	MsgInvalidRequest = "invalid request"
	// MsgCustomerAlreadyExists represents the error message indicating that the customer already exists in the system.
	MsgCustomerAlreadyExists = "customer already exists"
	// MsgInvalidCredentials represents the error message returned when login authentication fails due to invalid credentials.
	MsgInvalidCredentials = "invalid credentials"
	// MsgInternalError represents the error message returned when the system fails to log in a customer.
	MsgInternalError = "an unexpected error occurred"
	// MsgInvalidRefreshToken represents an error message indicating an invalid or expired refresh token.
	MsgInvalidRefreshToken = "invalid or expired refresh token"
	// MsgTokenMismatch represents the error message for a token mismatch scenario.
	MsgTokenMismatch = "token mismatch"
)

// RegisterCustomerRequest represents the request payload for registering a new customer.
type RegisterCustomerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// RegisterCustomerResponse represents the response returned after successfully registering a new customer.
type RegisterCustomerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginCustomerRequest represents the request payload for logging in a customer.
type LoginCustomerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginCustomerResponse represents the response payload for a successful customer login.
type LoginCustomerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // the number of seconds until the token expires
	TokenType    string `json:"token_type"`
}

// RefreshCustomerRequest represents a request to refresh customer information using tokens.
type RefreshCustomerRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
}

// RefreshCustomerResponse represents the response returned when refreshing a customer's token.
type RefreshCustomerResponse struct {
	LoginCustomerResponse
}

// ErrorResponse represents a standardized structure for API error responses containing code, message, and optional details.
type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

// Handler manages HTTP requests for auth-customer-related operations.
type Handler struct {
	logger  *zap.Logger
	service Service
}

// NewHandler creates a new instance of Handler.
func NewHandler(logger *zap.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// RegisterRoutes registers the customer-related HTTP routes.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1.0/customers/register", h.RegisterCustomer)
	router.POST("/v1.0/customers/login", h.LoginCustomer)
	router.POST("v1.0/customers/refresh", h.RefreshCustomer)
}

// RegisterCustomer handles the registration of a new customer.
func (h *Handler) RegisterCustomer(c *gin.Context) {
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("RegisterCustomer handler called")

	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Failed to bind request", zap.Error(err))
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RegisterCustomerInput(req)

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
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := RegisterCustomerResponse(output)
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
		Info("Customer registered successfully", zap.Any("customer", resp))
	c.JSON(http.StatusCreated, resp)
}

// LoginCustomer processes the login request for a customer using credentials provided in JSON format.
func (h *Handler) LoginCustomer(c *gin.Context) {
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("LoginCustomer handler called")

	var req LoginCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Failed to bind request", zap.Error(err))
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := LoginCustomerInput(req)
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
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := LoginCustomerResponse(output.TokenPair)
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
		Info("Customer logged in successfully")
	c.JSON(http.StatusOK, resp)
}

// RefreshCustomer handles the refreshing of a customer's authentication token.
func (h *Handler) RefreshCustomer(c *gin.Context) {
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("RefreshCustomer handler called")

	var req RefreshCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Failed to bind request", zap.Error(err))
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RefreshCustomerInput(req)
	output, err := h.service.RefreshCustomer(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Invalid refresh token provided")
			c.JSON(http.StatusUnauthorized, newErrorResponse(CodeInvalidRefreshToken, MsgInvalidRefreshToken))
			return
		} else if errors.Is(err, ErrTokenMismatch) {
			logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Warn("Token mismatch")
			c.JSON(http.StatusForbidden, newErrorResponse(CodeTokenMismatch, MsgTokenMismatch))
			return
		}

		logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).
			Error("Failed to refresh customer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := RefreshCustomerResponse{
		LoginCustomerResponse(output.TokenPair),
	}
	logctx.LoggerWithRequestInfo(c.Request.Context(), h.logger).Info("Customer refreshed successfully")
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
	field := strcase.ToSnake(fe.Field())
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
