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

	return nil
}
