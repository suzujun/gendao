package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestFile_createDirIfNotExist(t *testing.T) {
	require := assert.New(t)
	path, err := ioutil.TempDir("", "")
	require.NoError(err)

	assert := assert.New(t)
	assert.NoError(createDirIfNotExist(path))

	info, err := os.Stat(path)
	require.NoError(err)
	require.True(info.IsDir())
}

func TestFile_createFile_readFile(t *testing.T) {
	require := assert.New(t)
	dir, err := ioutil.TempDir("", "")
	require.NoError(err)
	path := dir + "/foo.txt"

	assert := assert.New(t)
	v, err := createFile(path, "hello world")
	assert.NoError(err)
	assert.True(v > 0)

	b, err := readFile(path)
	assert.NoError(err)
	assert.Equal(string(b), "hello world")
}

func TestFile_readFileJSON(t *testing.T) {
	require := assert.New(t)
	dir, err := ioutil.TempDir("", "")
	require.NoError(err)

	assert := assert.New(t)

	type TestData struct {
		StringValue string
		IntValue    int
		BoolValue   bool
	}
	var importData TestData

	// text file
	textPath := dir + "/data.txt"
	v, err := createFile(textPath, "hello world")
	assert.NoError(err)
	assert.True(v > 0)

	assert.Error(readFileJSON(textPath, &importData))
	assert.Empty(importData.StringValue)

	// json file
	data := TestData{
		StringValue: "hogege",
		IntValue:    1200,
		BoolValue:   true,
	}
	b, err := json.Marshal(data)
	assert.NoError(err)
	jsonPath := dir + "/data.json"
	v, err = createFile(jsonPath, b)
	assert.NoError(err)
	assert.True(v > 0)
	assert.NoError(readFileJSON(jsonPath, &importData))
	assert.Equal(importData, data)
}
