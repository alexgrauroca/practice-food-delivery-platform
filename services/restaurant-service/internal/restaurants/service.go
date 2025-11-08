package restaurants

import (
	"context"

	"github.com/alexgrauroca/practice-food-delivery-platform/pkg/log"
	"github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/staff"
)

// Service represents the interface for operations related to restaurant management.
//
//go:generate mockgen -destination=./mocks/service_mock.go -package=restaurants_mocks github.com/alexgrauroca/practice-food-delivery-platform/services/restaurant-service/internal/restaurants Service
type Service interface {
	RegisterRestaurant(ctx context.Context, input RegisterRestaurantInput) (RegisterRestaurantOutput, error)
}

type service struct {
	logger    log.Logger
	repo      Repository
	staffServ staff.Service
}

func NewService(logger log.Logger, repo Repository, staffServ staff.Service) Service {
	return &service{
		logger:    logger,
		repo:      repo,
		staffServ: staffServ,
	}
}

func (s service) RegisterRestaurant(ctx context.Context, input RegisterRestaurantInput) (RegisterRestaurantOutput, error) {
	logger := s.logger.WithContext(ctx)
	logger.Info(
		"registering restaurant",
		log.Field{Key: "vat_code", Value: input.Restaurant.VatCode},
		log.Field{Key: "name", Value: input.Restaurant.Name},
	)

	params := CreateRestaurantParams{
		VatCode:    input.Restaurant.VatCode,
		Name:       input.Restaurant.Name,
		LegalName:  input.Restaurant.LegalName,
		TaxID:      input.Restaurant.TaxID,
		TimezoneID: input.Restaurant.TimezoneID,
		Contact:    CreateContactParams(input.Restaurant.Contact),
	}
	_, err := s.repo.CreateRestaurant(ctx, params)
	if err != nil {
		logger.Error("failed to create restaurant", err)
		return RegisterRestaurantOutput{}, err
	}

	return RegisterRestaurantOutput{}, nil
}
