package mapperx

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgUUIDToString(id pgtype.UUID) (string, error) {
	if !id.Valid {
		return "", fmt.Errorf("uuid is null/invalid")
	}
	u, err := uuid.FromBytes(id.Bytes[:])
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func StringToPgUUID(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}
