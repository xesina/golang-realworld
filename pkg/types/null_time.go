package types

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

// NullTime represents a time.Time that may be null.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// MarshalJSON try to marshalize to json
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	}

	return []byte("null"), nil
}

// UnmarshalJSON try to unmarshal data from input
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == "null" {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}

	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
