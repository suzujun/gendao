package commands

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestUtil_WordConverter(t *testing.T) {

	// plane
	tests := map[string]string{
		"":                "",
		"id":              "id",
		"channel_id":      "channel_id",
		"test_channel_id": "test_channel_id",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("plane[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).ToString()
			assert.Equal(res, answer)
		})
	}

	// plane & lint
	tests = map[string]string{
		"":                "",
		"id":              "ID",
		"channel_id":      "channel_ID",
		"test_channel_id": "test_channel_ID",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("plane_lint[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Lint().ToString()
			assert.Equal(res, answer)
		})
	}

	// camelcase
	tests = map[string]string{
		"":                "",
		"id":              "id",
		"channel_id":      "channelId",
		"test_channel_id": "testChannelId",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("camelcase[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Camelcase().ToString()
			assert.Equal(res, answer)
		})
	}

	// camelcase & pluralize
	tests = map[string]string{
		"":                "",
		"id":              "ids",
		"channel_id":      "channelIds",
		"test_channel_id": "testChannelIds",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("camelcase_pluralize[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Camelcase().Pluralize().ToString()
			assert.Equal(res, answer)
		})
	}

	// camelcase & pluralize & lint
	tests = map[string]string{
		"":                "",
		"id":              "ids",
		"channel_id":      "channelIDs",
		"test_channel_id": "testChannelIDs",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("camelcase_pluralize_lint[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Camelcase().Pluralize().Lint().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase
	tests = map[string]string{
		"":                "",
		"id":              "Id",
		"channel_id":      "ChannelId",
		"test_channel_id": "TestChannelId",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase & pluralize
	tests = map[string]string{
		"":                "",
		"id":              "Ids",
		"channel_id":      "ChannelIds",
		"test_channel_id": "TestChannelIds",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase_pluralize[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().Pluralize().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase & lint
	tests = map[string]string{
		"":                 "",
		"id":               "ID",
		"channel_id":       "ChannelID",
		"channel_ids":      "ChannelIDs",
		"test_channel_id":  "TestChannelID",
		"test_channel_ids": "TestChannelIDs",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase_lint[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().Lint().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase & pluralize & lint
	tests = map[string]string{
		"":                "",
		"id":              "IDs",
		"channel_id":      "ChannelIDs",
		"test_channel_id": "TestChannelIDs",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase_pluralize_lint[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().Pluralize().Lint().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase & singularize
	tests = map[string]string{
		"":                 "",
		"ids":              "Id",
		"channel_ids":      "ChannelId",
		"test_channel_ids": "TestChannelId",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase_singularize[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().Singularize().ToString()
			assert.Equal(res, answer)
		})
	}

	// pascalcase & singularize & lint
	tests = map[string]string{
		"":                 "",
		"ids":              "ID",
		"channel_ids":      "ChannelID",
		"test_channel_ids": "TestChannelID",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("pascalcase_singularize_lint[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := NewWordConverter(test).Pascalcase().Singularize().Lint().ToString()
			assert.Equal(res, answer)
		})
	}
}

func TestUtil_getLint(t *testing.T) {

	// plane
	tests := map[string]string{
		"":      "",
		"id":    "ID",
		"ids":   "IDs",
		"hoge":  "hoge",
		"json":  "JSON",
		"jsons": "JSONs",
	}
	for test, answer := range tests {
		title := fmt.Sprintf("plane[%s]", test)
		t.Run(title, func(t *testing.T) {
			assert := assert.New(t)
			res := getLint(test)
			assert.Equal(res, answer)
		})
	}
}
