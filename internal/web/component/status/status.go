// Shared go code between components in the status package

package status

import "github.com/Se7enSe7enSe7en/tenant-manager/internal/constants"

type CircleProps struct {
	Status constants.PaymentStatus
}

type BadgeProps struct {
	Status constants.PaymentStatus
}

func statusColorClassName(paymentStatus constants.PaymentStatus) string {
	className := ""

	switch paymentStatus {
	case constants.PAID:
		className = "bg-(--bg-paid-color) text-(--paid-color)"
	case constants.UNPAID:
		className = "bg-(--bg-unpaid-color) text-(--unpaid-color)"
	case constants.LATE:
		className = "bg-(--bg-late-color) text-(--late-color)"
	}

	return className
	// note: notice how we didn't do early return in the switch, this is because
	// we need the "className" variable for tailwindcss intellisense to work
}

func statusDotAnimationClassName(paymentStatus constants.PaymentStatus) string {
	className := ""

	if paymentStatus == constants.LATE {
		className = "animate-ping"
	}

	return className
}
