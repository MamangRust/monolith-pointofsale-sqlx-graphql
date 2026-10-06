package convert

import (
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// EnvOr returns the value of the environment variable key if set, otherwise fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// NullableString converts an empty string to a nil pointer.
func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// PgTimestamp parses a Postgres timestamp string ("2006-01-02 15:04:05" or
// RFC3339) into a pgtype.Timestamp. Empty or unparsable input yields an
// invalid (zero) timestamp.
func PgTimestamp(s string) pgtype.Timestamp {
	if s == "" {
		return pgtype.Timestamp{}
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return pgtype.Timestamp{}
		}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}
