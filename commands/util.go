package commands

import (
	"database/sql"
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

func stringsContains(s []string, substr string) bool {
	for _, a := range s {
		if a == substr {
			return true
		}
	}
	return false
}
