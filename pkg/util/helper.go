package util

import (
	"database/sql"
)

// nullString converts a string to sql.NullString for empty string check
func NullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

// nullUint32 converts an uint32 to sql.NullInt32 for empty uint32 check
func NullUint32(value uint32) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	valueInt64 := int64(value)

	return sql.NullInt64{
		Int64: valueInt64,
		Valid: true,
	}
}

// nullUint64 converts an uint64 to sql.NullInt64 for empty uint64 check
func NullUint64(value uint64) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	valueInt64 := int64(value)

	return sql.NullInt64{
		Int64: valueInt64,
		Valid: true,
	}
}

// nullInt64 converts an int64 to sql.NullInt64 for empty int64 check
func NullInt64(value int64) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}

// nullFloat64 converts a float64 to sql.NullFloat64 for empty float64 check
func NullFloat64(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}
