// Code migrated from sqlc-generated sources; domain types kept for repository layer.
// source: consumer_inbox.sql

package db

type MarkConsumerInboxProcessedParams struct {
	ConsumerName       string `db:"consumer_name" json:"consumer_name"`
	EventKey           string `db:"event_key" json:"event_key"`
	ReservationVersion int64  `db:"reservation_version" json:"reservation_version"`
}

type ReleaseConsumerInboxParams struct {
	ConsumerName       string `db:"consumer_name" json:"consumer_name"`
	EventKey           string `db:"event_key" json:"event_key"`
	LastError          string `db:"last_error" json:"last_error"`
	ReservationVersion int64  `db:"reservation_version" json:"reservation_version"`
}

type ReserveConsumerInboxParams struct {
	ConsumerName  string `db:"consumer_name" json:"consumer_name"`
	EventKey      string `db:"event_key" json:"event_key"`
	Topic         string `db:"topic" json:"topic"`
	PartitionID   int32  `db:"partition_id" json:"partition_id"`
	MessageOffset int64  `db:"message_offset" json:"message_offset"`
}

type ReserveConsumerInboxRow struct {
	Reserved           bool  `db:"reserved" json:"reserved"`
	Processed          bool  `db:"processed" json:"processed"`
	ReservationVersion int64 `db:"reservation_version" json:"reservation_version"`
}
