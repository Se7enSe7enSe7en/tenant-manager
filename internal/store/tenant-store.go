package store

// this is the store/view model for the create tenant handler
type CreateTenantSignals struct {
	// for tenant table
	PropertyId  string `json:"property_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`

	// for lease table
	ExpectedRentDay int     `json:"expected_rent_day"`
	StartDate       string  `json:"start_date"`
	IsMonthAdvance  bool    `json:"is_month_advance"`
	DepositAmount   float64 `json:"deposit_amount"` // later converted to type decimal.Decimal
	NextDueDate     string  `json:"next_due_date"`
}

// note: make sure the name of the json field is the same in the front end, eg. `json:"field_name"`
