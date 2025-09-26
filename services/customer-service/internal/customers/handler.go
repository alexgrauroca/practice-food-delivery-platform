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
	// CodeInternalError represents the error code for an unspecified internal server error encountered in the system.
	CodeInternalError = "INTERNAL_ERROR"
	// CodeNotFound represents the error code indicating that the requested resource could not be found in the system.
	CodeNotFound = "NOT_FOUND"

	// MsgCustomerAlreadyExists represents the error message indicating that the customer already exists in the system.
	MsgCustomerAlreadyExists = "customer already exists"
	// MsgInternalError represents the error message returned when the system fails to log in a customer.
	MsgInternalError = "an unexpected error occurred"
	// MsgNotFound represents the error message indicating that the requested resource could not be found.
	MsgNotFound = "resource not found"
)

// Handler manages HTTP requests for customer-related operations.
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
	router.POST("/v1.0/customers", h.RegisterCustomer)
	router.GET("/v1.0/customers/:customerID", h.authMiddleware.RequireCustomer(), h.GetCustomer)
	router.PUT("/v1.0/customers/:customerID", h.authMiddleware.RequireCustomer(), h.UpdateCustomer)
}

// GetCustomerResponse represents the response returned after successfully retrieving a customer.
type GetCustomerResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetCustomer handles retrieving a customer by CustomerID.
func (h *Handler) GetCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("GetCustomer handler called")

	customerID := c.Param("customerID")

	output, err := h.service.GetCustomer(ctx, GetCustomerInput{CustomerID: customerID})
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Warn("Customer not found", log.Field{Key: "customerID", Value: customerID})
			c.JSON(http.StatusNotFound, customhttp.NewErrorResponse(CodeNotFound, MsgNotFound))
			return
		}
		if errors.Is(err, ErrCustomerIDMismatch) {
			logger.Warn("Customer CustomerID mismatch with the token", log.Field{Key: "customerID", Value: customerID})
			errResp := customhttp.NewErrorResponse(auth.CodeForbiddenError, auth.MessageForbiddenError)
			c.JSON(http.StatusForbidden, errResp)
			return
		}
		logger.Error("Failed to get customer", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := GetCustomerResponse(output)
	logger.Info("Customer retrieved successfully", log.Field{Key: "customer", Value: resp})
	c.JSON(http.StatusOK, resp)
}

// RegisterCustomerRequest represents the request payload for registering a new customer.
type RegisterCustomerRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Name        string `json:"name" binding:"required,max=100"`
	Address     string `json:"address" binding:"required,max=100"`
	City        string `json:"city" binding:"required,max=100"`
	PostalCode  string `json:"postal_code" binding:"required,min=5,max=32"`
	CountryCode string `json:"country_code" binding:"required,min=2,max=2"`
}

// RegisterCustomerResponse represents the response returned after successfully registering a new customer.
type RegisterCustomerResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
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
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := RegisterCustomerResponse(output)
	logger.Info("Customer registered successfully", log.Field{Key: "customer", Value: resp})
	c.JSON(http.StatusCreated, resp)
}

// UpdateCustomerRequest represents the request payload for updating an existing customer's information.
type UpdateCustomerRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Address     string `json:"address" binding:"required,max=100"`
	City        string `json:"city" binding:"required,max=100"`
	PostalCode  string `json:"postal_code" binding:"required,min=5,max=32"`
	CountryCode string `json:"country_code" binding:"required,min=2,max=2"`
}

// UpdateCustomerResponse represents the response returned after successfully updating a customer's information.
type UpdateCustomerResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateCustomer handles updating an existing customer's information.
func (h *Handler) UpdateCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("UpdateCustomer handler called")

	customerID := c.Param("customerID")

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	input := UpdateCustomerInput{
		CustomerID:  customerID,
		Name:        req.Name,
		Address:     req.Address,
		City:        req.City,
		PostalCode:  req.PostalCode,
		CountryCode: req.CountryCode,
	}

	output, err := h.service.UpdateCustomer(ctx, input)
	if err != nil {
		if errors.Is(err, ErrCustomerNotFound) {
			logger.Warn("Customer not found", log.Field{Key: "customerID", Value: customerID})
			c.JSON(http.StatusNotFound, customhttp.NewErrorResponse(CodeNotFound, MsgNotFound))
			return
		}
		if errors.Is(err, ErrCustomerIDMismatch) {
			logger.Warn("Customer CustomerID mismatch with the token", log.Field{Key: "customerID", Value: customerID})
			errResp := customhttp.NewErrorResponse(auth.CodeForbiddenError, auth.MessageForbiddenError)
			c.JSON(http.StatusForbidden, errResp)
			return
		}
		logger.Error("Failed to update customer", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(CodeInternalError, MsgInternalError))
		return
	}

	resp := UpdateCustomerResponse(output)
	logger.Info("Customer updated successfully", log.Field{Key: "customer", Value: resp})
	c.JSON(http.StatusOK, resp)
}
