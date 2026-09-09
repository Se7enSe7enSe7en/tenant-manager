package tenantcard

import "github.com/Se7enSe7enSe7en/tenant-manager/internal/constants"

type TenantCardProps struct {
	Id          string
	Name        string
	Unit        string
	Status      constants.PaymentStatus
	RentAmount  string
	NextDueDate string
	Email       *string
	PhoneNumber *string
	LeaseId     string
}
