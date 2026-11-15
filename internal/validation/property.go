package validation

import (
	"errors"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
)

// type CreatePropertyForm struct {
// 	Name       string
// 	RentAmount decimal.Decimal
// }

// func CheckCreatePropertyForm(form CreatePropertyForm) error {
// 	// TODO: validate Name, should have limit

// 	// validate RentAmount, max value in the DB is 99,999,999.99
// 	if form.RentAmount.GreaterThan(decimal.NewFromFloat(99_999_999.99)) {
// 		return errors.New("rent amount is too high, max amount is 99,999,999.99")
// 	}

// 	return nil
// }

var CreatePropertyError = struct {
	NoPropertyName    error
	InvalidRentAmount error
}{
	NoPropertyName:    errors.New("no property name"),
	InvalidRentAmount: errors.New("invalid rent amount, min amount is 0, max amount is 99,999,999.99"),
}

func CheckCreatePropertySignals(signals store.CreatePropertySignals) error {
	if signals.Name == "" {
		return CreatePropertyError.NoPropertyName
	}

	if signals.RentAmount < 0 || signals.RentAmount > 99_999_999.99 {
		return CreatePropertyError.InvalidRentAmount
	}

	return nil
}
