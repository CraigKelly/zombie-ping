package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := DecodeConfiguration([]byte(""))
	assert.Error(err)
	assert.Nil(config)
}

func TestSingleConfig(t *testing.T) {
	assert := assert.New(t)

	contents := `{ "Targets":[
        {"URL":"a", "PingSeconds":1 }
    ]}`

	config, err := DecodeConfiguration([]byte(contents))
	assert.NotNil(config)
	assert.NoError(err)

	assert.Equal(1, len(config.Targets))
	assert.Equal("a", config.Targets[0].URL)
	assert.Equal(1, config.Targets[0].PingSeconds)
}

func TestTempFileConfig(t *testing.T) {
	assert := assert.New(t)

	tmpfile, err := ioutil.TempFile("", "zombie-ping-config")
	assert.NoError(err)
	defer os.Remove(tmpfile.Name())

	contents := `{ "Targets":[
        {"URL":"a", "PingSeconds":1 },
        {"URL":"b", "PingSeconds":2 }
    ]}`
	assert.NoError(ioutil.WriteFile(tmpfile.Name(), []byte(contents), os.ModeTemporary))

	config, err := ReadConfiguration(tmpfile.Name())
	assert.NotNil(config)
	assert.NoError(err)

	assert.Equal(2, len(config.Targets))
	assert.Equal("a", config.Targets[0].URL)
	assert.Equal(1, config.Targets[0].PingSeconds)
	assert.Equal("b", config.Targets[1].URL)
	assert.Equal(2, config.Targets[1].PingSeconds)
}
