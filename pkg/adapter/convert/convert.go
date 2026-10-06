// Package convert holds the proto→domain conversion helpers shared by the
// gRPC client adapters under pkg/adapter/<domain>.
package convert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Timestamp parses a timestamp string emitted by a gRPC service into a
// pgtype.Timestamp, the type sqlc generates for the timestamp columns. It
// accepts the layouts the services are known to emit and returns the zero
// (invalid) Timestamp when the input is empty or cannot be parsed, so callers
// keep a usable model with only the textual fields set.
func Timestamp(s string) pgtype.Timestamp {
	if s == "" {
		return pgtype.Timestamp{}
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return pgtype.Timestamp{Time: t, Valid: true}
		}
	}
	return pgtype.Timestamp{}
}

// NullableString returns nil for an empty string, otherwise a pointer to it.
func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// NullableInt32 returns a pointer to i.
func NullableInt32(i int32) *int32 {
	return &i
}
