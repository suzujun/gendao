package commands

import (
	"strings"

	"github.com/suzujun/gendao/inflector"
)

type WordConverter struct {
	word        string
	singularize bool
	pluralize   bool
	camelcase   bool
	pascalcase  bool
	lint        bool
}

func NewWordConverter(word string) *WordConverter {
	return &WordConverter{word: word}
}

func (wc *WordConverter) Singularize() *WordConverter {
	wc.singularize = true
	wc.pluralize = false
	return wc
}

func (wc *WordConverter) Pluralize() *WordConverter {
	wc.singularize = false
	wc.pluralize = true
	return wc
}

func (wc *WordConverter) Camelcase() *WordConverter {
	wc.camelcase = true
	wc.pascalcase = false
	return wc
}

func (wc *WordConverter) Pascalcase() *WordConverter {
	wc.camelcase = false
	wc.pascalcase = true
	return wc
}

func (wc *WordConverter) Lint() *WordConverter {
	wc.lint = true
	return wc
}

func (wc *WordConverter) ToString() string {
	if wc.word == "" {
		return ""
	}
	sep := "_"
	items := strings.Split(wc.word, sep)
	lastIndex := len(items) - 1
	// singularize
	if wc.singularize {
		items[lastIndex] = inflector.Singularize(items[lastIndex])
	}
	// camel calse, pascal case
	if wc.camelcase || wc.pascalcase {
		for i, v := range items {
			items[i] = strings.Title(v)
		}
		if wc.camelcase && len(items[0]) > 0 {
			items[0] = strings.ToLower(items[0][:1]) + items[0][1:]
		}
		sep = ""
	}
	// lint
	if wc.lint && (!wc.camelcase || wc.camelcase && len(items) > 1) {
		items[lastIndex] = getLint(items[lastIndex])
	}
	// pluralize
	if wc.pluralize {
		items[lastIndex] = inflector.Pluralize(items[lastIndex])
	}
	return strings.Join(items, sep)
}

func getLint(word string) string {
	if word == "" {
		return ""
	}
	upper := strings.ToUpper(word)
	if lintMap[upper] {
		return upper
	} else if lintMap[upper[:len(word)-1]] && upper[len(word)-1:len(word)] == "S" {
		return upper[:len(word)-1] + "s"
	}
	return word
}

// @see https://github.com/golang/lint/blob/b8599f7d71e7fead76b25aeb919c0e2558672f4a/lint.go#L717-L759
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var lintMap = map[string]bool{
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
