package store

type CreatePropertySignals struct {
	Name       string  `json:"name"`        // Required
	RentAmount float64 `json:"rent_amount"` // Required, convert later to decimal.Decimal
	Address    string  `json:"address"`

	// for error handling
	FormError string `json:"form_error"`
}
