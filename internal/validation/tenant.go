package validation

import (
	"errors"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
)

func CheckCreateTenantForm(form store.CreateTenantSignals) error {
	// TODO: complete validation

	if form.PropertyId == "" {
		return errors.New("No property selected")
	}

	if form.ExpectedRentDay < 1 || form.ExpectedRentDay > 31 {
		return errors.New("Choose a day between 1 and 31")
	}

	// TODO: check ExpectedRentDay, no decimals, whole numbers only

	return nil
}
