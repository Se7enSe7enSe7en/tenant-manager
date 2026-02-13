package constants

type PaymentStatus int

const (
	PAID PaymentStatus = iota
	UNPAID
	LATE
)

func (ps PaymentStatus) String() string {
	switch ps {
	case PAID:
		return "Paid"
	case UNPAID:
		return "Unpaid"
	case LATE:
		return "Late"
	default:
		return ""
	}
}
