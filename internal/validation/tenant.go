package validation

import (
	"errors"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/model"
)

func CheckCreateTenantForm(form model.CreateTenantSignals) error {
	// TODO: complete validation

	if form.PropertyId == "" {
		return errors.New("No property selected")
	}

	if form.ExpectedRentDay < 1 || form.ExpectedRentDay > 31 {
		return errors.New("Choose a day between 1 and 31")
	}

	return nil
}
