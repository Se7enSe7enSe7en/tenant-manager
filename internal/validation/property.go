package validation

import (
	"errors"
)

type CreatePropertyForm struct {
	Name       string
	RentAmount float64
}

func CheckCreatePropertyForm(form CreatePropertyForm) error {
	// TODO: validate Name, should have limit

	// validate RentAmount, max value in the DB is 99,999,999.99
	if form.RentAmount > 99_999_999.99 {
		return errors.New("rent amount is too high, max amount is 99,999,999.99")
	}

	return nil
}
