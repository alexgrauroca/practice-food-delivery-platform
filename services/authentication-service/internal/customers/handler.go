package customers

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
	// CodeCustomerAlreadyExists represents the error code indicating the customer already exists in the system.
	CodeCustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"
	// MsgCustomerAlreadyExists represents the error message indicating that the customer already exists in the system.
	MsgCustomerAlreadyExists = "customer already exists"

	// CodeInvalidCredentials represents the error code for failed authentication due to invalid login credentials.
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	// MsgInvalidCredentials represents the error message returned when login authentication fails due to invalid credentials.
	MsgInvalidCredentials = "invalid credentials"

	// CodeInvalidRefreshToken represents the error code for an invalid or expired refresh token used in authentication processes.
	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	// MsgInvalidRefreshToken represents an error message indicating an invalid or expired refresh token.
	MsgInvalidRefreshToken = "invalid or expired refresh token"

	// CodeTokenMismatch represents an error code indicating a mismatch between the provided token and the expected value.
	CodeTokenMismatch = "TOKEN_MISMATCH"
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
		authRouter.POST("/customers", h.RegisterCustomer)
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
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RegisterCustomerInput(req)

	output, err := h.service.RegisterCustomer(ctx, input)
	if err != nil {
		if errors.Is(err, ErrCustomerAlreadyExists) {
			logger.Warn("Customer already exists", log.Field{Key: "email", Value: req.Email})
			c.JSON(http.StatusConflict, customhttp.NewErrorResponse(CodeCustomerAlreadyExists, MsgCustomerAlreadyExists))
			return
		}
		logger.Error("Failed to register customer", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(
			customhttp.CodeInternalError,
			customhttp.MsgInternalError,
		))
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
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := LoginCustomerInput(req)
	output, err := h.service.LoginCustomer(ctx, input)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			logger.Warn("Invalid credentials provided", log.Field{Key: "email", Value: req.Email})
			c.JSON(http.StatusUnauthorized, customhttp.NewErrorResponse(CodeInvalidCredentials, MsgInvalidCredentials))
			return
		}
		logger.Error("Failed to login customer", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(
			customhttp.CodeInternalError,
			customhttp.MsgInternalError,
		))
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
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := RefreshCustomerInput(req)
	output, err := h.service.RefreshCustomer(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			logger.Warn("Invalid refresh token provided")
			c.JSON(http.StatusUnauthorized, customhttp.NewErrorResponse(CodeInvalidRefreshToken, MsgInvalidRefreshToken))
			return
		} else if errors.Is(err, ErrTokenMismatch) {
			logger.Warn("Token mismatch")
			c.JSON(http.StatusForbidden, customhttp.NewErrorResponse(CodeTokenMismatch, MsgTokenMismatch))
			return
		}

		logger.Error("Failed to refresh customer", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(
			customhttp.CodeInternalError,
			customhttp.MsgInternalError,
		))
		return
	}

	resp := RefreshCustomerResponse{TokenPairResponse: TokenPairResponse(output.TokenPair)}
	logger.Info("Customer refreshed successfully")
	c.JSON(http.StatusOK, resp)
}
