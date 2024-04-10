package tests

import (
	"testing"

	"go-challenge/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Call the New function
	defaultTestConfig := &database.Config{
		Database: "testgochallenge",
		AppEnv:   "test",
	}

	service, err := database.New(defaultTestConfig)

	// Assert that the error is nil
	assert.Nil(t, err)

	// Assert that the returned service is not nil
	assert.NotNil(t, service)

	database.Service.Close(service)
	defaultTestConfig.Database = "postgres"
	errorTD := database.Service.TearDown(service, defaultTestConfig, "testgochallenge")

	// Assert that the error is nil
	assert.Nil(t, errorTD)

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
	database.Service.TearDown(service, &database.Config{
		Database: "postgres",
		AppEnv:   "test",
	}, "testgochallenge")

}
