package customers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
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

// TokenPairResponse represents the structure for holding both access and refresh tokens along with metadata.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // the number of seconds until the token expires
	TokenType    string `json:"token_type"`
}

// Handler manages HTTP requests for auth-customer-related operations.
type Handler struct {
	logger  log.Logger
	service Service
}

// NewHandler creates a new instance of Handler.
func NewHandler(logger log.Logger, service Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// RegisterRoutes registers the customer-related HTTP routes.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/v1.0/auth")
	{
		auth.POST("/customers", h.RegisterCustomer)
	}

	// Existing routes remain for backward compatibility
	router.POST("/v1.0/customers/register", h.RegisterCustomer)
	router.POST("/v1.0/customers/login", h.LoginCustomer)
	router.POST("v1.0/customers/refresh", h.RefreshCustomer)
}

// RegisterCustomerRequest represents the request payload for registering a new customer.
type RegisterCustomerRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Name       string `json:"name" binding:"required"`
}

// RegisterCustomerResponse represents the response returned after successfully registering a new customer.
type RegisterCustomerResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterCustomer handles the registration of a new customer.
func (h *Handler) RegisterCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("RegisterCustomer handler called")

	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RegisterCustomerInput(req)

	output, err := h.service.RegisterCustomer(ctx, input)
	if err != nil {
		if errors.Is(err, ErrCustomerAlreadyExists) {
			logger.Warn("Customer already exists", log.Field{Key: "email", Value: req.Email})
			c.JSON(http.StatusConflict, newErrorResponse(CodeCustomerAlreadyExists, MsgCustomerAlreadyExists))
			return
		}
		logger.Error("Failed to register customer", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := RegisterCustomerResponse(output)
	logger.Info("Customer registered successfully", log.Field{Key: "customer", Value: resp})
	c.JSON(http.StatusCreated, resp)
}

// LoginCustomerRequest represents the request payload for logging in a customer.
type LoginCustomerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginCustomerResponse represents the response payload for a successful customer login.
type LoginCustomerResponse struct {
	TokenPairResponse
}

// LoginCustomer processes the login request for a customer using credentials provided in JSON format.
func (h *Handler) LoginCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("LoginCustomer handler called")

	var req LoginCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := LoginCustomerInput(req)
	output, err := h.service.LoginCustomer(ctx, input)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			logger.Warn("Invalid credentials provided", log.Field{Key: "email", Value: req.Email})
			c.JSON(http.StatusUnauthorized, newErrorResponse(CodeInvalidCredentials, MsgInvalidCredentials))
			return
		}
		logger.Error("Failed to login customer", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := LoginCustomerResponse{TokenPairResponse: TokenPairResponse(output.TokenPair)}
	logger.Info("Customer logged in successfully")
	c.JSON(http.StatusOK, resp)
}

// RefreshCustomerRequest represents a request to refresh customer information using tokens.
type RefreshCustomerRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
}

// RefreshCustomerResponse represents the response returned when refreshing a customer's token.
type RefreshCustomerResponse struct {
	TokenPairResponse
}

// RefreshCustomer handles the refreshing of a customer's authentication token.
func (h *Handler) RefreshCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("RefreshCustomer handler called")

	var req RefreshCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := getErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RefreshCustomerInput(req)
	output, err := h.service.RefreshCustomer(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			logger.Warn("Invalid refresh token provided")
			c.JSON(http.StatusUnauthorized, newErrorResponse(CodeInvalidRefreshToken, MsgInvalidRefreshToken))
			return
		} else if errors.Is(err, ErrTokenMismatch) {
			logger.Warn("Token mismatch")
			c.JSON(http.StatusForbidden, newErrorResponse(CodeTokenMismatch, MsgTokenMismatch))
			return
		}

		logger.Error("Failed to refresh customer", err)
		c.JSON(http.StatusInternalServerError, newErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := RefreshCustomerResponse{TokenPairResponse: TokenPairResponse(output.TokenPair)}
	logger.Info("Customer refreshed successfully")
	c.JSON(http.StatusOK, resp)
}

// ErrorResponse represents a standardized structure for API error responses containing code, message, and optional details.
type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
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
