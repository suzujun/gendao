package helper

import (
	"io/ioutil"
	"os"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestFile_CreateDirIfNotExist(t *testing.T) {
	require := assert.New(t)
	path, err := ioutil.TempDir("", "")
	require.NoError(err)

	assert := assert.New(t)
	assert.NoError(CreateDirIfNotExist(path))

	info, err := os.Stat(path)
	require.NoError(err)
	require.True(info.IsDir())
}

func TestFile_CreateFile_readFile(t *testing.T) {
	require := assert.New(t)
	dir, err := ioutil.TempDir("", "")
	require.NoError(err)
	path := dir + "/foo.txt"

	assert := assert.New(t)
	v, err := CreateFile(path, "hello world")
	assert.NoError(err)
	assert.True(v > 0)

	b, err := ReadFile(path)
	assert.NoError(err)
	assert.Equal(string(b), "hello world")
}

func TestFile_ReadFileJSON(t *testing.T) {
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
	v, err := CreateFile(textPath, "hello world")
	assert.NoError(err)
	assert.True(v > 0)

	assert.Error(ReadFileJSON(textPath, &importData))
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
	v, err = CreateFile(jsonPath, b)
	assert.NoError(err)
	assert.True(v > 0)
	assert.NoError(ReadFileJSON(jsonPath, &importData))
	assert.Equal(importData, data)
}
