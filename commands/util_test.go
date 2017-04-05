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
}
