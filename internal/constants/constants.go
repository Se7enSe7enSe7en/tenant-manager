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

type TradeType int

const (
	RENT TradeType = iota
	DEPOSIT
)

func (t TradeType) String() string {
	switch t {
	case RENT:
		return "Rent"
	case DEPOSIT:
		return "Deposit"
	default:
		return ""

	}
}
