package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString represents a string that may be null.
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// MarshalJSON try to marshalize to json
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return []byte("null"), nil
}

// UnmarshalJSON try to unmarshal data from input
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	if err != nil {
		return err
	}
	ns.Valid = ns.String != ""

	return nil
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	temp := &sql.NullString{}
	if err := temp.Scan(value); err != nil {
		return err
	}

	ns.String = temp.String
	ns.Valid = temp.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// NewNullString return nullstring structure
func NewNullString(str string) NullString {
	return NullString{
		String: str,
		Valid:  true,
	}
}
