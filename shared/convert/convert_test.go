package convert

import (
	"testing"
	"time"
)

func TestNullableString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want *string
	}{
		{name: "empty becomes nil", in: "", want: nil},
		{name: "non-empty keeps value", in: "hello", want: ptr("hello")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableString(tt.in)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("NullableString(%q) = %v, want nil", tt.in, *got)
				}
				return
			}
			if got == nil {
				t.Fatalf("NullableString(%q) = nil, want %q", tt.in, *tt.want)
			}
			if *got != *tt.want {
				t.Fatalf("NullableString(%q) = %q, want %q", tt.in, *got, *tt.want)
			}
		})
	}
}

func TestPgTimestamp(t *testing.T) {
	valid := func(s string) bool {
		ts := PgTimestamp(s)
		return ts.Valid
	}
	tests := []struct {
		name string
		in   string
		want bool // valid
	}{
		{name: "empty is invalid", in: "", want: false},
		{name: "db format parses", in: "2026-01-02 15:04:05", want: true},
		{name: "rfc3339 parses", in: "2026-01-02T15:04:05Z", want: true},
		{name: "garbage is invalid", in: "not-a-date", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := valid(tt.in)
			if got != tt.want {
				t.Fatalf("PgTimestamp(%q).Valid = %v, want %v", tt.in, got, tt.want)
			}
		})
	}

	// Parsed time must round-trip the original layout.
	ts := PgTimestamp("2026-01-02 15:04:05")
	if !ts.Valid {
		t.Fatal("PgTimestamp(db format) not valid")
	}
	want := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	if !ts.Time.Equal(want) {
		t.Fatalf("PgTimestamp time = %v, want %v", ts.Time, want)
	}
}

func ptr(s string) *string { return &s }
