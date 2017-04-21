package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_stringsContains(t *testing.T) {
	assert := assert.New(t)
	strs := []string{"aaa", "bbb", "ccc"}

	// []strings: match
	assert.True(stringsContains(strs, "bbb"))

	// []strings: unmatch
	assert.False(stringsContains(strs, "ddd"))

	// []strings: unmatch
	assert.False(stringsContains(strs, ""))

	// string: match
	assert.True(stringsContains("hogefuga", "fu"))

	// string: match to empty
	assert.True(stringsContains("hogefuga", ""))

	// string: unmatch
	assert.False(stringsContains("hogefuga", "piyo"))

	// int: unmatch
	assert.False(stringsContains(12345, "234"))
}
