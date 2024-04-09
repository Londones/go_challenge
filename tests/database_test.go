package tests

import (
	"testing"

	"go-challenge/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Call the New function
	service, err := database.New(&database.Config{
		Database: "testgochallenge",
		AppEnv:   "test",
	})

	// Assert that the error is nil
	assert.Nil(t, err)

	// Assert that the returned service is not nil
	assert.NotNil(t, service)

	errorTD := database.Service.TearDown(service, "testgochallenge")

	// Assert that the error is nil
	assert.Nil(t, errorTD)

	database.Service.Close(service)
}

func TestNewError(t *testing.T) {
	// Call the New function
	service, err := database.New(&database.Config{
		Database: "testgochallenge",
		AppEnv:   "test",
	})

	// Assert that the error is nil
	assert.Nil(t, err)

	// Assert that the returned service is not nil
	assert.NotNil(t, service)

	// Call the New function again
	_, err = database.New(&database.Config{
		Database: "testgochallenge",
		AppEnv:   "test",
	})

	// Assert that the error is not nil
	assert.NotNil(t, err)

	database.Service.Close(service)
	database.Service.TearDown(service, "testgochallenge")

}
