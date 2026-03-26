package constants

type PaymentStatus int

const (
	UNPAID PaymentStatus = iota
	PAID
	LATE
)

func (ps PaymentStatus) String() string {
	switch ps {
	case UNPAID:
		return "Unpaid"
	case PAID:
		return "Paid"
	case LATE:
		return "Late"
	default:
		return ""
	}
}
