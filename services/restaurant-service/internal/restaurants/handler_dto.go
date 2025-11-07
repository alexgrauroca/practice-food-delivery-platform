package restaurants

import "time"

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
	Restaurant RestaurantResponse `json:"restaurant"`
	StaffOwner StaffOwnerResponse `json:"staff_owner"`
}

// RestaurantResponse represents the response returned after successfully retrieving a restaurant.
type RestaurantResponse struct {
	ID         string          `json:"id"`
	VatCode    string          `json:"vat_code"`
	Name       string          `json:"name"`
	LegalName  string          `json:"legal_name"`
	TaxID      string          `json:"tax_id"`
	TimezoneID string          `json:"timezone_id"`
	Contact    ContactResponse `json:"contact"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// ContactResponse represents the response returned after successfully retrieving a restaurant's contact information.
type ContactResponse struct {
	PhonePrefix string `json:"phone_prefix"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// StaffOwnerResponse represents the response returned after successfully retrieving a restaurant's staff owner.
type StaffOwnerResponse struct {
	ID          string    `json:"id"`
	Owner       bool      `json:"owner"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostalCode  string    `json:"postal_code"`
	CountryCode string    `json:"country_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
