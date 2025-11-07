package restaurants

import (
	"net/http"

	"github.com/gin-gonic/gin"

	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

// Handler manages HTTP requests for restaurant-related operations.
type Handler struct {
	logger log.Logger
}

// NewHandler creates a new instance of Handler.
func NewHandler(logger log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// RegisterRoutes registers the restaurant-related HTTP routes.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1.0/restaurants", h.RegisterRestaurant)
}

// RegisterRestaurantRequest represents the request payload for registering a new restaurant.
type RegisterRestaurantRequest struct {
	Restaurant struct {
		VatCode    string `json:"vat_code" binding:"required,max=40"`
		Name       string `json:"name" binding:"required,max=100"`
		LegalName  string `json:"legal_name" binding:"required,max=100"`
		TaxID      string `json:"tax_id" binding:"max=40"`
		TimezoneID string `json:"timezone_id" binding:"required,iana_tz"`
		Contact    struct {
			PhonePrefix string `json:"phone_prefix" binding:"required,phone_pref"`
			PhoneNumber string `json:"phone_number" binding:"required,phone_num"`
			Email       string `json:"email" binding:"required,email"`
			Address     string `json:"address" binding:"required,max=100"`
			City        string `json:"city" binding:"required,max=100"`
			PostalCode  string `json:"postal_code" binding:"required,min=5,max=32"`
			CountryCode string `json:"country_code" binding:"required,min=2,max=2"`
		} `json:"contact" binding:"required"`
	} `json:"restaurant" binding:"required"`
	StaffOwner struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		Name        string `json:"name" binding:"required,max=100"`
		Address     string `json:"address" binding:"required,max=100"`
		City        string `json:"city" binding:"required,max=100"`
		PostalCode  string `json:"postal_code" binding:"required,min=5,max=32"`
		CountryCode string `json:"country_code" binding:"required,min=2,max=2"`
	} `json:"staff_owner" binding:"required"`
}

// RegisterRestaurantResponse represents the response returned after successfully registering a new restaurant.
type RegisterRestaurantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// RegisterRestaurant handles the registration of a new restaurant and staff owner.
func (h *Handler) RegisterRestaurant(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.WithContext(ctx)

	logger.Info("RegisterRestaurant handler called")

	var req RegisterRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request", log.Field{Key: "error", Value: err.Error()})
		errResp := customhttp.GetErrorResponseFromValidationErr(err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
}
