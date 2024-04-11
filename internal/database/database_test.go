package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Call the New function
	defaultTestConfig := &Config{
		Database: "testgochallenge",
		AppEnv:   "test",
	}

	service, err := New(defaultTestConfig)

	// Assert that the error is nil
	assert.Nil(t, err)

	// Assert that the returned service is not nil
	assert.NotNil(t, service)
}
