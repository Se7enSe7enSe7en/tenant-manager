package utils

import (
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgtypeUuid(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)

	return uuid, err
}

func StringToPgtypeNumeric(s string) (pgtype.Numeric, error) {
	var num pgtype.Numeric
	err := num.Scan(s)

	return num, err
}

func PgtypeNumericToString(num pgtype.Numeric) string {
	const fallbackValue = "0"

	// If the database column value was NULL, handle the fallback state
	if !num.Valid {
		return fallbackValue
	}

	// Value() extracts the precise underlying string formatting natively
	val, err := num.Value()
	if err != nil {
		log.Printf("Failed to read numeric value: %v", err)
		return fallbackValue
	}

	// Type assert the driver.Value interface securely into a string
	if str, ok := val.(string); ok {
		return str
	}

	return fallbackValue
}
