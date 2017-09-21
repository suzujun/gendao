package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_StringsContains(t *testing.T) {
	assert := assert.New(t)
	strs := []string{"aaa", "bbb", "ccc"}

	// []strings: match
	assert.True(StringsContains(strs, "bbb"))

	// []strings: unmatch
	assert.False(StringsContains(strs, "ddd"))

	// []strings: unmatch
	assert.False(StringsContains(strs, ""))

	// string: match
	assert.True(StringsContains("hogefuga", "fu"))

	// string: match to empty
	assert.True(StringsContains("hogefuga", ""))

	// string: unmatch
	assert.False(StringsContains("hogefuga", "piyo"))

	// int: unmatch
	assert.False(StringsContains(12345, "234"))
}

func TestUtil_Unique(t *testing.T) {

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
			res := Unique(test.args)
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
