package commands

import (
	"database/sql"
	"strings"

	"github.com/suzujun/gendao/inflector"
)

type inflect int

var singularize inflect = 1
var pluralize inflect = 2

// ConvSingularize convert singularize
func ConvSingularize(str string) string {
	if str == "" {
		return ""
	}
	items := strings.Split(str, "_")
	last := len(items) - 1
	items[last] = inflector.Singularize(items[last])
	return strings.Join(items, "_")
}

// ConvPascalcase convert pascal case
func ConvPascalcase(str string, lastLint bool) string {
	return strings.Title(convCamelcaseOption(str, nil, lastLint))
}

// ConvPascalcaseSingularize convert pascal case and singularize
func ConvPascalcaseSingularize(str string, lastLint bool) string {
	return strings.Title(convCamelcaseOption(str, &singularize, lastLint))
}

// ConvPascalcasePluralize convert pascal case and pluralize
func ConvPascalcasePluralize(str string, lastLint bool) string {
	return strings.Title(convCamelcaseOption(str, &pluralize, lastLint))
}

// ConvCamelcase convert camel case
func ConvCamelcase(str string, lastLint bool) string {
	return convCamelcaseOption(str, nil, lastLint)
}

// ConvCamelcaseSingularize convert camel case and singularize
func ConvCamelcaseSingularize(str string, lastLint bool) string {
	return convCamelcaseOption(str, &singularize, lastLint)
}

// ConvCamelcasePluralize convert camel case and pluralize
func ConvCamelcasePluralize(str string, lastLint bool) string {
	return convCamelcaseOption(str, &pluralize, lastLint)
}

func convCamelcaseOption(str string, inflect *inflect, lint bool) string {
	items := strings.Split(str, "_")
	for i, v := range items {
		if i == len(items)-1 {
			if inflect != nil {
				if *inflect == singularize {
					v = inflector.Singularize(v)
				} else if *inflect == pluralize {
					v = inflector.Pluralize(v)
				}
			}
			if lint && commonInitialisms[strings.ToUpper(v)] {
				items[i] = strings.ToUpper(v)
				continue
			}
		}
		if i == 0 {
			items[i] = strings.ToLower(v)
		} else {
			items[i] = strings.Title(v)
		}
	}
	return strings.Join(items, "")
}

func stringsJoinPascalcase(items []string, sep string) string {
	res := []string{}
	for _, item := range items {
		res = append(res, strings.Title(item))
	}
	return strings.Join(res, sep)
}

func stringsFilter(items []string, fn func(string) bool) []string {
	res := []string{}
	for _, item := range items {
		if fn(item) {
			res = append(res, item)
		}
	}
	return res
}

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

// @see https://github.com/golang/lint/blob/b8599f7d71e7fead76b25aeb919c0e2558672f4a/lint.go#L717-L759
// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}
