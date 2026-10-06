// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: outbox_events.sql

package db

import (
	"time"
)

type ClaimPendingOutboxEventsParams struct {
	Limit         int32     `db:"limit" json:"limit"`
	NextAttemptAt time.Time `db:"next_attempt_at" json:"next_attempt_at"`
}

type CreateOutboxEventParams struct {
	Topic    string `db:"topic" json:"topic"`
	EventKey string `db:"event_key" json:"event_key"`
	Payload  []byte `db:"payload" json:"payload"`
}

type MarkOutboxEventFailedParams struct {
	OutboxID      int64     `db:"outbox_id" json:"outbox_id"`
	NextAttemptAt time.Time `db:"next_attempt_at" json:"next_attempt_at"`
}
