package utils

import (
	"time"

	"github.com/Se7enSe7enSe7en/tenant-manager/internal/constants"
)

func ComputeNextExpiryDate(prevDate time.Time, expectedRentDay int) time.Time {
	// set the start date to 1st day of the month
	startDateNoDay := time.Date(
		prevDate.Year(),
		prevDate.Month(),
		1,
		0, 0, 0, 0, // remove H:M:S:Ns
		prevDate.Location(),
	)

	// get next month
	oneMonthAfterStartDate := startDateNoDay.AddDate(0, 1, 0)

	expiryDate := time.Date(
		oneMonthAfterStartDate.Year(),
		oneMonthAfterStartDate.Month(),
		expectedRentDay,
		oneMonthAfterStartDate.Hour(),
		oneMonthAfterStartDate.Minute(),
		oneMonthAfterStartDate.Second(),
		oneMonthAfterStartDate.Nanosecond(),
		oneMonthAfterStartDate.Location(),
	)

	// capped normalization logic
	if expiryDate.Month() != oneMonthAfterStartDate.Month() {
		expiryDate = time.Date(
			oneMonthAfterStartDate.Year(),
			oneMonthAfterStartDate.Month()+1, // add 1 since we want the last day of the next month
			0,                                // "0" means the last day of the previous month (eg. month = Oct, day = 0 -> Sep 31 )
			oneMonthAfterStartDate.Hour(),
			oneMonthAfterStartDate.Minute(),
			oneMonthAfterStartDate.Second(),
			oneMonthAfterStartDate.Nanosecond(),
			oneMonthAfterStartDate.Location(),
		)
	}
	// note: by default golang normalizes the dates,
	// if the ExpectedRentDay is "31" on a month without 31 (eg. Feb)
	// it will roll over to the next month (Feb 31 = Mar 2), we don't want this behavior,
	// we want it to be capped to the last day of the month instead (Feb 31 = Feb 28)

	return expiryDate
}

func ComputeStatus(expiryDate time.Time) constants.PaymentStatus {
	now := time.Now()

	if now.Before(expiryDate) {
		return constants.PAID
	}

	// expiryDate + 1 month = 1 month late
	if now.After(expiryDate.AddDate(0, 1, 0)) {
		return constants.LATE
	}

	return constants.UNPAID
}
