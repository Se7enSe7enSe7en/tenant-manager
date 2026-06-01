package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgTypeUuid(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)

	return uuid, err
}
