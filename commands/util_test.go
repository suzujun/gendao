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

func TestUtil_unique(t *testing.T) {

	tests := []struct {
		title string
		args  []string
		want  []string
	}{
		{
			title: "empty",
			args:  []string{""},
			want:  []string{""},
		},
		{
			title: "one",
			args:  []string{"one"},
			want:  []string{"one"},
		},
		{
			title: "unduplicate",
			args:  []string{"aaa", "bbb", "ccc"},
			want:  []string{"aaa", "bbb", "ccc"},
		},
		{
			title: "duplicate",
			args:  []string{"aaa", "bbb", "bbb", "aaa", "aab"},
			want:  []string{"aaa", "bbb", "aab"},
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			assert := assert.New(t)
			res := unique(test.args)
			assert.Len(res, len(test.want))
			for _, v := range test.want {
				assert.Contains(res, v)
			}
			for _, v := range res {
				assert.Contains(test.want, v)
			}
		})
	}
}
