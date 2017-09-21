package helper

import (
	"database/sql"
	"strings"
)

type Uniquer struct {
	data map[string]bool
}

func ParseIntPointer(val *sql.NullInt64) *uint {
	if val.Valid {
		uintVal := uint(val.Int64)
		return &uintVal
	}
	return nil
}

func ParseStringPointer(val *sql.NullString) *string {
	if val.Valid {
		return &val.String
	}
	return nil
}

func StringsContains(s interface{}, substr string) bool {
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

func Unique(s []string) []string {
	uq := NewUniquer()
	for _, v := range s {
		uq.Add(v)
	}
	return uq.Uniq()
}

// -----------------
// Uniquer
// -----------------

func NewUniquer() *Uniquer {
	return &Uniquer{
		data: map[string]bool{},
	}
}
func (u *Uniquer) Add(s string) {
	u.data[s] = true
}

func (u *Uniquer) Uniq() []string {
	var i int
	res := make([]string, len(u.data))
	for v := range u.data {
		res[i] = v
		i++
	}
	return res
}
