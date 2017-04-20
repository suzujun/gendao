package commands

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_NewConfig(t *testing.T) {
	assert := assert.New(t)

	conf := NewConfig("", "", "")
	assert.Equal(conf, newConfig())

	conf = NewConfig("test-user", "test-pass", "test-db")
	assert.NotEqual(conf, newConfig())
	assert.Equal(conf.MysqlConfig.User, "test-user")
	assert.Equal(conf.MysqlConfig.Password, "test-pass")
	assert.Equal(conf.MysqlConfig.DbName, "test-db")
}

func TestUtil_Write(t *testing.T) {
	require := assert.New(t)
	path, err := ioutil.TempDir("", "")
	require.NoError(err)

	assert := assert.New(t)

	conf := NewConfig("", "", "")
	assert.NoError(conf.Write(path + "/config.json"))
}

func TestUtil_ParseJSON(t *testing.T) {
	assert := assert.New(t)

	data := NewConfig("test-user", "test-pass", "test-db")
	b, err := json.Marshal(data)
	assert.NoError(err)
	assert.NotNil(b)

	conf := Config{}
	assert.NoError(conf.ParseJSON(b))
	assert.Equal(conf, data)
}
