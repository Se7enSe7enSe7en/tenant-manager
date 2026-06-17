package model

// this is the view model for the create tenant handler
type CreateTenantSignals struct {
	PropertyId      string `json:"property_id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	ExpectedRentDay int    `json:"expected_rent_day"`
}
