package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_ConvPascalcase(t *testing.T) {
	assert := assert.New(t)
	var res string

	// case empty
	res = ConvPascalcase("", false)
	assert.Empty(res)

	// case one
	res = ConvPascalcase("id", false)
	assert.Equal(res, "Id")

	// case one and lastlint
	res = ConvPascalcase("id", true)
	assert.Equal(res, "ID")

	// case two
	res = ConvPascalcase("hello_id", false)
	assert.Equal(res, "HelloId")

	// case one and lastlint
	res = ConvPascalcase("hello_id", true)
	assert.Equal(res, "HelloID")
}

func TestUtil_ConvPascalcaseSingularize(t *testing.T) {
	assert := assert.New(t)
	var res string

	// case empty
	res = ConvPascalcaseSingularize("", false)
	assert.Empty(res)

	// case one
	res = ConvPascalcaseSingularize("ids", false)
	assert.Equal(res, "Id")

	// case one and lastlint
	res = ConvPascalcaseSingularize("ids", true)
	assert.Equal(res, "ID")

	// case two
	res = ConvPascalcaseSingularize("hello_ids", false)
	assert.Equal(res, "HelloId")

	// case one and lastlint
	res = ConvPascalcaseSingularize("hello_ids", true)
	assert.Equal(res, "HelloID")
}

func TestUtil_stringsContains(t *testing.T) {
	assert := assert.New(t)
	strs := []string{"aaa", "bbb", "ccc"}

	// []strings: match
	assert.True(stringsContains(strs, "bbb"))

	// []strings: unmatch
	assert.False(stringsContains(strs, "ddd"))

	// []strings: unmatch
	assert.False(stringsContains(strs, ""))
}

func TestUtil_ConvPascalcasePluralize(t *testing.T) {
	assert := assert.New(t)
	var res string

	// case empty
	res = ConvPascalcasePluralize("", false)
	assert.Empty(res)

	// case one
	res = ConvPascalcasePluralize("id", false)
	assert.Equal(res, "Ids")

	// case one and lastlint
	res = ConvPascalcasePluralize("id", true)
	assert.Equal(res, "Ids")

	// case two
	res = ConvPascalcasePluralize("hello_id", false)
	assert.Equal(res, "HelloIds")

	// case one and lastlint
	res = ConvPascalcasePluralize("hello_id", true)
	assert.Equal(res, "HelloIds")
}

// func camelcaseToArray(str string) []string {
// 	var bi int
// 	res := []string{}
// 	for i, c := range str {
// 		if 65 <= c && c <= 90 { // A-Z[65-90]
// 			res = append(res, strings.ToLower(str[bi:i]))
// 			bi = i
// 		}
// 	}
// 	if len(str) > 0 {
// 		res = append(res, strings.ToLower(str[bi:]))
// 	}
// 	return res
// }

// func stringsJoinPascalcase(items []string, sep string) string {
// 	res := []string{}
// 	for _, item := range items {
// 		res = append(res, strings.Title(item))
// 	}
// 	return strings.Join(res, sep)
// }

// func stringsFilter(items []string, fn func(string) bool) []string {
// 	res := []string{}
// 	for _, item := range items {
// 		if fn(item) {
// 			res = append(res, item)
// 		}
// 	}
// 	return res
// }

// func parseIntPointer(val *sql.NullInt64) *uint {
// 	if val.Valid {
// 		uintVal := uint(val.Int64)
// 		return &uintVal
// 	}
// 	return nil
// }

// func parseStringPointer(val *sql.NullString) *string {
// 	if val.Valid {
// 		return &val.String
// 	}
// 	return nil
// }
