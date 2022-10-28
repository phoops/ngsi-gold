package client_test

import (
	"testing"
	"time"

	"github.com/phoops/ngsild/client"
	"github.com/stretchr/testify/assert"
)

func TestClientConstructorSuccess(t *testing.T) {
	setURLOpt := client.SetURL("localhost:9090")
	c, err := client.New(setURLOpt)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

// URL Options

func TestClientMissingURL(t *testing.T) {
	c, err := client.New()
	assert.Error(t, err)
	assert.Nil(t, c)
	assert.ErrorIs(t, err, client.ErrMissingURL)
}

// Timeout
func TestClientTimeoutSuccess(t *testing.T) {
	setURLOpt := client.SetURL("localhost:9090")
	setTimeoutOpt := client.SetClientTimeout(10 * time.Second)
	c, err := client.New(setURLOpt, setTimeoutOpt)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestClientNegativeTimeout(t *testing.T) {
	setURLOpt := client.SetURL("localhost:9090")
	setTimeoutOpt := client.SetClientTimeout(-10 * time.Second)
	c, err := client.New(setURLOpt, setTimeoutOpt)
	assert.Error(t, err)
	assert.Nil(t, c)
	assert.ErrorIs(t, err, client.ErrNegativeTimeout)
}
