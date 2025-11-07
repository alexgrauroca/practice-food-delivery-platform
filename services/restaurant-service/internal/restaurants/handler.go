package restaurants

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	customhttp "github.com/alexgrauroca/practice-food-delivery-platform/pkg/http"
	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
)

const (
	// CodeRestaurantAlreadyExists represents the error code indicating the customer already exists in the system.
	CodeRestaurantAlreadyExists = "RESTAURANT_ALREADY_EXISTS"
	// MsgRestaurantAlreadyExists represents the error message indicating that the customer already exists in the system.
	MsgRestaurantAlreadyExists = "restaurant already exists"
)

// Handler manages HTTP requests for restaurant-related operations.
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

// RegisterRoutes registers the restaurant-related HTTP routes.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/v1.0/restaurants", h.RegisterRestaurant)
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

	input := RegisterRestaurantInput{
		Restaurant: RestaurantInput{
			VatCode:    req.Restaurant.VatCode,
			Name:       req.Restaurant.Name,
			LegalName:  req.Restaurant.LegalName,
			TaxID:      req.Restaurant.TaxID,
			TimezoneID: req.Restaurant.TimezoneID,
			Contact:    ContactInput(req.Restaurant.Contact),
		},
		StaffOwner: StaffOwnerInput(req.StaffOwner),
	}
	output, err := h.service.RegisterRestaurant(ctx, input)
	if err != nil {
		if errors.Is(err, ErrRestaurantAlreadyExists) {
			logger.Warn("Restaurant already exists", log.Field{Key: "vat_code", Value: req.Restaurant.VatCode})
			c.JSON(http.StatusConflict, customhttp.NewErrorResponse(
				CodeRestaurantAlreadyExists,
				MsgRestaurantAlreadyExists,
			))
			return
		}
		logger.Error("Failed to register restaurant", err)
		c.JSON(http.StatusInternalServerError, customhttp.NewErrorResponse(
			customhttp.CodeInternalError,
			customhttp.MsgInternalError,
		))
		return
	}

	resp := RegisterRestaurantResponse{
		Restaurant: RestaurantResponse{
			ID:         output.Restaurant.ID,
			VatCode:    output.Restaurant.VatCode,
			Name:       output.Restaurant.Name,
			LegalName:  output.Restaurant.LegalName,
			TaxID:      output.Restaurant.TaxID,
			TimezoneID: output.Restaurant.TimezoneID,
			Contact:    ContactResponse(output.Restaurant.Contact),
			CreatedAt:  output.Restaurant.CreatedAt,
			UpdatedAt:  output.Restaurant.UpdatedAt,
		},
		StaffOwner: StaffOwnerResponse(output.StaffOwner),
	}
	logger.Info("Restaurant registered successfully", log.Field{Key: "restaurant", Value: resp})
	c.JSON(http.StatusCreated, resp)
}
