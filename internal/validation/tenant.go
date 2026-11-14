package validation

import (
	"errors"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/store"
)

var CreateTenantError = struct {
	NoTenantName           error
	NoProperty             error
	InvalidExpectedRentDay error
}{
	NoTenantName:           errors.New("Tenant Name is required"),
	NoProperty:             errors.New("No property selected"),
	InvalidExpectedRentDay: errors.New("Choose a day between 1 and 31"),
}

func CheckCreateTenantForm(form store.CreateTenantSignals) error {
	// TODO: complete validation

	if form.PropertyId == "" {
		return CreateTenantError.NoProperty
	}

	if form.Name == "" {
		return CreateTenantError.NoTenantName
	}

	if form.ExpectedRentDay < 1 || form.ExpectedRentDay > 31 {
		return CreateTenantError.InvalidExpectedRentDay
	}

	// TODO: check ExpectedRentDay, no decimals, whole numbers only

	return nil
}
