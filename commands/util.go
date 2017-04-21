package commands

import (
	"database/sql"
	"strings"
)

func parseIntPointer(val *sql.NullInt64) *uint {
	if val.Valid {
		uintVal := uint(val.Int64)
		return &uintVal
	}
	return nil
}

func parseStringPointer(val *sql.NullString) *string {
	if val.Valid {
		return &val.String
	}
	return nil
}

func stringsContains(s interface{}, substr string) bool {
	switch d := s.(type) {
	case string:
		return strings.Contains(d, substr)
	case []string:
		for _, a := range d {
			if a == substr {
				return true
			}
		}
		return false
	default:
		return false
	}
}
