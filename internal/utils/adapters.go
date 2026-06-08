package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgTypeUuid(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)

	return uuid, err
}

func StringToPgTypeNumeric(s string) (pgtype.Numeric, error) {
	var num pgtype.Numeric
	err := num.Scan(s)

	return num, err
}
